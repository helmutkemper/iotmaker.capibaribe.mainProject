package capibaribe

import (
	"os"
	"runtime/debug"
	"strings"
)

type DebugLogger struct{}

func (d DebugLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, "multiple response.WriteHeader") {
		debug.PrintStack()
	}
	return os.Stderr.Write(p)
}
