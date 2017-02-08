package configspec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
)

type Spec struct {
	Foreman  Foreman  `mapstructure:"foreman"`
	Knife    Knife    `mapstructure:"knife"`
	Vsphere  Vsphere  `mapstructure:"vsphere"`
	Infoblox Infoblox `mapstructure:"infoblox"`
}

type Foreman struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Knife struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Vsphere struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Infoblox struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ParseFile parses the given configspec file.
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

	// Should be a list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: root should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"foreman",
		"knife",
		"vsphere",
		"infoblox",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	var spec Spec

	if o := list.Filter("foreman"); len(o.Items) > 0 {
		if err := parseForeman(&spec.Foreman, o); err != nil {
			return nil, fmt.Errorf("error parsing foreman block: %s", err)
		}
	}
	if o := list.Filter("knife"); len(o.Items) > 0 {
		if err := parseKnife(&spec.Knife, o); err != nil {
			return nil, fmt.Errorf("error parsing knife block: %s", err)
		}
	}
	if o := list.Filter("vsphere"); len(o.Items) > 0 {
		if err := parseVsphere(&spec.Vsphere, o); err != nil {
			return nil, fmt.Errorf("error parsing vsphere block: %s", err)
		}
	}
	if o := list.Filter("infoblox"); len(o.Items) > 0 {
		if err := parseInfoblox(&spec.Infoblox, o); err != nil {
			return nil, fmt.Errorf("error parsing infoblox block: %s", err)
		}
	}

	return &spec, nil
}

func parseForeman(result *Foreman, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("Only one %q block allowed", "foreman")
	}

	// Get our "foreman" object
	o := list.Items[0]

	valid := []string{
		"username",
		"password",
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

func parseKnife(result *Knife, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("Only one %q block allowed", "knife")
	}

	// Get our "knife" object
	o := list.Items[0]

	valid := []string{
		"username",
		"password",
	}
	if err := checkHCLKeys(o.Val, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var knife Knife
	if err := mapstructure.WeakDecode(m, &knife); err != nil {
		return err
	}

	*result = knife
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
		"username",
		"password",
	}
	if err := checkHCLKeys(listVal, valid); err != nil {
		return multierror.Prefix(err, "vsphere ->")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var vsphere Vsphere
	if err := mapstructure.WeakDecode(m, &vsphere); err != nil {
		return err
	}

	*result = vsphere
	return nil
}

func parseInfoblox(result *Infoblox, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block allowed", "infoblox")
	}

	// Get our infoblox object
	o := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := o.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	valid := []string{
		"username",
		"password",
	}
	if err := checkHCLKeys(listVal, valid); err != nil {
		return multierror.Prefix(err, "infoblox ->")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, o.Val); err != nil {
		return err
	}

	var infoblox Infoblox
	if err := mapstructure.WeakDecode(m, &infoblox); err != nil {
		return err
	}

	*result = infoblox
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
