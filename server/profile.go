package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/pprof"
)

func saveProfile(ptype string, f io.Writer, debug int) error {
	profile := pprof.Lookup(ptype)
	if profile == nil {
		return nil
	}

	if err := profile.WriteTo(f, debug); err != nil {
		return err
	}
	return nil
}

func startCPUProfile(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	pprof.StartCPUProfile(f)
	return nil
}

func stopCPUProfile() {
	pprof.StopCPUProfile()
}

func writeMemProfile(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	return pprof.WriteHeapProfile(f)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}

func startGoroutineProfile(port int) {
	go func() {
		http.HandleFunc("/", httpHandler)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
