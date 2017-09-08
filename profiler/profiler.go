package profiler

import (
	"flag"
	"log"
	"net/http"

	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	"../logp"
)

// cmdline flags
var memprofile, cpuprofile *bool
var httpprof *string
var cpuOut *os.File

func init() {
	memprofile = flag.Bool("memprofile", false, "Write memory profile to meminfo file")
	cpuprofile = flag.Bool("cpuprofile", false, "Write cpu profile cpuinfo file")
	httpprof = flag.String("httpprof", "", "Start pprof http server")
}

// BeforeRun takes care of necessary actions such as creating files
// before the beat should run.
func Run() {
	if *cpuprofile {
		cpuOut, err := os.Create("cpuinfo")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpuOut)
	}

	if *httpprof != "" {
		go func() {
			logp.Info("start pprof endpoint")
			logp.Info("finished pprof endpoint: %v", http.ListenAndServe(*httpprof, nil))
		}()
	}
}

// Cleanup handles cleaning up the runtime and OS environments. This includes
// tasks such as stopping the CPU profile if it is running.
func Cleanup() {
	if *cpuprofile {
		pprof.StopCPUProfile()
		cpuOut.Close()
	}

	if *memprofile {
		runtime.GC()
		writeHeapProfile("meminfo")
		debugMemStats()
	}
}

func debugMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logp.Debug("mem", "Memory stats: In use: %d Total (even if freed): %d System: %d",
		m.Alloc, m.TotalAlloc, m.Sys)
}

func writeHeapProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		logp.Err("Failed creating file %s: %s", filename, err)
		return
	}
	pprof.WriteHeapProfile(f)
	f.Close()

	logp.Info("Created memory profile file %s.", filename)
}
