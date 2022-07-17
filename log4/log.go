package log4

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const (
	gray = uint8(iota + 90)
	red
	green
	yellow
	blue
	purple
	sky
	white
)

var level = 0

func init() {
	// log.SetFlags(log.Lmicroseconds)
	// log.SetPrefix(fmt.Sprintf("\x1b[%dm", yellow))
}

func shortfile(depth int) string {
	_, longfile, line, _ := runtime.Caller(depth)
	strs := strings.Split(longfile, "/")
	short := strs[len(strs)-1]
	return fmt.Sprintf("\x1b[%dm%s:%d\x1b[0m", blue, short, line)
}

// Level show message
// 0 error warning success info debug
// 1 error warning success info
// 2 error warning success
// 3 error warning
// 4 error
func Level(lv int) {
	level = lv
}

func formatMessage(color uint8, title string, format interface{}, a ...interface{}) string {
	logger := []interface{}{shortfile(3), fmt.Sprintf("\x1b[%dm", color), "[" + title + "]"}
	if fstr, ok := format.(string); ok {
		if strings.Count(fstr, "%") == len(a) {
			logger = append(logger, fmt.Sprintf(fstr, a...), "\x1b[0m")
			log.Println(logger...)
			return fmt.Sprintln(logger...)
		}
	}
	logger = append(logger, format)
	logger = append(logger, a...)
	logger = append(logger, "\x1b[0m")
	log.Println(logger...)
	return fmt.Sprintln(logger...)
}

// Error Error
func Error(format interface{}, a ...interface{}) string {
	return formatMessage(red, "error", format, a...)
}

// Warning Warning
func Warning(format interface{}, a ...interface{}) string {
	if level > 3 {
		return ""
	}
	return formatMessage(purple, "warning", format, a...)
}

// Success Success
func Success(format interface{}, a ...interface{}) string {
	if level > 2 {
		return ""
	}
	return formatMessage(green, "success", format, a...)
}

// Info Info
func Info(format interface{}, a ...interface{}) string {
	if level > 1 {
		return ""
	}
	return formatMessage(sky, "info", format, a...)
}

// Debug Debug
func Debug(format interface{}, a ...interface{}) string {
	if level > 0 {
		return ""
	}
	return formatMessage(white, "debug", format, a...)
}
