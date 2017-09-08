package config

import (
	"time"

	"../cfgfile"
	"../common"
	"../logp"
)

var ConfigSingleton Config

type Config struct {
	Logging logp.Logging `yaml:"logging"`
	Output  map[string]*common.Config
	Shipper ShipperConfig
}

type ShipperConfig struct {
	common.EventMetadata `config:",inline"` // Fields and tags to add to each event.
	Name                 string
	RefreshTopologyFreq  time.Duration `config:"refresh_topology_freq"`
	Ignore_outgoing      bool          `config:"ignore_outgoing"`
	Topology_expire      int           `config:"topology_expire"`
	Geoip                common.Geoip  `config:"geoip"`

	// internal publisher queue sizes
	QueueSize     *int `config:"queue_size"`
	BulkQueueSize *int `config:"bulk_queue_size"`
	MaxProcs      *int `config:"max_procs"`
}

func Init(conf string) error {
	cnf := conf
	if cfg, err := cfgfile.Load(cnf); err != nil {
		return err
	} else {
		cfg.Unpack(&ConfigSingleton)
	}
	return nil
}
