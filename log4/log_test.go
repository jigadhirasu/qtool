package log4

import (
	"log"
	"testing"
)

func TestLog4ding(t *testing.T) {
	Error("%s", "string")
	Warning("%s", "string")
	Success("%s", "string")
	Info("string")
	Debug("string")

	Error("abc", "ced", 123, map[string]string{"x": "y"})
	Error("%d%s", 1, "abcb")

	log.Println("Hello World")
}
