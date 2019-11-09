package capibaribe

import "net/http"

type healthCheck struct {
	Path   string   `yaml:"path"   json:"path"`
	Header []header `yaml:"header" json:"header"`
	Body   string   `yaml:"body"   json:"body"`
}

func (el *healthCheck) VerifyPathToValidatePathIntoHost(path string) bool {
	return el.Path == path
}

func (el *healthCheck) WriteDataToOutputEndpoint(w http.ResponseWriter, r *http.Request) {
	for _, headerData := range el.Header {
		w.Header().Add(headerData.Key, headerData.Value)
	}

	w.Write([]byte(el.Body))
}
