package instabench

type Report struct {
	Config  map[string]interface{} `json:"config"`
	Results []*Results             `json:"results"`

	ExtraConfig map[string]interface{} `json:"extra_config"`
}

type ExtraConfigExporter interface {
	ExtraConfig() map[string]interface{}
}
