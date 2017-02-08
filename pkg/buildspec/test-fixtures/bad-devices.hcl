spec "default" {
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
    }
}
