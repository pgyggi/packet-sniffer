package common

import (
	"encoding/json"
	//	"fmt"
	//	"io/ioutil"
	//	"time"
)

// 		  ds/ss     -userId    -dbid/sid  -infos
var DBMap map[string]map[string]map[string]map[string]interface{}

func GetDBMap(userId string, dbId string) map[string]interface{} {
	return DBMap[DASH_BOARD_DS][userId][dbId]
}

func GetAllDBMap() map[string]map[string]map[string]map[string]interface{} {
	return DBMap
}

func GetDBByUser(userId string) map[string]map[string]interface{} {
	return DBMap[DASH_BOARD_DS][userId]
}

func GetAllDBMapByUser(userId string) map[string]map[string]interface{} {
	return DBMap[DASH_BOARD_DS][userId]
}

func SetDBMap(userId string, dbId string, valMap map[string]interface{}) {
	if DBMap == nil {
		DBMap = make(map[string]map[string]map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_DS] == nil {
		DBMap[DASH_BOARD_DS] = make(map[string]map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_DS][userId] == nil {
		DBMap[DASH_BOARD_DS][userId] = make(map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_DS][userId][dbId] == nil {
		DBMap[DASH_BOARD_DS][userId][dbId] = make(map[string]interface{})
	}
	DBMap[DASH_BOARD_DS][userId][dbId] = valMap
}

func SetAllDBMap(valMap map[string]map[string]map[string]map[string]interface{}) {
	DBMap = valMap
}
func GetSearchMap(userId string, searchId string) map[string]interface{} {
	if DBMap == nil {
		DBMap = make(map[string]map[string]map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_SS] == nil {
		DBMap[DASH_BOARD_SS] = make(map[string]map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_SS][userId] == nil {
		DBMap[DASH_BOARD_SS][userId] = make(map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_SS][userId][searchId] == nil {
		DBMap[DASH_BOARD_SS][userId][searchId] = make(map[string]interface{})
	}
	return DBMap[DASH_BOARD_SS][userId][searchId]
}

func GetSearchByUser(userId string) map[string]map[string]interface{} {
	return DBMap[DASH_BOARD_SS][userId]
}

func SetSearchMap(userId string, searchId string, valMap map[string]interface{}) {
	if DBMap[DASH_BOARD_SS] == nil {
		DBMap[DASH_BOARD_SS] = make(map[string]map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_SS][userId] == nil {
		DBMap[DASH_BOARD_SS][userId] = make(map[string]map[string]interface{})
	}
	if DBMap[DASH_BOARD_SS][userId][searchId] == nil {
		DBMap[DASH_BOARD_SS][userId][searchId] = make(map[string]interface{})
	}

	DBMap[DASH_BOARD_SS][userId][searchId] = valMap
}
func GetDBInfoStr() string {
	b, _ := json.Marshal(DBMap)
	if string(b) != "" && string(b) != "null" {
		return string(b)
	} else {
		return ""
	}
}

func GetDBByUserStr(userId string, dbId string) string {
	b, _ := json.Marshal(DBMap[DASH_BOARD_DS][userId][dbId])
	if string(b) != "" && string(b) != "null" {
		return string(b)
	} else {
		return ""
	}
}

func GetSearchInfoByUserStr(userId string, searchid string) string {
	b, _ := json.Marshal(DBMap[DASH_BOARD_SS][userId][searchid])
	if string(b) != "" && string(b) != "null" {
		return string(b)
	} else {
		return ""
	}
}
