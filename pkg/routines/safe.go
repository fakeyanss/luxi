package routines

import (
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		s := string(debug.Stack())
		log.Error().Msgf("err=%s, stack=%s", fmt.Sprint(p), s)
	}
}
