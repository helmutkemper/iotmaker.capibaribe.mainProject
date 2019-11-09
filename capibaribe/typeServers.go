package capibaribe

const KListMaxLength = 50

type servers struct {
	Analytics
	Host     string  `yaml:"host"       json:"host"`
	Weight   float64 `yaml:"weight"     json:"weight"`
	OverLoad int     `yaml:"overLoad"   json:"overLoad"`
}

func NewServerStruct(host string, weight float64, overLoad int) servers {

	ret := servers{}

	ret.Analytics = NewAnalytics()

	ret.Host = host
	ret.Weight = weight
	ret.OverLoad = overLoad

	return ret
}
