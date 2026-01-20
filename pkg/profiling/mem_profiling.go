package profiling

import (
	"io"
	"log"
	"runtime"
	"runtime/pprof"
	"time"
)

var memProfilingInterval = time.Second

var pprofWriteHeapProfile = func(w io.Writer) error {
	return pprof.WriteHeapProfile(w)
}

var DoMemProfiling = func(memProfFile string) (close func()) {

	writeMemProfile := func() {
		f, err := osCreate(memProfFile)
		if err != nil {
			log.Printf("could not create memory profile: %v", err)
			return
		}
		defer func() {
			_ = f.Close()
		}()
		runtime.GC() // get up-to-date statistics
		if err = pprofWriteHeapProfile(f); err != nil {
			log.Printf("could not write memory profile: %v", err)
		}
	}

	go func() {
		for {
			time.Sleep(memProfilingInterval)
			writeMemProfile()
		}
	}()

	return writeMemProfile
}
