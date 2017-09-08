package setting

import (
	"../../manage/client"
	"encoding/json"

	"log"
)

//
func receiveParse(msg []byte) {

	var conf Conf
	json.Unmarshal(msg, &conf)

	if !validate(conf) {
		return
	}

	var conf_nil Conf
	if conf_nil != conf {
		log.Println("--------[Lived]", client.Lived, conf)

		kafkaHandle(conf)
		globalHandle(conf)
	}
}

//
func validate(conf Conf) bool {

	return true
}

var Restart bool = false

//
func globalHandle(conf Conf) {

	i := 0
	if "" != conf.Device {
		conf_global.Device = conf.Device
		i++
	}

	//
	if i > 0 {
		Restart = true
	}
}

//
func kafkaHandle(conf Conf) {

	hosts := conf.Kafka.Hosts
	if "" != hosts {
		conf_global.Kafka.Hosts = hosts
	}

	topic := conf.Kafka.Topic
	if "" != topic {
		conf_global.Kafka.Topic = topic
	}
}
