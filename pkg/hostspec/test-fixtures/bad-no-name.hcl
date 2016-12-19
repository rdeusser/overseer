spec {
    virtual {
        cpus = 2
        cores = 1
        memory = 8096
    }

    vsphere {
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
}
