package packet

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/tsg/gopacket/layers"

	"../../common/droppriv"
	"../../common/errors"
	"../../config"
	"../../logp"
	"../filter"
	"../publisher"
	"./decoder"
	"./flows"
	"./procs"
	"./protos"
	"./protos/icmp"
	"./protos/tcp"
	"./protos/udp"
	"./publish"
	"./settings"
	"./sniffer"

	_ "./protos/http"
)

type Packet struct {
	Pub         *publish.PacketbeatPublisher
	PbConfig    settings.Settings
	CmdLineArgs CmdLineArgs
	Sniff       *sniffer.SnifferSetup
	filters     *filter.FilterList // Filters
	services    []interface {
		Start()
		Stop()
	}
}

type CmdLineArgs struct {
	File         *string
	Loop         *int
	OneAtAtime   *bool
	TopSpeed     *bool
	Dumpfile     *string
	PrintDevices *bool
	WaitShutdown *int
}

var cmdLineArgs CmdLineArgs

const (
	defaultQueueSize     = 2048
	defaultBulkQueueSize = 0
)

func init() {
	cmdLineArgs = CmdLineArgs{
		File:         flag.String("I", "", "Read packet data from specified file"),
		Loop:         flag.Int("l", 1, "Loop file. 0 - loop forever"),
		OneAtAtime:   flag.Bool("O", false, "Read packets one at a time (press Enter)"),
		TopSpeed:     flag.Bool("t", false, "Read packets as fast as possible, without sleeping"),
		Dumpfile:     flag.String("dump", "", "Write all captured packets to this libpcap file"),
		PrintDevices: flag.Bool("devices", false, "Print the list of devices and exit"),
		WaitShutdown: flag.Int("waitstop", 0, "Additional seconds to wait before shutting down"),
	}
}

func New() *Packet {

	pb := &Packet{}
	pb.CmdLineArgs = cmdLineArgs

	return pb
}

func (pb *Packet) HandleFlags() error {
	// -devices CLI flag
	if *pb.CmdLineArgs.PrintDevices {
		devs, err := sniffer.ListDeviceNames(true)
		if err != nil {
			return fmt.Errorf("Error getting devices list: %v\n", err)
		}
		if len(devs) == 0 {
			fmt.Printf("No devices found.")
			if runtime.GOOS != "windows" {
				fmt.Printf(" You might need sudo?\n")
			} else {
				fmt.Printf("\n")
			}
		}
		for i, dev := range devs {
			fmt.Printf("%d: %s\n", i, dev)
		}
		return errors.GracefulExit
	}
	return nil
}

// Loads the beat specific config and overwrites params based on cmd line
func (pb *Packet) Config() error {

	// CLI flags over-riding config
	if *pb.CmdLineArgs.TopSpeed {
		pb.PbConfig.Interfaces.TopSpeed = true
	}

	if len(*pb.CmdLineArgs.File) > 0 {
		pb.PbConfig.Interfaces.File = *pb.CmdLineArgs.File
	}

	pb.PbConfig.Interfaces.Loop = *pb.CmdLineArgs.Loop
	pb.PbConfig.Interfaces.OneAtATime = *pb.CmdLineArgs.OneAtAtime

	if len(*pb.CmdLineArgs.Dumpfile) > 0 {
		pb.PbConfig.Interfaces.Dumpfile = *pb.CmdLineArgs.Dumpfile
	}

	filters, err := filter.New(pb.PbConfig.Filter)
	if err != nil {
		return fmt.Errorf("error initializing filters: %v", err)
	} else {
		pb.filters = filters
	}

	// assign global singleton as it is used in protocols
	// TODO: Refactor
	settings.SettingSingleton = pb.PbConfig

	return nil
}

