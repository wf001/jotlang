package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

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
func GREEN(body string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", body)
}

func DebugTokens(tok *types.Token) {
	Debug(BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		Debug(BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}

func DebugNode(node *types.Node, depth int) {
	Debug(BLUE(fmt.Sprintf("%s %p %#+v %#+v", strings.Repeat("\t", depth), node, node.Kind, node.Val)))

	for ; node != nil; node = node.Next {
		switch node.Kind {
		case types.ND_ADD:
			DebugNode(node.Lhs, depth+1)
			DebugNode(node.Rhs, depth+1)
		}
	}
}

func DebugMessage(message string) {
	Debug("", GREEN(message))
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
	SetLevelWarning()
	SetFormatter()
}

func SetOutputFile(fileName string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(f)
	defer f.Close()
}

func SetFormatter() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
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
