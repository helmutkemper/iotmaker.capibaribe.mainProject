package capibaribe

import "crypto/tls"

type ssl struct {
	Enabled                  bool          `yaml:"enabled"`
	Certificate              string        `yaml:"certificate"`
	CertificateKey           string        `yaml:"certificateKey"`
	Version                  sslVersion    `yaml:"version"`
	CurvePreferences         []tls.CurveID `yaml:"curvePreferences"`
	PreferServerCipherSuites bool          `yaml:"preferServerCipherSuites"`
	CipherSuites             []uint16      `yaml:"cipherSuites"`
	X509                     sslX509       `yaml:"x509"`
}

type sslVersion struct {
	Min uint16 `yaml:"min"`
	Max uint16 `yaml:"max"`
}

type sslX509 struct {
	Certificate    string `yaml:"certificate"`
	CertificateKey string `yaml:"certificateKey"`
}
