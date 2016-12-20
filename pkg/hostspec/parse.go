package hostspec

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
)

type Spec struct {
	Name    string
	Foreman Foreman `mapstructure:"foreman"`
	Chef    Chef    `mapstructure:"chef"`
	Vsphere Vsphere `mapstructure:"vsphere"`
	Virtual Virtual `mapstructure:"virtual"`
}

type Devices struct {
	Disks    []*Disk    `mapstructure:"disk"`
	Networks []*Network `mapstructure:"network"`
	SCSIs    []*SCSI    `mapstructure:"scsi"`
}

type Foreman struct {
	Hostgroup         string `mapstructure:"hostgroup"`
	Location          string `mapstructure:"location"`
	Organization      string `mapstructure:"organization"`
	Environment       string `mapstructure:"environment"`
	ComputeProfile    string `mapstructure:"compute_profile"`
	ArchitectureID    int    `mapstructure:"architecture_id"`
	ComputeResource   string `mapstructure:"compute_resource"`
	DomainID          int    `mapstructure:"domain_id"`
	OperatingSystemID int    `mapstructure:"operating_system_id"`
	PartitionTableID  int    `mapstructure:"partition_table_id"`
	Medium            string `mapstructure:"medium"`
}

type Chef struct {
	Environment string   `mapstructure:"environment"`
	BaseRole    string   `mapstructure:"base_role"`
	RunList     []string `mapstructure:"run_list"`
}

type Vsphere struct {
	Domain     string  `mapstructure:"domain"`
	Cluster    string  `mapstructure:"cluster"`
	Datastore  string  `mapstructure:"datastore"`
	Folder     string  `mapstructure:"folder"`
	Datacenter string  `mapstructure:"datacenter"`
	Devices    Devices `mapstructure:"device"`
}

type Virtual struct {
	CPUs   int `mapstructure:"cpus"`
	Cores  int `mapstructure:"cores"`
	Memory int `mapstructure:"memory"`
}

type Disk struct {
	DeviceName string
	DeviceType string
	Size       int `mapstructure:"size"`
}

type Network struct {
	DeviceName string
	DeviceType string
	BuildVLAN  string `mapstructure:"build_vlan"`
	VLAN       string `mapstructure:"vlan"`
	SwitchType string `mapstructure:"switch_type"`
}

type SCSI struct {
	DeviceName string
	DeviceType string
	Type       string `mapstructure:"type"`
}

func ParseDir(path, spec string) (*Spec, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		hostspec, err := ParseFile(file.Name())
		if err != nil {
			return nil, err
		}

		if hostspec.Name != spec {
			continue
		} else {
			return hostspec, nil
		}
	}

	return nil, nil
}

// ParseFile parses the given hostspec file.
func ParseFile(path string) (*Spec, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

// Parse parses the hostspec from the given io.Reader.
//
// Due to current internal limitations, the entire contents of the
// io.Reader will be copied into memory first before parsing.
func Parse(r io.Reader) (*Spec, error) {
	// Copy the reader into an in-memory buffer first since HCL requires it.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	// Parse the buffer
	root, err := hcl.Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing: %s", err)
	}
	buf.Reset()

	// Top-level item should be a list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: root should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"spec",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	var spec Spec

	// Parse the spec out
	matches := list.Filter("spec")
	if len(matches.Items) == 0 {
		return nil, fmt.Errorf("%q stanza not found", "spec")
	}
	if err := parseSpec(&spec, matches); err != nil {
		return nil, fmt.Errorf("error parsing spec block: %s", err)
	}

	spec.Name = matches.Items[0].Keys[0].Token.Value().(string)

	return &spec, nil
}

