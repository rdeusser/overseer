# Overseer
[![GoDoc](https://godoc.org/github.com/iamthemuffinman/overseer?status.svg)](https://godoc.org/github.com/iamthemuffinman/overseer)
[![Build Status](https://travis-ci.org/iamthemuffinman/overseer.svg?branch=master)](https://travis-ci.org/iamthemuffinman/overseer) [![Go Report Card](https://goreportcard.com/badge/github.com/iamthemuffinman/overseer)](https://goreportcard.com/report/github.com/iamthemuffinman/overseer)

A provisioning toolkit.

Overseer uses something called a "hostspec" to determine how to build a physical or virtual machine.
Server names are read from the command line or another kind of spec called a "buildspec".

A hostspec looks like this:
```hcl
spec "indy.prod.kafka" {
    vsphere {
        cpus = 2
        cores = 1
        memory = 8096
        domain = "qa.local"
        cluster = "cluster01"
        datastore = "ds01"
        folder = "folder01"
        datacenter = "dc01"

        device "disk" "Hard disk 1" {
            size = 40
        }

        device "network" "Network adapter 1" {
            build_vlan = "dv-build"
            vlan = "dv-appservers"
            switch_type = "distributed"
        }

        device "scsi" "SCSI controller 1" {
            type = "paravirtual"
        }
    }

    foreman {
        hostgroup = "hg01"
        location = "location01"
        organization = "org01"
        environment = "env01"
        compute_profile = "compute01"
        architecture_id = 6
        compute_resource = "lol"
        domain_id = 6
        operating_system_id = 2
        partition_table_id = 6
        medium = "centos-7"
    }

    chef {
        base_role = "baserole01"
        run_list = [
            "role01",
            "role02"
        ]
    }
}
```

A buildspec for a physical host:
```hcl
hello.qa.local 1C:29:DF:E5:AA:B5
lol.qa.local 52:65:06:7A:C5:C8
with1234.qa.local 37:25:61:C8:B5:9C
nope.qa.local 19:62:AD:A7:92:BA
sometimes@#$@#%123135.qa.local E5:CF:60:13:C2:3E
```

A buildspec for a virtual host:
```hcl
hello.qa.local
lol.qa.local
with1234.qa.local
nope.qa.local
sometimes@#$@#%123135.qa.local
```

## Overseer kinda seems like Terraform?
Yeah, they do share some similarities. The hostspec concept was taken from how SaltStack uses profiles.
The one big difference and the reason I created this was because Terraform currently needs to maintain state.
Overseer does not and will never maintain state of any kind. The idea here is that you pass a list of hostnames
(or use a buildspec) and a hostspec that describes the kind of build you want and it'll go through and create
all of the necessary resources for you.
