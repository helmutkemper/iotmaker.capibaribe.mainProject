package capibaribe

type changeHeaders struct {
	Number int      `yaml:"number"`
	Header []header `yaml:"header"`
	Rate   float64  `yaml:"rate"`
}
