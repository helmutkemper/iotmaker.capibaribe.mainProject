package capibaribe

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type transport struct {
	RoundTripper http.RoundTripper
	Project      *Project
}

func (el *transport) roundTripReadBody(req *http.Request) (*http.Response, []byte, error) {
	var resp *http.Response
	var err error
	var inBody []byte

	resp, err = el.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, nil, err
	}

	inBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return resp, nil, err
	}

	return resp, inBody, err
}

func (el *transport) roundTripCopyBody(inBody []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(inBody))
}

func (el *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return el.RoundTripper.RoundTrip(req)
}
