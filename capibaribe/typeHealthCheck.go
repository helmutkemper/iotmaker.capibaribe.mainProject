package capibaribe

type healthCheck struct {
	Enabled         bool  `yaml:"enabled"`
	Interval        int   `yaml:"interval"`
	Fails           int   `yaml:"fails"`
	Passes          int   `yaml:"passes"`
	Uri             int   `yaml:"rui"`
	SuspendInterval int   `yaml:"suspendInterval"`
	Match           match `yaml:"match"`
}
