package log

import (
	"fmt"
	"regexp"
	"runtime"
  "github.com/sirupsen/logrus"
)


func Info(value interface{}, args ...string) {
	// Default format is "%s"
	format := "%s"
	
	// If a format is provided in args, use it
	if len(args) > 0 {
		format = args[0]
	}

	// Get file and line number
	_, file, line, _ := runtime.Caller(1)
	reg := "[/]"
	files := regexp.MustCompile(reg).Split(file, -1)

	// Print file name and line number
	fmt.Printf("%s %d| ", files[len(files)-1], line)
	
	// Print the value using the specified or default format
	fmt.Printf(format, value)
	fmt.Println()
}
