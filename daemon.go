package main

import (
	"log"
	"github.com/sevlyar/go-daemon"
)


const (
    PIDFILE = "/var/run/altervoice.pid"
   // LOGFILE = "/var//log/altervoice.log"
)

type dae struct {
    pidfile string 
    logfile string
    curdir string
}


func (d *dae) daemond() error {
    cntxt :=  &daemon.Context{
        PidFileName: d.pidfile,
        PidFilePerm: 0644,
        LogFileName: "",
        LogFilePerm: 0,
        WorkDir: d.curdir,
        Umask: 027,
    }

    _, err := cntxt.Reborn()

    defer cntxt.Release()

    if err != nil {
        log.Printf("Reborn() failed\n")
        return err
    }

    return nil

}

