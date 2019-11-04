package capibaribe

type healthCheck struct {
	Enabled         bool  `yaml:"enabled"         json:"enabled"`
	Interval        int   `yaml:"interval"        json:"interval"`
	Fails           int   `yaml:"fails"           json:"fails"`
	Passes          int   `yaml:"passes"          json:"passes"`
	Uri             int   `yaml:"rui"             json:"rui"`
	SuspendInterval int   `yaml:"suspendInterval" json:"suspendInterval"`
	Match           match `yaml:"match"           json:"match"`
}
