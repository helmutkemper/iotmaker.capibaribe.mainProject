package capibaribe

type servers struct {
	consecutiveErrors  int
	consecutiveSuccess int
	errors             int
	success            int

	lastRoundError bool

	Host     string  `yaml:"host"`
	Weight   float64 `yaml:"weight"`
	OverLoad int     `yaml:"overLoad"`
}
