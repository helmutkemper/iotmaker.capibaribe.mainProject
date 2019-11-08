package capibaribe

import (
	"time"
)

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

	timeOnStarEvent time.Time
}

func (el *analytics) StartTimeCounter() {
	el.timeOnStarEvent = time.Now()
}

func (el *analytics) getTimeInterval() time.Duration {
	return time.Since(el.timeOnStarEvent)
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

func (el *analytics) AddExecutionTimeWithSuccess() {
	duration := el.getTimeInterval()

	el.calculateMaxExecutionSuccessDuration(duration)
	el.calculateMinExecutionSuccessDuration(duration)
	el.addExecutionTimeToEntireList(duration, false)
	el.addExecutionTimeToSuccessList(duration)
	el.calculateExecutionSuccessDurationAverage()
}

func (el *analytics) AddExecutionTimeWithError() {
	duration := el.getTimeInterval()

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
