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

        device "disk" "Hard disk 1" {
            size = 200
        }

        device "network" "Network adapter 1" {
            build_vlan = "dv-build"
            vlan = "dv-appservers"
            switch_type = "distributed"
        }

        device "network" "Network adapter 1" {
            vlan = "dv-appservers"
            switch_type = "distributed"
        }

        device "scsi" "SCSI controller 1" {
            type = "paravirtual"
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
        run_list = [
            "role01",
            "role02"
        ]
    }
}
