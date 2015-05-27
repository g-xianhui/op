package main

import (
	"io"
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
