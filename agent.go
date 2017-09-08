package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"./collector"
	"./common"
	"./config"
	"./logp"
	"./profiler"
	//"./manage/setting"
	"flag"
)

const (
	Name       string = "Agent"
	PROCESS_ID string = "process.pid"
)

var code = 0

var exitChan = make(chan int)
var (
	packet                = flag.String("p", "/etc/packet-sniffer/packet.yml", "conf file for packet.yml")
	conf         = flag.String("conf", "/etc/packet-sniffer/config.yml", "conf file for config.yml")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	//load config.yml to system
	if err := config.Init(*conf); err != nil {
		fmt.Printf("Load config file error: %+v\n", err)
		common.ExitProcess(err)
	}
	err := logp.Init(Name, &config.ConfigSingleton.Logging)

	if err != nil {
		fmt.Printf("Init system error: %+v\n", err)
		common.ExitProcess(err)
	}

	//logp.SetStderr()

	profiler.Run()

	writePIDFile()
	//listen system call os.Interrupt,os.Kill,syscall.SIGTERM
	go listenToSystemSignels()

	if err := initAgent(); err != nil {
		logp.Critical("Init agent error: %v", err)
	}

	//setting.InitOnline()

	<-exitChan

}

func listenToSystemSignels() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)
	signal.Notify(signalChan, syscall.SIGTERM)
	sig := <-signalChan
	logp.Info("Shutdown server.", sig)
	exitChan <- 1
}

func writePIDFile() {
	pidFile := PROCESS_ID
	err := os.MkdirAll(filepath.Dir(pidFile), 0700)
	if err != nil {
		logp.Err("Failed to verify pid directory: %v", err)
	}
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		logp.Err("Failed to write pidfile: %v", err)
	}
}

//start to collect data and load packet.yml to system

func initAgent() error {
	return collector.Init(*packet)
}
