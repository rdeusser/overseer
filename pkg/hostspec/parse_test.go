package hostspec

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		File     string
		Expected *Spec
		Err      bool
	}{
		{
			"basic.hcl",
			&Spec{
				Name: "default",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
					Devices: Devices{
						Disks: []*Disk{
							&Disk{
								DeviceName: "Hard disk 1",
								DeviceType: "disk",
								Size:       40,
							},
						},
						Networks: []*Network{
							&Network{
								DeviceName: "Network adapter 1",
								DeviceType: "network",
								BuildVLAN:  "dv-build",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
						},
						SCSIs: []*SCSI{
							&SCSI{
								DeviceName: "SCSI controller 1",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
						},
					},
				},
			},
			false,
		},
		{
			"bad-no-name.hcl",
			nil,
			true,
		},
		{
			"bad-devices.hcl",
			&Spec{
				Name: "default",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
				},
			},
			false,
		},
		{
			"bad-host-options.hcl",
			&Spec{
				Name: "default",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
				},
			},
			false,
		},
		{
			"bad-device-name.hcl",
			nil,
			true,
		},
		{
			"foreman.hcl",
			&Spec{
				Name: "default",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
					Devices: Devices{
						Disks: []*Disk{
							&Disk{
								DeviceName: "Hard disk 1",
								DeviceType: "disk",
								Size:       40,
							},
						},
						Networks: []*Network{
							&Network{
								DeviceName: "Network adapter 1",
								DeviceType: "network",
								BuildVLAN:  "dv-build",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
						},
						SCSIs: []*SCSI{
							&SCSI{
								DeviceName: "SCSI controller 1",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
						},
					},
				},
				Foreman: Foreman{
					Hostgroup:         "hg01",
					Location:          "location01",
					Organization:      "org01",
					Environment:       "env01",
					ComputeProfile:    "compute01",
					ArchitectureID:    6,
					ComputeResource:   "lol",
					DomainID:          6,
					OperatingSystemID: 2,
					PartitionTableID:  6,
					Medium:            "centos-7",
				},
			},
			false,
		},
		{
			"chef.hcl",
			&Spec{
				Name: "default",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
					Devices: Devices{
						Disks: []*Disk{
							&Disk{
								DeviceName: "Hard disk 1",
								DeviceType: "disk",
								Size:       40,
							},
						},
						Networks: []*Network{
							&Network{
								DeviceName: "Network adapter 1",
								DeviceType: "network",
								BuildVLAN:  "dv-build",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
						},
						SCSIs: []*SCSI{
							&SCSI{
								DeviceName: "SCSI controller 1",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
						},
					},
				},
				Chef: Chef{
					BaseRole: "baserole01",
					RunList: []string{
						"role01",
						"role02",
					},
				},
			},
			false,
		},
		{
			"complete.hcl",
			&Spec{
				Name: "indy.prod.kafka",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
					Devices: Devices{
						Disks: []*Disk{
							&Disk{
								DeviceName: "Hard disk 1",
								DeviceType: "disk",
								Size:       40,
							},
						},
						Networks: []*Network{
							&Network{
								DeviceName: "Network adapter 1",
								DeviceType: "network",
								BuildVLAN:  "dv-build",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
						},
						SCSIs: []*SCSI{
							&SCSI{
								DeviceName: "SCSI controller 1",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
						},
					},
				},
				Foreman: Foreman{
					Hostgroup:         "hg01",
					Location:          "location01",
					Organization:      "org01",
					Environment:       "env01",
					ComputeProfile:    "compute01",
					ArchitectureID:    6,
					ComputeResource:   "lol",
					DomainID:          6,
					OperatingSystemID: 2,
					PartitionTableID:  6,
					Medium:            "centos-7",
				},
				Chef: Chef{
					BaseRole: "baserole01",
					RunList: []string{
						"role01",
						"role02",
					},
				},
			},
			false,
		},
		{
			"multiple-devices.hcl",
			&Spec{
				Name: "indy.prod.kafka",
				Vsphere: Vsphere{
					Domain:     "qa.local",
					Cluster:    "cluster01",
					Datastore:  "ds01",
					Folder:     "folder01",
					Datacenter: "dc01",
					CPUs:       2,
					Cores:      1,
					Memory:     8096,
					Devices: Devices{
						Disks: []*Disk{
							&Disk{
								DeviceName: "Hard disk 1",
								DeviceType: "disk",
								Size:       40,
							},
							&Disk{
								DeviceName: "Hard disk 2",
								DeviceType: "disk",
								Size:       200,
							},
						},
						Networks: []*Network{
							&Network{
								DeviceName: "Network adapter 1",
								DeviceType: "network",
								BuildVLAN:  "dv-build",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
							&Network{
								DeviceName: "Network adapter 2",
								DeviceType: "network",
								BuildVLAN:  "",
								VLAN:       "dv-appservers",
								SwitchType: "distributed",
							},
						},
						SCSIs: []*SCSI{
							&SCSI{
								DeviceName: "SCSI controller 1",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
							&SCSI{
								DeviceName: "SCSI controller 2",
								DeviceType: "scsi",
								Type:       "paravirtual",
							},
						},
					},
				},
				Foreman: Foreman{
					Hostgroup:         "hg01",
					Location:          "location01",
					Organization:      "org01",
					Environment:       "env01",
					ComputeProfile:    "compute01",
					ArchitectureID:    6,
					ComputeResource:   "lol",
					DomainID:          6,
					OperatingSystemID: 2,
					PartitionTableID:  6,
					Medium:            "centos-7",
				},
				Chef: Chef{
					BaseRole: "baserole01",
					RunList: []string{
						"role01",
						"role02",
					},
				},
			},
			false,
		},
		{
			"bad-same-device-name.hcl",
			nil,
			true,
		},
	}

	for _, tt := range cases {
		t.Logf("Testing parse on: %s", tt.File)

		path, err := filepath.Abs(filepath.Join("./test-fixtures", tt.File))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		actual, err := ParseFile(path)
		if (err != nil) != tt.Err {
			t.Fatalf("file: %s\n\n%s", tt.File, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.Expected) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tt.File, actual, tt.Expected)
		}
	}
}
