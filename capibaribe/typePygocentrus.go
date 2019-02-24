package capibaribe

type pygocentrus struct {
	Enabled          bool            `yaml:"enabled"`
	Delay            rateMaxMin      `yaml:"delay"`
	DontRespond      rateMaxMin      `yaml:"dontRespond"`
	ChangeLength     float64         `yaml:"changeLength"`
	ChangeContent    changeContent   `yaml:"changeContent"`
	DeleteContent    float64         `yaml:"deleteContent"`
	ChangeHeaders    []changeHeaders `yaml:"changeHeaders"`
	successfulAttack bool            `yaml:"-"`
}

func (el *pygocentrus) SetAttack() {
	el.successfulAttack = true
}

func (el *pygocentrus) GetAttack() bool {
	return el.successfulAttack
}

func (el *pygocentrus) ClearAttack() {
	el.successfulAttack = false
}
