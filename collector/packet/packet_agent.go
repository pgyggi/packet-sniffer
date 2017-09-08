package packet

import (
	"fmt"

	"../../cfgfile"
	"../../config"
	"../../logp"
	"../../manage/setting"
	"./settings"
)

type PacketColl struct {
	packet  *Packet
	running bool
}

var collector *PacketColl = &PacketColl{}

func NewCollector() *PacketColl {
	return collector
}

func (collector *PacketColl) IsRuning() bool {
	return collector.running
}

func (collector *PacketColl) Init(content []byte) error {

	//Load all agent config

	cnf, err := cfgfile.LoadByte(content)

	if err != nil {
		return fmt.Errorf("error loading packetbeat config: %v", err)
	}

	settings := settings.Settings{}

	err = cnf.Unpack(&settings)
	if err != nil {
		return fmt.Errorf("error init packetbeat config: %v", err)
	}

	setting.ResetConf(&settings)

	//Init packet agent
	packet := New()
	packet.PbConfig = settings

	//packet.HandleFlags()

	err = packet.Config()
	if err != nil {
		return err
	}
	//collect all packet.yml and setup packet.
	defer packet.Cleanup()
	if err = packet.Setup(config.ConfigSingleton); err != nil {
		logp.Critical("Config packet sniffer error:%v", err)
		return err
	}

	collector.packet = packet

	return nil
}

func (collector *PacketColl) Startup() error {

	var err error
	go func() {
		collector.running = true
		if err = collector.packet.Run(); err != nil {
			collector.running = false
			logp.Critical("Init packet sniffer error:%v", err)
		}
	}()

	return err
}

func (collector *PacketColl) Shutdown() error {
	collector.packet.Stop()
	collector.running = false
	return nil
}
