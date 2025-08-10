package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"
)

const (
    PROCNAME = "altervoice"
)

type SerMgt struct{}

type Startarg struct {
    Addr string
    Port string
    DebugFlag int
}

var instance_process  = Startarg{
    Addr: "0.0.0.0",
    Port: "9000",
    DebugFlag: 0,
}

func IsFileExits(filename string) bool {
    if _,err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }

    return true
}

func Looptask() (err error) {
    router := newRouter()
    
    srv := &http.Server{
        Addr: instance_process.Addr + ":" + instance_process.Port,
        Handler: router,
    }
    err = srv.ListenAndServe()

    if err != nil {
        return err
    }

    return nil
}

func (srv *SerMgt) Start(d *dae) (err error){
    if err = InitLog(d); err != nil {
        LogPrint(LEVEL_FATAL,"[initlog] failed\n")
        return err
    }

    if d == nil {
        LogPrint(LEVEL_INFO,"the %s will start in the frontground\n")
    }else {

        if err = d.daemond(); err != nil {
            LogPrint(LEVEL_FATAL,"[daemond] failed:%s\n",err.Error())
            return err
        }
    }


    if err = Looptask(); err != nil {
        LogPrint(LEVEL_FATAL,"[looptask] failed:%s\n",err.Error())
        return err
    }

    return nil

}

func (srv *SerMgt) Stop(d *dae) (err error){
    if d == nil {
        fmt.Printf("the %s is running in the frontground",PROCNAME)
        return nil
    }

    fp,err := os.Open(d.pidfile)

    if err != nil {
        LogPrint(LEVEL_ERROR,"[stop] failed:%s\n",err.Error())
        return err
    }

    var pid int

    if _,err := fmt.Fscanf(fp,"%d",&pid); err != nil {
        LogPrint(LEVEL_ERROR,"[stop] failed:%s",err.Error())
        return err
    }

    for{
        if err = syscall.Kill(pid,syscall.SIGKILL); err == syscall.ENOENT{
            break
        }
        time.Sleep(1*time.Second)
    }

    return nil
}

func (srv *SerMgt) Status(d *dae) (err error) {
    if d == nil {
        fmt.Printf("the %s is not running the background,so you have to check the status by yourself\n",PROCNAME)
        return nil
    }

    var pid string
    if IsFileExits(d.pidfile) {
        fp,err := os.Open(d.pidfile)

        if err != nil {
            fmt.Fprintf(os.Stdin,"\033[1;5;31m%s Status\033[0m:\ncheck status failed:%s\n",PROCNAME,err.Error())
            return err
        }
        
        fmt.Fscanf(fp,"%s",&pid)

        var procname string

        fp1,err1 := os.Open("/proc/"+pid+"comm")

        if err1 == nil {
            fmt.Fscanf(fp1,"%s",&procname)
            if procname == PROCNAME {
                fmt.Fprintf(os.Stdout,"%s Status:\n\tthe status:running\n")
            }

            fmt.Fprintf(os.Stdout,"\033[1;5;31m%s Status\033[0m:\n%s is already exits\n",PROCNAME,PROCNAME)
        }

        fmt.Fprintf(os.Stdout,"\033[1;5;31m%s Status\033[0m:\n%s is already exits\n",PROCNAME,PROCNAME)
        return err1 
    }

    fmt.Fprintf(os.Stdout,"%s Status:\n\tthe status:stopped\n")
    
    return nil
}

func (srv *SerMgt) Restart(d *dae) (err error){
    if d == nil {
        fmt.Printf("the %s is not running background,so you can mgt by yourself\n",PROCNAME)
        return nil
    }
    if err = srv.Stop(d); err != nil {
        return err
    }

    if err = srv.Start(d); err != nil {
        return err

    }

    return nil
}

func (srv *SerMgt) Usage() {
    fmt.Fprintf(os.Stdout,"\033[1;5;31musage\033[0m:\n%s [start | stop | restart | status]\n",PROCNAME)
    os.Exit(1)
}

