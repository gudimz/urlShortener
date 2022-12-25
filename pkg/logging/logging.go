package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var entry *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{entry}
}

type writerHook struct {
	Writers   []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writers {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func init() {
	logger := logrus.New()
	logger.SetReportCaller(true)

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		// Print the function name and line number
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	})

	// create dir with logs
	err := os.MkdirAll("logs", 0750)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// create logfile
	logFile, err := os.OpenFile("logs/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// because we need write logs in stdout and log file
	logger.SetOutput(io.Discard)

	// write logs in log file and stdout
	logger.AddHook(&writerHook{
		Writers:   []io.Writer{logFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logger.SetLevel(logrus.TraceLevel)
	entry = logrus.NewEntry(logger)
}
