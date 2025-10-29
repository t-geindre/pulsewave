//go:build pprof
// +build pprof

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		fmt.Println("starting pprof server, http://localhost:6060/debug/pprof/")
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			panic(err)
		}
	}()
}
