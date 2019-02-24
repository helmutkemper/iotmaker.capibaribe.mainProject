package capibaribe

type delay struct {
	Rate float64 `yaml:"rate"`
	Min  int     `yaml:"min"`
	Max  int     `yaml:"max"`
}
