package filter

import (
	"strings"

	"../../common"
	"../../logp"
)

type Replace struct {
	Src string `config:"src"`
	New string `config:"new"`
}

type ReplaceFieldsConfig struct {
	Replacement []*Replace `config:"replace"`
}

type ReplaceFields struct {
	Replacement []*Replace
}

/* SysOidFields methods */
func NewReplaceFields(replacement []*Replace) *ReplaceFields {
	return &ReplaceFields{replacement}
}

func (f *ReplaceFields) Filter(event common.MapStr) (common.MapStr, error) {

	for ek, ev := range event {
		for _, v := range f.Replacement {
			if strings.Contains(ek, v.Src) {
				newkey := strings.Replace(ek, v.Src, v.New, -1)
				event[newkey] = ev
				delete(event, ek)
			}
		}
	}

	logp.Debug("filter", "after replace filter:%v", event)

	return event, nil
}

func (f *ReplaceFields) String() string {

	str := "replace_fields="
	for _, v := range f.Replacement {
		str += v.Src + " replace with " + v.New + " "
	}

	return str
}
