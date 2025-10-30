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
	l := logger()
	go func() {
		l.Info().
			Str("url", fmt.Sprintf("http://localhost:%s/debug/pprof/", pprofPort)).
			Msg("starting pprof server")

		err := http.ListenAndServe(fmt.Sprintf("localhost:%s", pprofPort), nil)
		if err != nil {
			panic(err)
		}
	}()
}