// parseSpec parses the hostspec
func parseSpec(result *Spec, list *ast.ObjectList) error {
	list = list.Children()
	if len(list.Items) != 1 {
		return fmt.Errorf("only one %q block allowed, got %d spec blocks", "spec", len(list.Items))
	}

	// Get our "spec" object
	o := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := o.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	valid := []string{
		"foreman",
		"chef",
		"vsphere",
		"virtual",
	}
	if err := checkHCLKeys(listVal, valid); err != nil {
		return multierror.Prefix(err, "spec ->")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	delete(m, "foreman")
	delete(m, "chef")
	delete(m, "vsphere")
	delete(m, "virtual")

	var spec Spec
	if err := mapstructure.WeakDecode(m, &spec); err != nil {
		return err
	}

	// Parse out foreman fields
	if o := listVal.Filter("foreman"); len(o.Items) > 0 {
		if err := parseForeman(&spec.Foreman, o); err != nil {
			return multierror.Prefix(err, "foreman ->")
		}
	}

	// Parse out chef fields
	if o := listVal.Filter("chef"); len(o.Items) > 0 {
		if err := parseChef(&spec.Chef, o); err != nil {
			return multierror.Prefix(err, "chef ->")
		}
	}

	// Parse out vsphere fields
	if o := listVal.Filter("vsphere"); len(o.Items) > 0 {
		if err := parseVsphere(&spec.Vsphere, o); err != nil {
			return multierror.Prefix(err, "vsphere ->")
		}
	}

	// Parse out virtual fields
	if o := listVal.Filter("virtual"); len(o.Items) > 0 {
		if err := parseVirtual(&spec.Virtual, o); err != nil {
			return multierror.Prefix(err, "virtual ->")
		}
	}

	*result = spec
	return nil
}

func parseForeman(result *Foreman, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("Only one %q block allowed", "foreman")
	}

	// Get our "foreman" object
	o := list.Items[0]

	valid := []string{
		"hostgroup",
		"location",
		"organization",
		"environment",
		"compute_profile",
		"architecture_id",
		"compute_resource",
		"domain_id",
		"operating_system_id",
		"partition_table_id",
		"medium",
	}
	if err := checkHCLKeys(o.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var foreman Foreman
	if err := mapstructure.WeakDecode(m, &foreman); err != nil {
		return err
	}

	*result = foreman
	return nil
}

func parseChef(result *Chef, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("Only one %q block allowed", "chef")
	}

	// Get our "chef" object
	o := list.Items[0]

	valid := []string{
		"environment",
		"base_role",
		"run_list",
	}
	if err := checkHCLKeys(o.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var chef Chef
	if err := mapstructure.WeakDecode(m, &chef); err != nil {
		return err
	}

	*result = chef
	return nil
}

func parseVsphere(result *Vsphere, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block allowed", "vsphere")
	}

	// Get our vsphere object
	o := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := o.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	valid := []string{
		"domain",
		"cluster",
		"datastore",
		"folder",
		"datacenter",
		"device",
	}
	if err := checkHCLKeys(listVal, valid); err != nil {
		return multierror.Prefix(err, "vsphere ->")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	delete(m, "device")

	var vsphere Vsphere
	if err := mapstructure.WeakDecode(m, &vsphere); err != nil {
		return err
	}

	// Parse out device fields
	if o := listVal.Filter("device"); len(o.Items) > 0 {
		if err := parseDevices(&vsphere.Devices, o); err != nil {
			return multierror.Prefix(err, "device ->")
		}
	}

	*result = vsphere
	return nil
}

func parseVirtual(result *Virtual, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block allowed", "virtual")
	}

	// Get our virtual object
	o := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := o.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	valid := []string{
		"cpus",
		"cores",
		"memory",
	}
	if err := checkHCLKeys(listVal, valid); err != nil {
		return multierror.Prefix(err, "virtual ->")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var virtual Virtual
	if err := mapstructure.WeakDecode(m, &virtual); err != nil {
		return err
	}

	*result = virtual
	return nil
}

func parseDevices(result *Devices, list *ast.ObjectList) error {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil
	}

	var devices Devices

	seen := make(map[string]struct{})
	for _, item := range list.Items {
		if len(item.Keys) != 2 {
			return fmt.Errorf("%q must be followed by exactly two strings: a type and a name", "disk")
		}

		t := item.Keys[0].Token.Value().(string)
		n := item.Keys[1].Token.Value().(string)

		if _, ok := seen[n]; ok {
			return fmt.Errorf("key names should be unique: %q is defined more than once", n)
		}
		seen[n] = struct{}{}

		switch t {
		case "disk":
			if err := parseDisks(&devices.Disks, item); err != nil {
				return multierror.Prefix(err, "disk ->")
			}
		case "network":
			if err := parseNetworks(&devices.Networks, item); err != nil {
				return multierror.Prefix(err, "network ->")
			}
		case "scsi":
			if err := parseSCSIs(&devices.SCSIs, item); err != nil {
				return multierror.Prefix(err, "scsi ->")
			}
		}
	}

	*result = devices
	return nil
}

func parseDisks(result *[]*Disk, item *ast.ObjectItem) error {
	t := item.Keys[0].Token.Value().(string)
	n := item.Keys[1].Token.Value().(string)

	valid := []string{
		"size",
	}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	var disk Disk
	disk.DeviceName = n
	disk.DeviceType = t

	if err := mapstructure.WeakDecode(m, &disk); err != nil {
		return err
	}

	*result = append(*result, &disk)
	return nil
}

func parseNetworks(result *[]*Network, item *ast.ObjectItem) error {
	t := item.Keys[0].Token.Value().(string)
	n := item.Keys[1].Token.Value().(string)

	valid := []string{
		"build_vlan",
		"vlan",
		"switch_type",
	}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	var network Network
	network.DeviceName = n
	network.DeviceType = t

	if err := mapstructure.WeakDecode(m, &network); err != nil {
		return err
	}

	*result = append(*result, &network)
	return nil
}

func parseSCSIs(result *[]*SCSI, item *ast.ObjectItem) error {
	t := item.Keys[0].Token.Value().(string)
	n := item.Keys[1].Token.Value().(string)

	valid := []string{
		"type",
	}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	var scsi SCSI
	scsi.DeviceName = n
	scsi.DeviceType = t

	if err := mapstructure.WeakDecode(m, &scsi); err != nil {
		return err
	}

	*result = append(*result, &scsi)
	return nil
}

func checkHCLKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot check HCL keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key: %s", key))
		}
	}

	return result
}
