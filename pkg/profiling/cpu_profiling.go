package profiling

import (
	"log"
	"runtime/pprof"
)

var pprofStartCPUProfile = pprof.StartCPUProfile

var DoCPUProfiling = func(cpuProfFile string) (close func()) {
	f, err := osCreate(cpuProfFile)
	if err != nil {
		log.Printf("could not create CPU profile: %v", err)
		return func() {}
	}
	close = func() {
		pprof.StopCPUProfile()
		_ = f.Close() // error handling omitted for brevity
	}
	if err = pprofStartCPUProfile(f); err != nil {
		log.Printf("could not start CPU profile: %v", err)
		return func() {}
	}
	return close
}
