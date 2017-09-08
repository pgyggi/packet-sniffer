package filter

import (
	"strings"

	"../../common"
	"../../logp"
	"github.com/vjeantet/grok"
)

type Pattern struct {
	Key   string `config:"key"`
	Value string `config:"value"`
}

type GrokFieldsConfig struct {
	Patterns  []Pattern `config:"register"`
	Fields    []string  `config:"fields"`
	Splite    string    `config:"splite"`
	Pattern   string    `config:"pattern"`
	DelOrigin bool      `config:"delete_origin"`
}

type GrokFields struct {
	Patterns  []Pattern
	Fields    []string
	Splite    string
	Pattern   string
	DelOrigin bool
}

/* SysOidFields methods */
func NewGrokFields(patterns []Pattern, fields []string, splite string, pattern string, delOrigin bool) *GrokFields {
	return &GrokFields{patterns, fields, splite, pattern, delOrigin}
}

func (f *GrokFields) Filter(event common.MapStr) (common.MapStr, error) {

	logp.Debug("grok_filter", "Before filter :%v", event)

	g, _ := grok.New()
	for _, pattern := range f.Patterns {
		g.AddPattern(pattern.Key, pattern.Value)
	}

	var vals []string

	for _, field := range f.Fields {
		if val, exist := event[field]; exist {

			if f.Splite != "" {
				vals = strings.Split(val.(string), f.Splite)
			} else {
				vals = append(vals, val.(string))
			}
			for _, tmp := range vals {
				if vals, err := g.Parse(f.Pattern, tmp); err != nil {
					logp.Err("Parse grok error:%v", err)
				} else {
					for k, v := range vals {
						event[k] = v
					}
				}
			}
			if f.DelOrigin {
				delete(event, field)
			}
		}
	}

	logp.Debug("grok_filter", "After filter :%v", event)

	return event, nil
}

func (f *GrokFields) String() string {
	str := "grok_fields init."
	return str
}
