package settings

import (
	"time"

	"../../../common"
	"../../../common/droppriv"
	"../../filter"
	"../procs"
)

type Settings struct {
	Modes      ModeConf
	Interfaces InterfacesConfig
	Flows      *Flows
	Protocols  map[string]*common.Config
	RunOptions droppriv.RunOptions
	Procs      procs.ProcsConfig
	Filter     []filter.FilterConfig
}

type ModeConf struct {
	Id     string
	Online bool //[default false]
	Host   string
}

type InterfacesConfig struct {
	Device         string
	Type           string
	File           string
	With_vlans     bool
	Bpf_filter     string
	Snaplen        int
	Buffer_size_mb int
	TopSpeed       bool
	Dumpfile       string
	OneAtATime     bool
	Loop           int
}

type Flows struct {
	Timeout string
	Period  string
}

type ProtocolCommon struct {
	Ports              []int         `config:"ports"`
	SendRequest        bool          `config:"send_request"`
	SendResponse       bool          `config:"send_response"`
	TransactionTimeout time.Duration `config:"transaction_timeout"`
}

// Config Singleton
var SettingSingleton Settings
