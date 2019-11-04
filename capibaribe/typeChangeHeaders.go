package capibaribe

type changeHeaders struct {
	Number int      `yaml:"number"  json:"number"`
	Header []header `yaml:"header"  json:"header"`
	Rate   float64  `yaml:"rate"    json:"rate"`
}
