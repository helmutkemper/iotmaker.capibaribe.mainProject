package capibaribe

import (
	"errors"
	"math/rand"
	"time"
)

type changeContent struct {
	ChangeRateMin  float64 `yaml:"changeRateMin"  json:"changeRateMin"`
	ChangeRateMax  float64 `yaml:"changeRateMax"  json:"changeRateMax"`
	ChangeBytesMin int     `yaml:"changeBytesMin" json:"changeBytesMin"`
	ChangeBytesMax int     `yaml:"changeBytesMax" json:"changeBytesMax"`
	Rate           float64 `yaml:"rate"           json:"rate"`
}

func (el *changeContent) prepare() error {
	if el.Rate == 0.0 {
		return nil
	}

	if el.ChangeRateMin == el.ChangeRateMax && el.ChangeBytesMin == el.ChangeBytesMax && el.ChangeRateMin == 0.0 {
		el.Rate = 0.0
		return errors.New("pygocentrus attack > changeContent > rate set to zero")
	}

	if el.ChangeRateMin > el.ChangeRateMax {
		return errors.New("pygocentrus attack > changeContent > rate > the minimum value is greater than the maximum value")
	}

	if el.ChangeBytesMin > el.ChangeBytesMax {
		return errors.New("pygocentrus attack > changeContent > bytes > the minimum value is greater than the maximum value")
	}

	if (el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0) && (el.ChangeBytesMin != 0.0 || el.ChangeBytesMax != 0.0) {
		return errors.New("pygocentrus attack > changeContent > you must choose option rate change or option bytes change")
	}

	return nil
}

func (el *changeContent) GetRandomByMaxMin(length int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	if el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0 {
		var changeMin = int(float64(length) * el.ChangeRateMin)
		var changeMax = int(float64(length) * el.ChangeRateMax)

		return r1.Intn(changeMax-changeMin) + changeMin
	}

	return r1.Intn(el.ChangeBytesMax-el.ChangeBytesMin) + el.ChangeBytesMin
}

func (el *changeContent) GetRandomByLength(length int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(length)
}
