package capibaribe

type header struct {
	Key   string `yaml:"key"    json:"key"`
	Value string `yaml:"value"  json:"value"`
}

type headerMonitor struct {
	Key   string `yaml:"key"    json:"key"`
	Value string `yaml:"value"  json:"value"`
	Type  string `yaml:"type"  json:"type"`
}
