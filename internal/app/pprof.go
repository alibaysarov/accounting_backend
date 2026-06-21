package app

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func startPprofServer(addr string) {
	log.Println("starting pprof on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("pprof server error:", err)
	}
}
