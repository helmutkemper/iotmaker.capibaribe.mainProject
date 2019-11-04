package capibaribe

import (
	"net"
	"time"
)

type pygocentrusListenFunc func(inData []byte) (int, []byte)

type Listen struct {
	InProtocol  string                `yaml:"inProtocol"  json:"inProtocol"`
	InAddress   string                `yaml:"inAddress"   json:"inAddress"`
	OutProtocol string                `yaml:"outProtocol" json:"outProtocol"`
	OutAddress  string                `yaml:"outAddress"  json:"outAddress"`
	Pygocentrus pygocentrus           `yaml:"pygocentrus" json:"pygocentrus"`
	stop        bool                  `yaml:"-"           json:"-"`
	attack      pygocentrusListenFunc `yaml:"-"           json:"-"`
}

func (el *Listen) Listen() error {
	listener, err := net.Listen(el.InProtocol, el.InAddress)
	if err != nil {
		return err
	}

	for {

		if el.stop == true {
			el.stop = false
			return nil
		}

		inConn, err := listener.Accept()
		if err != nil {
			return err
		}

		go el.handleRequest(inConn)
	}
}

func (el *Listen) handleRequest(inConn net.Conn) error {
	outConn, err := net.Dial(el.OutProtocol, el.OutAddress)
	if err != nil {
		return err
	}

	chan1 := el.chanFromConn(inConn)
	chan2 := el.chanFromConn(outConn)

	if el.Pygocentrus.Enabled == true {

		var randAttack int

		var list = make([]pygocentrusListenFunc, 0)

		if el.Pygocentrus.Delay.Rate != 0.0 {

			if el.Pygocentrus.Delay.Rate >= inLineRand().Float64() {
				list = append(list, el.pygocentrusDelay)
			}

		}

		if el.Pygocentrus.DontRespond.Rate != 0.0 {

			if el.Pygocentrus.DontRespond.Rate >= inLineRand().Float64() {
				list = append(list, el.pygocentrusDontRespond)
			}

		}

		if el.Pygocentrus.DeleteContent != 0.0 {

			if el.Pygocentrus.DeleteContent >= inLineRand().Float64() {
				list = append(list, el.pygocentrusDeleteContent)
			}

		}

		if el.Pygocentrus.ChangeContent.Rate != 0.0 {

			if el.Pygocentrus.ChangeContent.Rate >= inLineRand().Float64() {
				list = append(list, el.pygocentrusChangeContent)
			}

		}

		listLength := len(list)
		if listLength != 0 {
			el.Pygocentrus.SetAttack()
			randAttack = inLineRand().Intn(len(list))
			el.attack = list[randAttack]
		}
	}

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return nil
			} else {
				_, err = outConn.Write(b1)
				if err != nil {
					return err
				}
			}
		case b2 := <-chan2:
			if b2 == nil {
				return nil
			} else {
				_, err = inConn.Write(b2)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (el *Listen) chanFromConn(conn net.Conn) chan []byte {
	channel := make(chan []byte)

	go func() {
		buff := make([]byte, 1024)

		for {
			n, err := conn.Read(buff)

			if el.attack != nil {
				n, buff = el.attack(buff)
			}

			if n > 0 {
				res := make([]byte, n)
				copy(res, buff[:n])
				channel <- res
			}
			if err != nil {
				channel <- nil
				break
			}
		}
	}()

	return channel
}

func (el *Listen) pygocentrusDelay(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were delayed by a pygocentrus attack: delay content", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)

	return len(inData), inData

}

func (el *Listen) pygocentrusDontRespond(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: dont respond", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)
	return 0, nil

}

func (el *Listen) pygocentrusDeleteContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	n := len(inData)
	inData = make([]byte, n)

	return n, inData

}

func (el *Listen) pygocentrusChangeContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

	length := len(inData)
	forLength := el.Pygocentrus.ChangeContent.GetRandomByMaxMin(length)
	for i := 0; i != forLength; i += 1 {
		indexChange := el.Pygocentrus.ChangeContent.GetRandomByLength(length)
		inData = append(append(inData[:indexChange], byte(inLineRand().Intn(255))), inData[indexChange+1:]...)
	}

	return length, inData

}
