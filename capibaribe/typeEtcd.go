package capibaribe

import (
	"time"
)

type etcd struct {
	DialTimeOut    time.Duration `yaml:"dialTimeout"`
	RequestTimeout time.Duration `yaml:"requestTimeout"`
	Connection     []string      `yaml:"connection"`
	ConfigKey      string        `yaml:"configKey"`
}

func (el *etcd) Prepare() {
	if el.ConfigKey == "" {
		el.ConfigKey = "capibaribe-config-yaml-file"
	}

	if len(el.Connection) == 0 {
		el.Connection = []string{"127.0.0.1:2379"}
	}

	if el.DialTimeOut == 0 {
		el.DialTimeOut = 2 * time.Second
	}

	if el.RequestTimeout == 0 {
		el.RequestTimeout = 2 * time.Second
	}
}
