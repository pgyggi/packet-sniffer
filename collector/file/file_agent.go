package file

type FileColl struct {
	running bool
}

var collector *FileColl = &FileColl{}

func NewCollector() *FileColl {
	return collector
}

func (collector *FileColl) IsRuning() bool {
	return false
}

func (collector *FileColl) Init(content []byte) error {

	//	//Load all agent config

	//	cnf, err := cfgfile.LoadByte(content)

	//	if err != nil {
	//		return fmt.Errorf("error loading packetbeat config: %v", err)
	//	}

	//	settings := settings.Settings{}

	//	err = cnf.Unpack(&settings)
	//	if err != nil {
	//		return fmt.Errorf("error init packetbeat config: %v", err)
	//	}

	//	//Init packet agent
	//	packet := New()
	//	packet.PbConfig = settings
	//	packet.HandleFlags()

	//	err = packet.Config()
	//	if err != nil {
	//		return err
	//	}

	//	defer packet.Cleanup()
	//	if err = packet.Setup(config.ConfigSingleton); err != nil {
	//		logp.Critical("Config packet sniffer error:%v", err)
	//		return err
	//	}

	//	collector.packet = packet

	return nil
}

func (collector *FileColl) Startup() error {
	return nil
}

func (collector *FileColl) Shutdown() error {
	return nil
}
