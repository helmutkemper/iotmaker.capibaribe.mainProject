package capibaribe

type status struct {
	ExpReg string   `yaml:"expReg"`
	Value  int      `yaml:"value"`
	In     []maxMin `yaml:"in"`
	NotIn  []maxMin `yaml:"notIn"`
}
