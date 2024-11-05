package log

import (
	"os"
	"regexp"
	"runtime"

	"github.com/sirupsen/logrus"
)

var DEFAULT_FORMAT = "%+v"

func getCaller() (string, int) {
	_, file, line, _ := runtime.Caller(3)
	files := regexp.MustCompile("[/]").Split(file, -1)

	return files[len(files)-1], line
}
func getLogrus() *logrus.Entry {
	file, line := getCaller()

	return logrus.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	})
}


func Debug(value interface{}, args ...string) {
	format := DEFAULT_FORMAT

	if len(args) > 0 {
		format = args[0]
	}
  getLogrus().Debugf(format, value)
}

func Info(value interface{}, args ...string) {
	format := DEFAULT_FORMAT
	if len(args) > 0 {
		format = args[0]
	}

  getLogrus().Infof(format, value)
}
func Warn(value interface{}, args ...string) {
	format := DEFAULT_FORMAT
	if len(args) > 0 {
		format = args[0]
	}

  getLogrus().Warnf(format, value)
}
func Error(value interface{}, args ...string) {
	format := DEFAULT_FORMAT
	if len(args) > 0 {
		format = args[0]
	}

  getLogrus().Errorf(format, value)
}

func Panic(value interface{}, args ...string) {
	format := DEFAULT_FORMAT
	if len(args) > 0 {
		format = args[0]
	}

  getLogrus().Panicf(format, value)
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
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
