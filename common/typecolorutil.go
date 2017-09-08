package common

import (
	"encoding/json"
	"github.com/lucasb-eyer/go-colorful"
	"io/ioutil"
	"log"
	"sync"
)

var Typecolor map[string]interface{}
var TypecolorInit = false

func ReadTCFile() error {
	filecontent, err := ioutil.ReadFile("typecolors.yml")
	if err == nil {
		//struct 到json str
		err := json.Unmarshal(filecontent, &Typecolor)
		if err == nil {
			return nil
		}
		return err
	}
	return err
}

var lock sync.Mutex

func WriteTCFile() error {
	lock.Lock()
	content, err := json.Marshal(Typecolor)
	if err == nil {
		ioutil.WriteFile("typecolors.yml", content, 0666)
	}
	log.Println("WriteTCFile error:", err)
	lock.Unlock()
	return err
}

func CheckTypeColor(sslice []interface{}) error {
	if !TypecolorInit {
		log.Println("CheckTypeColor init start!")
		err := ReadTCFile()
		if err != nil {
			log.Println("typecolor init exp:", err)
		} else {
			TypecolorInit = true
		}
		log.Println("CheckTypeColor int end!", TypecolorInit)
	}
	tmap := Typecolor["standard"].(map[string]interface{})
	tslice := Typecolor["standby"].([]interface{})
	flag := false
	for _, alarmType := range sslice {
		if tmap[alarmType.(string)] == nil {
			if len(tslice) > 0 {
				tmap[alarmType.(string)] = tslice[0]
				if len(tslice) == 1 {
					tslice = tslice[:0]
				} else {
					tslice = tslice[:1]
				}
			} else {
				tempcolor := colorful.WarmColor().Hex()
				tmap[alarmType.(string)] = tempcolor
			}
			flag = true
		}
	}
	if flag {
		Typecolor["standby"] = tslice
		WriteTCFile()
	}
	return nil
}
