package sentryhook

import (
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

// NewStacktrace creates a stacktrace using `runtime.Callers`.
func NewStacktraceForHook() *sentry.Stacktrace {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(1, pcs)

	if n == 0 {
		return nil
	}

	frames := extractFrames(pcs[:n])
	frames = filterFrames(frames)

	stacktrace := sentry.Stacktrace{
		Frames: frames,
	}

	return &stacktrace
}

func extractFrames(pcs []uintptr) []sentry.Frame {
	var frames []sentry.Frame
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()

		frames = append([]sentry.Frame{
			sentry.NewFrame(callerFrame),
		}, frames...)

		if !more {
			break
		}
	}

	return frames
}


// filterFrames filters out stack frames that are not meant to be reported to
// Sentry. Those are frames internal to the SDK or Go.
func filterFrames(frames []sentry.Frame) []sentry.Frame {
	if len(frames) == 0 {
		return nil
	}

	filteredFrames := make([]sentry.Frame, 0, len(frames))

	for _, frame := range frames {
		// Skip Go internal frames.
		if frame.Module == "runtime" || frame.Module == "testing" {
			continue
		}
		// Skip Sentry internal frames, except for frames in _test packages (for
		// testing).
		if strings.HasPrefix(frame.Module, "github.com/getsentry/sentry-go") &&
			!strings.HasSuffix(frame.Module, "_test") {
			continue
		}
		// Skip Logrus Sentry Hook
		if strings.HasPrefix(frame.Module, "github.com/prixa-ai/logrus-sentry-hook") &&
			!strings.HasSuffix(frame.Module, "_test") {
			continue
		}
		// Skip Logrus
		if strings.HasPrefix(frame.Module, "github.com/sirupsen/logrus") &&
			!strings.HasSuffix(frame.Module, "_test") {
			continue
		}

		filteredFrames = append(filteredFrames, frame)
	}

	return filteredFrames
}
