package capibaribe

type ssl struct {
	Enabled                  bool        `yaml:"enabled"                  json:"enabled"`
	Certificate              string      `yaml:"certificate"              json:"certificate"`
	CertificateKey           string      `yaml:"certificateKey"           json:"certificateKey"`
	Version                  sslVersion  `yaml:"version"                  json:"version"`
	CurvePreferences         interface{} `yaml:"curvePreferences"         json:"curvePreferences"`
	PreferServerCipherSuites bool        `yaml:"preferServerCipherSuites" json:"preferServerCipherSuites"`
	CipherSuites             interface{} `yaml:"cipherSuites"             json:"cipherSuites"`
	X509                     sslX509     `yaml:"x509"                     json:"x509"`
}

type sslVersion struct {
	Min uint16 `yaml:"min" json:"min"`
	Max uint16 `yaml:"max" json:"max"`
}

type sslX509 struct {
	Certificate    string `yaml:"certificate"     json:"certificate"`
	CertificateKey string `yaml:"certificateKey"  json:"certificateKey"`
}
