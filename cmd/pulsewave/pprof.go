//go:build pprof
// +build pprof

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

// pprofPort override it at compilation with -ldflags "-X main.pprofPort=XXXX"
var pprofPort = "6060"

func init() {
	go func() {
		fmt.Printf("Starting pprof server, http://localhost:%s/debug/pprof/", pprofPort)
		err := http.ListenAndServe(fmt.Sprintf("localhost:%s", pprofPort), nil)
		if err != nil {
			panic(err)
		}
	}()
}
