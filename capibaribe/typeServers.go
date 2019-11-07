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
	NumberCurrentExecutions      int64
	ExecutionDurationMax         time.Duration
	ExecutionDurationMin         time.Duration
	ExecutionDurationList        []durationList
	ExecutionDurationSuccessList []durationList
	ExecutionDurationErrorList   []durationList
	ExecutionDurationAverage     time.Duration
	ExecutionDateList            []time.Time
	ExecutionDateSuccessList     []time.Time
	ExecutionDateErrorList       []time.Time
	ConsecutiveErrors            int
	ConsecutiveSuccess           int
	TotalErrorsCounter           int
	TotalSuccessCounter          int
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
	el.AddExecutionTimeWithError(elapsedTime)
}

func (el *servers) OnExecutionEndWithSuccessEvent(elapsedTime time.Duration) {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementSuccessCounters()
	el.ResetConsecutiveErrosCounter()
	el.ResetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionDateToEntireExecutionDateList()
	el.AddExecutionDateToSuccessExecutionDateList()
	el.AddExecutionTimeWithSuccess(elapsedTime)
}

func (el *servers) ResetConsecutiveErrosCounter() {
	el.ConsecutiveErrors = 0
}

func (el *servers) IncrementErrosCounters() {
	el.ConsecutiveErrors += 1
	el.TotalErrorsCounter += 1
}

func (el *servers) ResetConsecutiveSuccessCounter() {
	el.ConsecutiveSuccess = 0
}

func (el *servers) IncrementSuccessCounters() {
	el.ConsecutiveSuccess += 1
	el.TotalSuccessCounter += 1
}

func (el *servers) ResetRouteHasErrorOnLastRoundFlag() {
	el.lastRoundError = false
}

func (el *servers) SetRouteHasErrorOnLastRoundFlag() {
	el.lastRoundError = false
}

func (el *servers) OnExecutionStartCurrentExecutionsConterIncrementOne() {
	el.NumberCurrentExecutions += 1
}

func (el *servers) OnExecutionEndCurrentExecutionsConterDecrementOne() {
	el.NumberCurrentExecutions -= 1
}

func (el *servers) AddExecutionDateToEntireExecutionDateList() {
	el.ExecutionDateList = append(el.ExecutionDateList, time.Now())
	if len(el.ExecutionDateList) > KListMaxLength {
		el.ExecutionDateList = el.ExecutionDateList[1:]
	}
}

func (el *servers) AddExecutionDateToSuccessExecutionDateList() {
	el.ExecutionDateSuccessList = append(el.ExecutionDateSuccessList, time.Now())
	if len(el.ExecutionDateSuccessList) > KListMaxLength {
		el.ExecutionDateSuccessList = el.ExecutionDateSuccessList[1:]
	}
}

func (el *servers) AddExecutionDateToErrorExecutionDateList() {
	el.ExecutionDateErrorList = append(el.ExecutionDateErrorList, time.Now())
	if len(el.ExecutionDateErrorList) > KListMaxLength {
		el.ExecutionDateErrorList = el.ExecutionDateErrorList[1:]
	}
}

func (el *servers) AddExecutionTimeWithSuccess(duration time.Duration) {
	el.addExecutionTime(duration, false)
}

func (el *servers) AddExecutionTimeWithError(duration time.Duration) {
	el.addExecutionTime(duration, true)
}

func (el *servers) addExecutionTime(duration time.Duration, error bool) {

	if duration > el.ExecutionDurationMax {
		el.ExecutionDurationMax = duration
	}

	if el.ExecutionDurationMin == 0 {
		el.ExecutionDurationMin = duration
	} else if el.ExecutionDurationMin > duration {
		el.ExecutionDurationMin = duration
	}

	el.ExecutionDurationList = append(el.ExecutionDurationList, durationList{Duration: duration, Error: error})
	if len(el.ExecutionDurationList) > KListMaxLength {
		el.ExecutionDurationList = el.ExecutionDurationList[1:]
	}

	el.ExecutionDurationAverage = 0
	for _, durationEvent := range el.ExecutionDurationList {
		el.ExecutionDurationAverage += durationEvent.Duration
	}

	el.ExecutionDurationAverage = time.Duration(int64(el.ExecutionDurationAverage) / int64(len(el.ExecutionDurationList)))
}

func NewServerStruct() servers {
	ret := servers{}
	ret.ExecutionDurationList = make([]durationList, 0)
	ret.ExecutionDurationSuccessList = make([]durationList, 0)
	ret.ExecutionDurationErrorList = make([]durationList, 0)

	ret.ExecutionDateList = make([]time.Time, 0)
	ret.ExecutionDateSuccessList = make([]time.Time, 0)
	ret.ExecutionDateErrorList = make([]time.Time, 0)

	return ret
}
