package filter

import (
	"fmt"

	"../../common"
	"../../logp"
)

type Condition struct {
	Key      string      `config:"key"`
	Value    interface{} `config:"value"`
	AddKey   string      `config:"add_key"`
	AddValue interface{} `config:"add_value"`
}

type AddFieldsConfig struct {
	Fields     map[string]string `config:"fields"`
	Conditions []*Condition      `config:"conditions"`
}

type AddFields struct {
	Fields     map[string]string
	Conditions []*Condition
}

/* AddFields methods */
func NewAddFields(fields map[string]string, conditions []*Condition) *AddFields {
	return &AddFields{fields, conditions}
}

func (f *AddFields) Filter(event common.MapStr) (common.MapStr, error) {

	for key, value := range f.Fields {
		hasKey, err := event.HasKey(key)
		if err != nil {
			return event, fmt.Errorf("Fail to check the key %s: %s", key, err)
		}
		if hasKey {
			logp.Debug("Field:%v exists, ignore.", key)
		} else {
			event[key] = value
		}
	}

	return event, nil
}

func (f *AddFields) String() string {

	str := "add_fields="
	for k, v := range f.Fields {
		str += k + ":" + v + " "
	}

	return str
}
