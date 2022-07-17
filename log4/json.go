package log4

import (
	"encoding/json"
	"log"
)

func JSON(data interface{}) {
	log.Println(shortfile(2), "-------------------------------------------")
	logger := []interface{}{shortfile(2)}
	b, _ := json.Marshal(data)
	logger = append(logger, string(b))
	logger = append(logger, "\x1b[0m")
	log.Println(logger...)
}
