package capibaribe

type pygocentrus struct {
	Enabled       bool            `yaml:"enabled"`
	DontRespond   float64         `yaml:"dontRespond"`
	ChangeLength  float64         `yaml:"changeLength"`
	ChangeContent changeContent   `yaml:"changeContent"`
	DeleteContent float64         `yaml:"deleteContent"`
	ChangeHeaders []changeHeaders `yaml:"changeHeaders"`
}
