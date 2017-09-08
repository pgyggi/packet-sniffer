package common

import (
	"fmt"
	"os"

	"../logp"
	"../profiler"
)

const (
	AGENT  = "agent"
	SERVER = "server"
	LB     = "lb"
)

const (
	DATA_SOURCE string = "di_datasource"
	RIIL        string = "RIIL"
	NTA         string = "NTA"
	LBS         string = "LBS"
	Packet      string = "Packet"
)

var Role string

func ExitProcess(err error) {
	profiler.Cleanup()
	code := 1
	if err != nil && code != 0 {
		// logp may not be initialized so log the err to stderr too.
		logp.Critical("Exiting: %v", err)
		fmt.Fprintf(os.Stderr, "Exiting: %v\n", err)
	}
	os.Exit(code)
}
