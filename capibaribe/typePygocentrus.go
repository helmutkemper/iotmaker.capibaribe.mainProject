package capibaribe

type pygocentrus struct {
	Enabled          bool            `yaml:"enabled"        json:"enabled"`
	Delay            rateMaxMin      `yaml:"delay"          json:"delay"`
	DontRespond      rateMaxMin      `yaml:"dontRespond"    json:"dontRespond"`
	ChangeLength     float64         `yaml:"changeLength"   json:"changeLength"`
	ChangeContent    changeContent   `yaml:"changeContent"  json:"changeContent"`
	DeleteContent    float64         `yaml:"deleteContent"  json:"deleteContent"`
	ChangeHeaders    []changeHeaders `yaml:"changeHeaders"  json:"changeHeaders"`
	successfulAttack bool            `yaml:"-"              json:"-"`
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
