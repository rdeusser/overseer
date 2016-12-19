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

        device "disk" {
            size = 40
        }

        device "network" {
            build_vlan = "dv-build"
            vlan = "dv-appservers"
            switch_type = "distributed"
        }

        device "scsi" {
            type = "paravirtual"
        }
    }
}
