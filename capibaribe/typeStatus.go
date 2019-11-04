package capibaribe

type status struct {
	ExpReg string   `yaml:"expReg"  json:"expReg"`
	Value  int      `yaml:"value"   json:"value"`
	In     []maxMin `yaml:"in"      json:"in"`
	NotIn  []maxMin `yaml:"notIn"   json:"notIn"`
}
