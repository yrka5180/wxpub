package extra

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Default(level ...logrus.Level) {
	logrus.AddHook(new(GoIDFieldHook))
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{ // 由于统一使用日志收集器，因此使用json格式
		CallerPrettyfier: CallerPrettyfier,
	})

	// 日志级别默认为Info(整数4)
	_level := logrus.InfoLevel
	if len(level) > 0 {
		_level = level[0]
	}
	logrus.SetLevel(_level)
}

func SetLogger(logger *logrus.Logger) {
	logger.AddHook(new(GoIDFieldHook))
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{ // 由于统一使用日志收集器，因此使用json格式
		CallerPrettyfier: CallerPrettyfier,
	})
}

var fieldSeq = map[string]int{
	"time":   0,
	"nanoTS": 1,
	"level":  2,
	"file":   3,
	"func":   4,
	"goID":   5,
}

func CallerPrettyfier(f *runtime.Frame) (string, string) {
	a := strings.SplitN(path.Base(f.Function), ".", 2)
	return a[1], fmt.Sprintf("%s/%s:%d", a[0], path.Base(f.File), f.Line)
}

type GoIDFieldHook struct{}

func (h *GoIDFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GoIDFieldHook) Fire(e *logrus.Entry) error {
	e.Data["goID"] = fmt.Sprint(goID())
	e.Data["nanoTS"] = e.Time.UnixNano()
	return nil
}

var routinePrefix = []byte("goroutine ")

// Get go routine ID from runtime stack
func goID() string {
	buf := make([]byte, 32)
	runtime.Stack(buf, false)
	return string(bytes.Fields(bytes.TrimPrefix(buf, routinePrefix))[0])
}
