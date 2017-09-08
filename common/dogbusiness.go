package common

import (
	"encoding/json"
	//	"fmt"
	//	"io/ioutil"
	//	"log"
	//	"time"
	"sync"
)

//资源ip,port,type信息
var ResourceSet map[string]bool

var resLock sync.Mutex

func SetResourceSet(key string) {
	resLock.Lock()
	if ResourceSet == nil {
		ResourceSet = make(map[string]bool)
	}
	ResourceSet[key] = true
	resLock.Unlock()
}

func GetResourceSet() map[string]bool {
	return ResourceSet
}

//{
//  "user1": {
//    "businessName1": {
//      "key1": {
//        "color": "yellow",
//        "data": [
//          {
//            "position": "p1",
//            "infos": {
//              "ip": "192.168.0.1",
//              "port": "80",
//              "type": "mysql"
//            }
//          },
//          {
//            "position": "p1",
//            "infos": {
//              "ip": "192.168.0.1",
//              "port": "80",
//              "type": "mysql"
//            }
//          }
//        ]
//      }
//    }
//  }
//}
//业务应用信息		  - userId   -servicename - infos
var BusinessMap map[string]map[string]interface{}

var bsLock sync.Mutex

func GetBusinessMap(userId string, businessId string) interface{} {
	return BusinessMap[userId][businessId]
}

func GetBusinessMaps() map[string]map[string]interface{} {
	return BusinessMap
}
func GetAllBusinessMap(userId string) map[string]interface{} {
	return BusinessMap[userId]
}

func SetBusinessMap(userId string, businessId string, infos interface{}) {
	bsLock.Lock()
	if BusinessMap == nil {
		BusinessMap = make(map[string]map[string]interface{})
	}
	if BusinessMap[userId] == nil {
		BusinessMap[userId] = make(map[string]interface{})
	}
	BusinessMap[userId][businessId] = infos
	bsLock.Unlock()
}

func SetAllBusinessMap(bsMap map[string]map[string]interface{}) {
	bsLock.Lock()
	BusinessMap = bsMap
	bsLock.Unlock()
}

func GetBusinessInfoStr(userId string, businessId string) string {
	tempMap := make(map[string]map[string]interface{})
	if tempMap[userId] == nil {
		tempMap[userId] = make(map[string]interface{})
	}
	tempMap[userId][businessId] = BusinessMap[userId][businessId]
	b, _ := json.Marshal(tempMap)
	if string(b) != "" && string(b) != "null" {
		return string(b)
	} else {
		return ""
	}
}
