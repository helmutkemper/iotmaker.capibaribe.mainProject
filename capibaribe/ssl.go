package capibaribe

type ssl struct {
	Enabled        bool     `yaml:"enabled"`
	Certificate    string   `yaml:"certificate"`
	CertificateKey string   `yaml:"certificateKey"`
	SllProtocols   []string `yaml:"sslProtocols"`
}
