package capibaribe

import (
	"bytes"
	"github.com/helmutkemper/seelog"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type transport struct {
	RoundTripper http.RoundTripper
	Project      *Project
}

type pygocentrusFunc func(req *http.Request) (resp *http.Response, err error)

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

func (el *transport) PygocentrusDelay(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were delayed by a pygocentrus attack: delay content", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(inLineIntRange(el.Project.Pygocentrus.Delay.Min, el.Project.Pygocentrus.Delay.Max)) * time.Microsecond)

	resp.Body = el.roundTripCopyBody(inBody)
	return resp, nil

}

func (el *transport) PygocentrusDontRespond(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: dont respond", req.RemoteAddr, req.RequestURI)
	return nil, nil

}

func (el *transport) PygocentrusDeleteContent(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}

	inBody = make([]byte, len(inBody))

	resp.Body = el.roundTripCopyBody(inBody)
	return resp, nil

}

func (el *transport) PygocentrusChangeContent(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}

	length := len(inBody)
	forLength := el.Project.Pygocentrus.ChangeContent.GetRandomByMaxMin(length)
	for i := 0; i != forLength; i += 1 {
		indexChange := el.Project.Pygocentrus.ChangeContent.GetRandomByLength(length)
		inBody = append(append(inBody[:indexChange], byte(inLineRand().Intn(255))), inBody[indexChange+1:]...)
	}

	resp.Body = el.roundTripCopyBody(inBody)
	return resp, nil

}

func (el *transport) PygocentrusChangeLength(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: change length", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(inBody))

	randLength := inLineRand().Intn(len(inBody))

	resp.ContentLength = int64(randLength)
	//resp.Header.Set("Content-Length", strconv.Itoa(randLength))
	return resp, nil

}

// todo: fazer
//func (el *transport) PygocentrusChangeHeaders(req *http.Request) (resp *http.Response, err error) {}

func (el *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if el.Project.Pygocentrus.Enabled == true {

		var randAttack int

		var list = make([]pygocentrusFunc, 0)

		if el.Project.Pygocentrus.DontRespond != 0.0 {

			list = append(list, el.PygocentrusDontRespond)

		}

		if el.Project.Pygocentrus.DeleteContent != 0.0 {

			list = append(list, el.PygocentrusDeleteContent)

		}

		if el.Project.Pygocentrus.ChangeContent.Rate != 0.0 {

			list = append(list, el.PygocentrusChangeContent)

		}

		if el.Project.Pygocentrus.ChangeLength != 0.0 {

			list = append(list, el.PygocentrusChangeLength)

		}

		/* todo: fazer
		   if el.Project.Pygocentrus.ChangeHeaders[0].Rate != 0.0 {}
		*/

		randAttack = inLineRand().Intn(len(list))
		return list[randAttack](req)

	}

	return el.RoundTripper.RoundTrip(req)
}
