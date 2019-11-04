package capibaribe

type servers struct {
	consecutiveErrors  int
	consecutiveSuccess int
	errors             int
	success            int

	lastRoundError bool

	Host     string  `yaml:"host"       json:"host"`
	Weight   float64 `yaml:"weight"     json:"weight"`
	OverLoad int     `yaml:"overLoad"   json:"overLoad"`
}
