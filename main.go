package main

import (
	"flag"
	"os"
)


func main() {
    var daemonMode bool
    var d *dae = nil

    flag.BoolVar(&daemonMode,"daemon",false,"the service will run in the background")
    flag.Parse()

    if daemonMode {
        d.pidfile = PIDFILE
        d.logfile = ""
        d.curdir = "/"
    }

    srv := &SerMgt{}

    if len(os.Args) == 1 {
        srv.Usage()
    }else if os.Args[1] == "start" {
        if err := srv.Start(d); err != nil {
            os.Exit(1)
        }         
    }else if os.Args[1] == "stop" {
        if err := srv.Stop(d); err != nil {
            os.Exit(1)
        }
        os.Exit(0)
    }else if os.Args[1] == "status" {
        if err := srv.Status(d); err != nil {
            os.Exit(1)
        }
        os.Exit(0)
    }else if os.Args[1] == "restart" {
        if err := srv.Restart(d); err != nil {
            os.Exit(1)
        }
    }else {
        srv.Usage()
    }
}
