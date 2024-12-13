package apperr

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/tecchu11/nrgo-std/nrslog"
)

type stacktrace []uintptr

func caller(size, skip int) stacktrace {
	pc := make([]uintptr, size)
	n := runtime.Callers(skip+1, pc)
	return pc[:n:n]
}

// String implements [fmt.Stringer].
func (s stacktrace) String() string {
	var str strings.Builder
	frames := runtime.CallersFrames(s)
	for {
		f, more := frames.Next()
		str.WriteString(f.Function)
		str.WriteString("(")
		str.WriteString(f.File)
		str.WriteString(":")
		str.WriteString(strconv.Itoa(f.Line))
		str.WriteString(")\n")
		if !more {
			break
		}
	}
	return str.String()
}

// LogValue implements [slog.LogValuer].
func (s stacktrace) LogValue() slog.Value {
	return slog.StringValue(s.String())
}

// NRAttributes implements [nrslog.Attributer].
func (s stacktrace) NRAttribute() map[string]string {
	return map[string]string{
		"stacktrace": s.String(),
	}
}

var (
	_ fmt.Stringer      = (stacktrace)(nil)
	_ slog.LogValuer    = (stacktrace)(nil)
	_ nrslog.Attributer = (stacktrace)(nil)
)
