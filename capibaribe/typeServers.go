package capibaribe

import (
	"time"
)

const KListMaxLength = 10

type durationList struct {
	Duration time.Duration
	Error    bool
}

func (el *durationList) SuccessEvent(duration time.Duration) {
	el.Duration = duration
	el.Error = false
}

func (el *durationList) ErrorEvent(duration time.Duration) {
	el.Duration = duration
	el.Error = true
}

type servers struct {
	numberCurrentExecutions      int64
	executionDurationMax         time.Duration
	executionDurationMin         time.Duration
	executionDurationList        []durationList
	executionDurationSuccessList []durationList
	executionDurationErrorList   []durationList
	executionDurationAverage     time.Duration
	executionDateList            []time.Time
	executionDateSuccessList     []time.Time
	executionDateErrorList       []time.Time
	consecutiveErrors            int
	consecutiveSuccess           int
	totalErrorsCounter           int
	totalSuccessCounter          int
	lastRoundError               bool
	Host                         string  `yaml:"host"       json:"host"`
	Weight                       float64 `yaml:"weight"     json:"weight"`
	OverLoad                     int     `yaml:"overLoad"   json:"overLoad"`
}

func (el *servers) OnExecutionStartEvent() {
	el.AddExecutionDateToEntireExecutionDateList()
	el.OnExecutionStartCurrentExecutionsConterIncrementOne()
}

func (el *servers) OnExecutionEndWithErrorEvent(elapsedTime time.Duration) {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementErrosCounters()
	el.ResetConsecutiveSuccessCounter()
	el.SetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionDateToEntireExecutionDateList()
	el.AddExecutionDateToErrorExecutionDateList()
	el.AddExecutionTime(elapsedTime, true)
}

func (el *servers) OnExecutionEndWithSuccessEvent(elapsedTime time.Duration) {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementSuccessCounters()
	el.ResetConsecutiveErrosCounter()
	el.ResetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionDateToEntireExecutionDateList()
	el.AddExecutionDateToSuccessExecutionDateList()
	el.AddExecutionTime(elapsedTime, false)
}

func (el *servers) ResetConsecutiveErrosCounter() {
	el.consecutiveErrors = 0
}

func (el *servers) IncrementErrosCounters() {
	el.consecutiveErrors += 1
	el.totalErrorsCounter += 1
}

func (el *servers) ResetConsecutiveSuccessCounter() {
	el.consecutiveSuccess = 0
}

func (el *servers) IncrementSuccessCounters() {
	el.consecutiveSuccess += 1
	el.totalSuccessCounter += 1
}

func (el *servers) ResetRouteHasErrorOnLastRoundFlag() {
	el.lastRoundError = false
}

func (el *servers) SetRouteHasErrorOnLastRoundFlag() {
	el.lastRoundError = false
}

func (el *servers) OnExecutionStartCurrentExecutionsConterIncrementOne() {
	el.numberCurrentExecutions += 1
}

func (el *servers) OnExecutionEndCurrentExecutionsConterDecrementOne() {
	el.numberCurrentExecutions -= 1
}

func (el *servers) AddExecutionDateToEntireExecutionDateList() {
	el.executionDateList = append(el.executionDateList, time.Now())
	if len(el.executionDateList) > KListMaxLength {
		el.executionDateList = el.executionDateList[1:]
	}
}

func (el *servers) AddExecutionDateToSuccessExecutionDateList() {
	el.executionDateSuccessList = append(el.executionDateSuccessList, time.Now())
	if len(el.executionDateSuccessList) > KListMaxLength {
		el.executionDateSuccessList = el.executionDateSuccessList[1:]
	}
}

func (el *servers) AddExecutionDateToErrorExecutionDateList() {
	el.executionDateErrorList = append(el.executionDateErrorList, time.Now())
	if len(el.executionDateErrorList) > KListMaxLength {
		el.executionDateErrorList = el.executionDateErrorList[1:]
	}
}

func (el *servers) AddExecutionTime(duration time.Duration, error bool) {

	if duration > el.executionDurationMax {
		el.executionDurationMax = duration
	}

	if el.executionDurationMin == 0 {
		el.executionDurationMin = duration
	} else if el.executionDurationMin > duration {
		el.executionDurationMin = duration
	}

	el.executionDurationList = append(el.executionDurationList, durationList{Duration: duration, Error: error})
	if len(el.executionDurationList) > KListMaxLength {
		el.executionDurationList = el.executionDurationList[1:]
	}

	el.executionDurationAverage = 0
	for _, durationEvent := range el.executionDurationList {
		el.executionDurationAverage += durationEvent.Duration
	}

	el.executionDurationAverage = time.Duration(int64(el.executionDurationAverage) / int64(len(el.executionDurationList)))
}

func NewServerStruct() servers {
	ret := servers{}
	ret.executionDurationList = make([]durationList, 0)
	ret.executionDurationSuccessList = make([]durationList, 0)
	ret.executionDurationErrorList = make([]durationList, 0)

	ret.executionDateList = make([]time.Time, 0)
	ret.executionDateSuccessList = make([]time.Time, 0)
	ret.executionDateErrorList = make([]time.Time, 0)

	return ret
}
