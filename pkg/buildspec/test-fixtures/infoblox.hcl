spec "default" {
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

    infoblox {
        subnet = "192.168.1.0/24"
        zone = "qa.local"
    }
}
