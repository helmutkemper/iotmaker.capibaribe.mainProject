package capibaribe

import (
	"encoding/json"
	"net/http"
)

type analyticsCheck struct {
	Path string `yaml:"path"   json:"path"`
}

func (el *analyticsCheck) VerifyPathToValidatePathIntoHost(path string) bool {
	return el.Path == path
}

func (el *analyticsCheck) WriteDataToOutputEndpoint(w http.ResponseWriter, proxy *proxy) {
	w.Header().Add("Content-Type", "application/json")

	out, _ := json.Marshal(proxy)

	w.Write(out)
}
