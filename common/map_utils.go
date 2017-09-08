package common

import (
	"reflect"
	"sort"
)

func GetSortedKeysInt64(datas map[string]int64) []string {
	keys := make([]string, 0)
	for k, _ := range datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSortedKeysFloat64(datas map[string]float64) []string {
	keys := make([]string, 0)
	for k, _ := range datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSortedKeys(datas map[string]interface{}) []string {
	keys := make([]string, 0)
	for k, _ := range datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSortedKeysReflect(datas reflect.Value, asc bool) []string {
	keys := make([]string, 0)

	if datas.Kind() == reflect.Map {
		for _, v1 := range datas.MapKeys() {
			keys = append(keys, v1.String())
		}

	}
	if asc {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	return keys
}
