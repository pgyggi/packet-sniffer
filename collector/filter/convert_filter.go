package filter

import (
	"../../common"
	"../../logp"
)

type Convert struct {
	Key   []string `config:"keys"`
	Value string   `config:"value"`
}

type ConvertFieldsConfig struct {
	Converts     []*Convert `config:"converts"`
	DeleteOrigin bool       `config:"delete_origin"`
}

type ConvertFields struct {
	Converts     []*Convert
	DeleteOrigin bool
}

/* AddFields methods */
func NewConvertFields(fields []*Convert, deleteOrigin bool) *ConvertFields {
	return &ConvertFields{fields, deleteOrigin}
}

func (f *ConvertFields) Filter(event common.MapStr) (common.MapStr, error) {
	for _, convert := range f.Converts {
		for _, key := range convert.Key {
			if val, exist := event[key]; exist {
				event[convert.Value] = val
				if f.DeleteOrigin {
					delete(event, key)
				}
			}
		}
	}

	logp.Debug("filter", "after convert filter:%v", event)
	return event, nil

}

func (f *ConvertFields) String() string {

	str := "convert_fields= "
	for _, v := range f.Converts {
		for _, key := range v.Key {
			str += key + ","
		}
		str += " => " + v.Value
	}

	return str
}
