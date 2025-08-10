package main

import(
    "net/http"
    "github.com/gorilla/mux"
    "time"
)

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var gRoutes = Routes{
    Route{
        Name:"alter_voice",
        Method:"POST",
        Pattern:"/alter_voice",
        HandlerFunc:AlterVoice,
    },
}

func handlerWithLog(inner http.Handler,name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w,r)

        if instance_process.DebugFlag == 1 {
            LogPrint(LEVEL_INFO,
                "%s\t%s\t%s\t%s\n",
                r.Method,
                r.RequestURI,
                name,
                time.Since(start),
            )
        }

        
    })
}

func newRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)

    for _,route := range gRoutes {
        handler := http.TimeoutHandler(handlerWithLog(route.HandlerFunc,route.Name),time.Second*4,"MONITOR_PROCESS TIMEOUT")

        router.Methods(route.Method).
        Path(route.Pattern).
        Name(route.Name).
        Handler(handler)
    }

    return router
}
