package capibaribe

type match struct {
	Status []status `yaml:"status"`
	Header []header `yaml:"header"`
	Body   []string `yaml:"body"`
}
