package capibaribe

import (
	"time"
)

const KListMaxLength = 10

type executionInfo struct {
	Duration time.Duration
	Error    bool
	Date     time.Time
}

func (el *executionInfo) SuccessEvent(duration time.Duration) {
	el.Duration = duration
	el.Error = false
	el.Date = time.Now()
}

func (el *executionInfo) ErrorEvent(duration time.Duration) {
	el.Duration = duration
	el.Error = true
	el.Date = time.Now()
}

type analytics struct {
	NumberCurrentExecutions         int64
	ExecutionSuccessDurationMax     time.Duration
	ExecutionSuccessDurationMin     time.Duration
	ExecutionDurationEntireList     []executionInfo
	ExecutionDurationSuccessList    []executionInfo
	ExecutionDurationErrorList      []executionInfo
	ExecutionDurationSuccessAverage time.Duration
	ConsecutiveErrors               int
	ConsecutiveSuccess              int
	TotalErrorsCounter              int
	TotalSuccessCounter             int
	LastRoundError                  bool
}

type servers struct {
	analytics
	Host     string  `yaml:"host"       json:"host"`
	Weight   float64 `yaml:"weight"     json:"weight"`
	OverLoad int     `yaml:"overLoad"   json:"overLoad"`
}

func (el *servers) OnExecutionStartEvent() {
	el.OnExecutionStartCurrentExecutionsConterIncrementOne()
}

func (el *servers) OnExecutionEndWithErrorEvent(elapsedTime time.Duration) {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementErrosCounters()
	el.ResetConsecutiveSuccessCounter()
	el.SetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionTimeWithError(elapsedTime)
}

func (el *servers) OnExecutionEndWithSuccessEvent(elapsedTime time.Duration) {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementSuccessCounters()
	el.ResetConsecutiveErrosCounter()
	el.ResetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionTimeWithSuccess(elapsedTime)
}

func (el *analytics) ResetConsecutiveErrosCounter() {
	el.ConsecutiveErrors = 0
}

func (el *analytics) IncrementErrosCounters() {
	el.ConsecutiveErrors += 1
	el.TotalErrorsCounter += 1
}

func (el *analytics) ResetConsecutiveSuccessCounter() {
	el.ConsecutiveSuccess = 0
}

func (el *analytics) IncrementSuccessCounters() {
	el.ConsecutiveSuccess += 1
	el.TotalSuccessCounter += 1
}

func (el *analytics) ResetRouteHasErrorOnLastRoundFlag() {
	el.LastRoundError = false
}

func (el *analytics) SetRouteHasErrorOnLastRoundFlag() {
	el.LastRoundError = false
}

func (el *analytics) OnExecutionStartCurrentExecutionsConterIncrementOne() {
	el.NumberCurrentExecutions += 1
}

func (el *analytics) OnExecutionEndCurrentExecutionsConterDecrementOne() {
	el.NumberCurrentExecutions -= 1
}

func (el *analytics) AddExecutionTimeWithSuccess(duration time.Duration) {
	el.calculateMaxExecutionSuccessDuration(duration)
	el.calculateMinExecutionSuccessDuration(duration)
	el.addExecutionTimeToEntireList(duration, false)
	el.addExecutionTimeToSuccessList(duration)
	el.calculateExecutionSuccessDurationAverage()
}

func (el *analytics) AddExecutionTimeWithError(duration time.Duration) {
	el.addExecutionTimeToEntireList(duration, true)
	el.addExecutionTimeToErrorList(duration)
}

func (el *analytics) addExecutionTimeToEntireList(duration time.Duration, error bool) {
	el.ExecutionDurationEntireList = append(el.ExecutionDurationEntireList, executionInfo{Duration: duration, Error: error})
	if len(el.ExecutionDurationEntireList) > KListMaxLength {
		el.ExecutionDurationEntireList = el.ExecutionDurationEntireList[1:]
	}
}

func (el *analytics) addExecutionTimeToErrorList(duration time.Duration) {
	el.ExecutionDurationErrorList = append(el.ExecutionDurationErrorList, executionInfo{Duration: duration, Error: true})
	if len(el.ExecutionDurationErrorList) > KListMaxLength {
		el.ExecutionDurationErrorList = el.ExecutionDurationErrorList[1:]
	}
}

func (el *analytics) addExecutionTimeToSuccessList(duration time.Duration) {
	el.ExecutionDurationSuccessList = append(el.ExecutionDurationSuccessList, executionInfo{Duration: duration, Error: false})
	if len(el.ExecutionDurationSuccessList) > KListMaxLength {
		el.ExecutionDurationSuccessList = el.ExecutionDurationSuccessList[1:]
	}
}

func (el *analytics) calculateExecutionSuccessDurationAverage() {
	el.ExecutionDurationSuccessAverage = 0
	for _, durationEvent := range el.ExecutionDurationEntireList {
		el.ExecutionDurationSuccessAverage += durationEvent.Duration
	}

	el.ExecutionDurationSuccessAverage = time.Duration(int64(el.ExecutionDurationSuccessAverage) / int64(len(el.ExecutionDurationEntireList)))
}

func (el *analytics) calculateMaxExecutionSuccessDuration(duration time.Duration) {
	if duration > el.ExecutionSuccessDurationMax {
		el.ExecutionSuccessDurationMax = duration
	}
}

func (el *analytics) calculateMinExecutionSuccessDuration(duration time.Duration) {
	if el.ExecutionSuccessDurationMin == 0 || el.ExecutionSuccessDurationMin > duration {
		el.ExecutionSuccessDurationMin = duration
	}
}

func NewAnalytics() analytics {
	ret := analytics{}
	ret.ExecutionDurationEntireList = make([]executionInfo, 0)
	ret.ExecutionDurationSuccessList = make([]executionInfo, 0)
	ret.ExecutionDurationErrorList = make([]executionInfo, 0)

	return ret
}

func NewServerStruct(host string, weight float64, overLoad int) servers {

	ret := servers{}

	ret.analytics = NewAnalytics()

	ret.Host = host
	ret.Weight = weight
	ret.OverLoad = overLoad

	return ret
}
