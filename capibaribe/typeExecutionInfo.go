package capibaribe

import (
	"time"
)

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
