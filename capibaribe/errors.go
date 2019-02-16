package capibaribe

import (
	"github.com/helmutkemper/seelog"
	"runtime"
)

func HandleCriticalError(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)

		seelog.Criticalf("critical code error in [%s:%d] %v", fn, line, err)
	}
	return
}
