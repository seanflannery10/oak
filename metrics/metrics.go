package metrics

import (
	"expvar"
	"github.com/seanflannery10/ossa/version"
	"runtime"
	"time"
)

func Common() {
	expvar.NewString("version").Set(version.Get())

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}