// Setup packetbeat
func (pb *Packet) Setup(config config.Config) error {

	if err := procs.ProcWatcher.Init(pb.PbConfig.Procs); err != nil {
		logp.Critical(err.Error())
		return err
	}

	queueSize := defaultQueueSize
	if config.Shipper.QueueSize != nil {
		queueSize = *config.Shipper.QueueSize
	}
	bulkQueueSize := defaultBulkQueueSize
	if config.Shipper.BulkQueueSize != nil {
		bulkQueueSize = *config.Shipper.BulkQueueSize
	}

	publisher, err := publisher.New("Packet", config.Output, config.Shipper)
	if err != nil {
		return fmt.Errorf("error initializing publisher: %v", err)
	}

	publisher.RegisterFilter(pb.filters)
	pb.Pub = publish.NewPublisher(publisher, queueSize, bulkQueueSize)
	pb.Pub.Start()

	logp.Info("main", "Initializing protocol plugins")
	if err := protos.Protos.Init(false, pb.Pub, pb.PbConfig.Protocols); err != nil {
		return fmt.Errorf("Initializing protocol analyzers failed: %v", err)
	}

	logp.Info("main", "Initializing sniffer")
	if err := pb.setupSniffer(); err != nil {
		return fmt.Errorf("Initializing sniffer failed: %v", err)
	}

	// This needs to be after the sniffer Init but before the sniffer Run.
	if err := droppriv.DropPrivileges(settings.SettingSingleton.RunOptions); err != nil {
		return err
	}
	return nil
}

func (pb *Packet) setupSniffer() error {
	cfg := &pb.PbConfig

	withVlans := cfg.Interfaces.With_vlans
	_, withICMP := cfg.Protocols["icmp"]
	filter := cfg.Interfaces.Bpf_filter
	if filter == "" && cfg.Flows == nil {
		filter = protos.Protos.BpfFilter(withVlans, withICMP)
	}

	pb.Sniff = &sniffer.SnifferSetup{}
	return pb.Sniff.Init(false, pb.makeWorkerFactory(filter))
}

func (pb *Packet) makeWorkerFactory(filter string) sniffer.WorkerFactory {
	return func(dl layers.LinkType) (sniffer.Worker, string, error) {
		var f *flows.Flows
		var err error

		if pb.PbConfig.Flows != nil {
			f, err = flows.NewFlows(pb.Pub, pb.PbConfig.Flows)
			if err != nil {
				return nil, "", err
			}
		}

		var icmp4 icmp.ICMPv4Processor
		var icmp6 icmp.ICMPv6Processor
		if cfg, exists := pb.PbConfig.Protocols["icmp"]; exists {
			icmp, err := icmp.New(false, pb.Pub, cfg)
			if err != nil {
				return nil, "", err
			}

			icmp4 = icmp
			icmp6 = icmp
		}

		tcp, err := tcp.NewTcp(&protos.Protos)
		if err != nil {
			return nil, "", err
		}

		udp, err := udp.NewUdp(&protos.Protos)
		if err != nil {
			return nil, "", err
		}

		worker, err := decoder.NewDecoder(f, dl, icmp4, icmp6, tcp, udp)
		if err != nil {
			return nil, "", err
		}

		if f != nil {
			pb.services = append(pb.services, f)
		}
		return worker, filter, nil
	}
}

func (pb *Packet) Run() error {

	// start services
	for _, service := range pb.services {
		service.Start()
	}

	var wg sync.WaitGroup
	errC := make(chan error, 1)

	// Run the sniffer in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := pb.Sniff.Run()
		if err != nil {
			errC <- fmt.Errorf("Sniffer main loop failed: %v", err)
		}
	}()

	logp.Debug("main", "Waiting for the sniffer to finish")
	wg.Wait()
	select {
	default:
	case err := <-errC:
		return err
	}

	// kill services
	for _, service := range pb.services {
		service.Stop()
	}

	waitShutdown := pb.CmdLineArgs.WaitShutdown
	if waitShutdown != nil && *waitShutdown > 0 {
		time.Sleep(time.Duration(*waitShutdown) * time.Second)
	}

	return nil
}

func (pb *Packet) Cleanup() error {

	// TODO:
	// pb.TransPub.Stop()

	return nil
}

// Called by the Beat stop function
func (pb *Packet) Stop() {
	pb.Sniff.Stop()
}
