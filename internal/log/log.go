package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/wf001/modo/pkg/types"
)

var DEFAULT_FORMAT = "%+v"

func getCaller() (string, string, int) {
	pc, file, line, _ := runtime.Caller(3)
	files := regexp.MustCompile("[/]").Split(file, -1)
	funcName := regexp.MustCompile("[/]").Split(runtime.FuncForPC(pc).Name(), -1)

	return funcName[len(funcName)-1], files[len(files)-1], line
}
func getLogrus() *logrus.Entry {
	funcName, file, line := getCaller()

	return logrus.WithFields(logrus.Fields{
		"func": funcName,
		"file": file,
		"line": line,
	})
}

func BLUE(body string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", body)
}

func DebugToken(tok *types.Token) {
	for ; tok != nil; tok = tok.Next {
		Debug(BLUE(fmt.Sprintf("token: %p %#+v", tok, tok)))
	}
}

func DebugMessage(message string) {
	Debug("", message)
}

func Debug(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	getLogrus().Debugf(format, value...)
}

func Info(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Infof(format, value...)
}

func Warn(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Warnf(format, value...)
}

func Error(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Errorf(format, value...)
}

func Panic(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Panicf(format, value...)
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
}

func SetOutputFile(fileName string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(f)
	defer f.Close()
}

func SetLevelError() {
	logrus.SetLevel(logrus.ErrorLevel)
}
func SetLevelWarning() {
	logrus.SetLevel(logrus.WarnLevel)
}
func SetLevelInfo() {
	logrus.SetLevel(logrus.InfoLevel)
}
func SetLevelDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}
