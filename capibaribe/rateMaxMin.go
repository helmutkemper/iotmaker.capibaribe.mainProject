package capibaribe

type rateMaxMin struct {
	Rate float64 `yaml:"rate" json:"rate"`
	Min  int     `yaml:"min"  json:"min"`
	Max  int     `yaml:"max"  json:"max"`
}
