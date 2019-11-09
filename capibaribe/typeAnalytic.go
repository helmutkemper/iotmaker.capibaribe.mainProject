package capibaribe

import (
	"time"
)

// bt_br: Struct com a funcionalidade de registrar contadores e temporizadores em eventos de fim de execução, tanto para eventos com sucesso ou erro.
// en: A Struct that has the functionality of register, counters and timers, in the event of the end of the run, for both events, with a success or an error.
type Analytics struct {
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

// bt_br: Esta função deve ser chamada antes do início da execução
// en: This function must be called prior to the start of the run
func (el *Analytics) OnExecutionStart() {
	el.onExecutionStartCurrentExecutionsConterIncrementOne()
	el.startTimeCounter()
}

// pt_br: Esta função deve deve ser chamada ao final de uma execução com sucesso
// en: This function must be called at the end of a successful execution
func (el *Analytics) OnExecutionEndWithSuccess() {
	el.addExecutionTimeWithSuccess()
	el.incrementSuccessCounters()
	el.resetConsecutiveErrosCounter()
	el.resetRouteHasErrorOnLastRoundFlag()
	el.onExecutionEndCurrentExecutionsConterDecrementOne()
}

// bt_br: Esta função deve ser chamada ao final de uma execução com erro
// en: This function should be called at the end of the run with an error
func (el *Analytics) OnExecutionEndWithError() {
	el.addExecutionTimeWithError()
	el.resetConsecutiveSuccessCounter()
	el.incrementErrosCounters()
	el.setRouteHasErrorOnLastRoundFlag()
	el.onExecutionEndCurrentExecutionsConterDecrementOne()
}

func (el *Analytics) GetLastRoundError() bool {
	return el.LastRoundError
}

func (el *Analytics) startTimeCounter() {
	el.timeOnStarEvent = time.Now()
}

func (el *Analytics) getTimeInterval() time.Duration {
	return time.Since(el.timeOnStarEvent)
}

func (el *Analytics) resetConsecutiveErrosCounter() {
	el.ConsecutiveErrors = 0
}

func (el *Analytics) incrementErrosCounters() {
	el.ConsecutiveErrors += 1
	el.TotalErrorsCounter += 1
}

func (el *Analytics) resetConsecutiveSuccessCounter() {
	el.ConsecutiveSuccess = 0
}

func (el *Analytics) incrementSuccessCounters() {
	el.ConsecutiveSuccess += 1
	el.TotalSuccessCounter += 1
}

func (el *Analytics) resetRouteHasErrorOnLastRoundFlag() {
	el.LastRoundError = false
}

func (el *Analytics) setRouteHasErrorOnLastRoundFlag() {
	el.LastRoundError = true
}

/**
@see OnExecutionStart()
*/
func (el *Analytics) onExecutionStartCurrentExecutionsConterIncrementOne() {
	el.NumberCurrentExecutions += 1
}

/**
@see OnExecutionEndWithError()
@see OnExecutionEndWithSuccess()
*/
func (el *Analytics) onExecutionEndCurrentExecutionsConterDecrementOne() {
	el.NumberCurrentExecutions -= 1
}

/**
@OnExecutionEndWithSuccess()
*/
func (el *Analytics) addExecutionTimeWithSuccess() {
	duration := el.getTimeInterval()

	el.calculateMaxExecutionSuccessDuration(duration)
	el.calculateMinExecutionSuccessDuration(duration)
	el.addExecutionTimeToEntireList(duration, false)
	el.addExecutionTimeToSuccessList(duration)
	el.calculateExecutionSuccessDurationAverage()
}

/**
@OnExecutionEndWithError()
*/
func (el *Analytics) addExecutionTimeWithError() {
	duration := el.getTimeInterval()

	el.addExecutionTimeToEntireList(duration, true)
	el.addExecutionTimeToErrorList(duration)
}

/**
@see addExecutionTimeWithSuccess()
@see addExecutionTimeWithError()
*/
func (el *Analytics) addExecutionTimeToEntireList(duration time.Duration, error bool) {
	el.ExecutionDurationEntireList = append(el.ExecutionDurationEntireList, executionInfo{Duration: duration, Error: error, Date: time.Now()})
	if len(el.ExecutionDurationEntireList) > KListMaxLength {
		el.ExecutionDurationEntireList = el.ExecutionDurationEntireList[1:]
	}
}

/**
@addExecutionTimeWithError()
*/
func (el *Analytics) addExecutionTimeToErrorList(duration time.Duration) {
	el.ExecutionDurationErrorList = append(el.ExecutionDurationErrorList, executionInfo{Duration: duration, Error: true, Date: time.Now()})
	if len(el.ExecutionDurationErrorList) > KListMaxLength {
		el.ExecutionDurationErrorList = el.ExecutionDurationErrorList[1:]
	}
}

/**
@see addExecutionTimeWithSuccess()
*/
func (el *Analytics) addExecutionTimeToSuccessList(duration time.Duration) {
	el.ExecutionDurationSuccessList = append(el.ExecutionDurationSuccessList, executionInfo{Duration: duration, Error: false, Date: time.Now()})
	if len(el.ExecutionDurationSuccessList) > KListMaxLength {
		el.ExecutionDurationSuccessList = el.ExecutionDurationSuccessList[1:]
	}
}

/**
@see addExecutionTimeWithSuccess()
*/
func (el *Analytics) calculateExecutionSuccessDurationAverage() {
	el.ExecutionDurationSuccessAverage = 0
	for _, durationEvent := range el.ExecutionDurationEntireList {
		el.ExecutionDurationSuccessAverage += durationEvent.Duration
	}

	el.ExecutionDurationSuccessAverage = time.Duration(int64(el.ExecutionDurationSuccessAverage) / int64(len(el.ExecutionDurationEntireList)))
}

/**
@see addExecutionTimeWithSuccess()
*/
func (el *Analytics) calculateMaxExecutionSuccessDuration(duration time.Duration) {
	if duration > el.ExecutionSuccessDurationMax {
		el.ExecutionSuccessDurationMax = duration
	}
}

/**
@see addExecutionTimeWithSuccess()
*/
func (el *Analytics) calculateMinExecutionSuccessDuration(duration time.Duration) {
	if el.ExecutionSuccessDurationMin == 0 || el.ExecutionSuccessDurationMin > duration {
		el.ExecutionSuccessDurationMin = duration
	}
}
