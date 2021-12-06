package runtime

import (
	"sync"
)

type RuntimeOptions struct {
	Debug                  bool
	HashTrackingFunc       func(collection string, key string, data []byte) (string, string)
	DisableParam           bool
	DisableParamExceptions []string
}

var (
	doOnce         sync.Once
	runtimeOptions *RuntimeOptions
)

func SetRuntimeOptions(debug bool, hashTrackingFunc func(string, string, []byte) (string, string), disableParam bool, disableParamExceptions []string) {
	doOnce.Do(func() {
		runtimeOptions = &RuntimeOptions{
			Debug:                  debug,
			HashTrackingFunc:       hashTrackingFunc,
			DisableParam:           disableParam,
			DisableParamExceptions: disableParamExceptions,
		}
	})
}

func GetRuntimeOptions() *RuntimeOptions {
	return runtimeOptions
}
