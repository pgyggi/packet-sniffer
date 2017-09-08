package filter

type DropFieldsConfig struct {
	Fields []string `config:"fields"`
}

type IncludeFieldsConfig struct {
	Fields []string `config:"fields"`
}

type FilterConfig struct {
	DropFields    *DropFieldsConfig    `config:"drop_fields_filter"`
	IncludeFields *IncludeFieldsConfig `config:"include_fields_filter"`
	AddFields     *AddFieldsConfig     `config:"add_fields_filter"`
	ConvertFields *ConvertFieldsConfig `config:"convert_fields_filter"`
	ReplaceFields *ReplaceFieldsConfig `config:"replace_fields_filter"`
	GrokFields    *GrokFieldsConfig    `config:"grok_filter"`
	IPFields      *IPFieldsConfig      `config:"ip_filter"`
}

// fields that should be always exported
var MandatoryExportedFields = []string{"timestamp", "type"}
