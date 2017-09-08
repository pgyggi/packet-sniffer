package collector

import (
	"io/ioutil"

	"../logp"
	//"./file"
	"../manage/setting"
	"./packet"
	"time"

	"log"

	"net"
	//_ "net/http/pprof"
	//"net/http"
)

var Name = "Agent"

func Init(conf string) error {
	loadConfigFromFile(conf)
	return nil
}

func loadConfigFromFile(conf string,reline ...bool) {

	//Load packet capture configuration
	if content, err := ioutil.ReadFile(conf); err != nil {
		logp.Err("Load packet config error:%v", err)
	} else {
		collector := packet.NewCollector()

		go restart(collector,conf)

		if err := collector.Init(content); err != nil {
			logp.Err("Init packet collector error:%v", err)
			return
		} else if err := collector.Startup(); err != nil {
			logp.Err("Run packet collector error:%v", err)
			return
		}
	}

	//Load file capture configuration
	//	if content, err := ioutil.ReadFile("file.yml"); err != nil {
	//		logp.Err("Load file config error:%v", err)
	//	} else {
	//		collector := file.NewCollector()
	//		if err := collector.Init(content); err != nil {
	//			logp.Err("Init file collector error:%v", err)
	//		} else if err := collector.Startup(); err != nil {
	//			logp.Err("Run file collector error:%v", err)
	//		}
	//	}

	//
	if 0 == len(reline) || !reline[0] {
		setting.InitOnline()
	}

	/*
		//
		go func() {
			log.Println(http.ListenAndServe(get_internal()+":3999", nil))
		}()
	*/
}

//
func restart(p *packet.PacketColl,conf string) {

	for {
		time.Sleep(2e9)
		if setting.Restart {
			logp.Info(" Restart..................")
			setting.Restart = false
			/*
				p.Shutdown()
				err:=p.Startup()
				if nil!=err{
					logp.Err(" Restart Err..................",err)
				}
				logp.Info(" Restart IsRuning.................",p.IsRuning())
			*/
			p.Shutdown()
			loadConfigFromFile(conf,true)
			return
		}
	}
}

//
func get_internal() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
