package setting

import (
	"../../collector/packet/settings"
	"../../manage/client"
	"github.com/Shopify/sarama"
	"strings"
)

//
func InitOnline() {

	host := settings.SettingSingleton.Modes.Host
	if "" != host {
		settings.SettingSingleton.Modes.Online = true
		//c:=client.NewClient("172.16.2.8:8081", "/agentmanager/test",0,0)
		//c:=client.NewClient("172.16.2.9:8080", "/whecho",2e9,0)
		//return
		c := client.NewClient(host, "/whecho", 2e9, 0)
		c.SetReadHandle(receiveParse)
		c.InitWSConn()
		c.SendMsg("{\"type\":\"ping..\"}")
		c.SendMsg("{\"Device\":\"121212121\",\"Kafka\":{\"Topic\":\"wh22\",\"Hosts\":\"127.0.0.1,192.168.0.1\"}}")
		//
	}
}

//
var conf_global Conf

//
func LoadKafkaMessage(msg *sarama.ProducerMessage) {

	topic := conf_global.Kafka.Topic
	if "" != topic {
		msg.Topic = topic
	}
}

//
func LoadKafkaSetting() []string {

	var hs []string
	hosts := conf_global.Kafka.Hosts
	if "" != hosts {
		hs = strings.Split(hosts, ",")

	}

	return hs
}

func ResetConf(setting *settings.Settings) {

	device := conf_global.Device
	if "" != device {
		setting.Interfaces.Device = device
	}
}
