package errs

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	goModPkgPath = "/go/pkg/mod/"
	skipFrames   = 2
)

var StackTraceMaxDepth int = 5

type stacktraceFrame struct {
	file     string
	function string
	pc       uintptr
	line     int
}

// String returns text representation of frame
func (frame *stacktraceFrame) String() string {
	currentFrame := fmt.Sprintf("%v:%v", frame.file, frame.line)
	if frame.function != "" {
		currentFrame = fmt.Sprintf("%v:%v %v()", frame.file, frame.line, frame.function)
	}

	return currentFrame
}

type stacktrace struct {
	frames []stacktraceFrame
}

// Error returns empty string
func (st *stacktrace) Error() string {
	return st.String("")
}

// String returns text reprentation of stacktrace
func (st *stacktrace) String(deepestFrame string) string {
	var str string

	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	for _, frame := range st.frames {
		if frame.file != "" {
			currentFrame := frame.String()
			if currentFrame == deepestFrame {
				break
			}

			newline()

			str += "  --- " + currentFrame
		}
	}

	return str
}

// newStacktrace creates Errs stacktrace
func newStacktrace() *stacktrace {
	var frames []stacktraceFrame

	for i := skipFrames; len(frames) < StackTraceMaxDepth+skipFrames; i++ {
		proc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		file = strings.TrimPrefix(file, goModPkgPath)

		f := runtime.FuncForPC(proc)
		if f == nil {
			break
		}

		function := shortFuncName(f)

		isGoPkg := len(runtime.GOROOT()) > 0 &&
			strings.Contains(file, runtime.GOROOT())

		if !isGoPkg {
			frames = append(frames, stacktraceFrame{
				pc:       proc,
				file:     file,
				function: function,
				line:     line,
			})
		}
	}

	return &stacktrace{
		frames: frames,
	}
}

// shortFuncName shortens func names removes special symbols
func shortFuncName(f *runtime.Func) string {
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
