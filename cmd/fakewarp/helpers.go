package main

import (
	"bytes"
	"encoding/json"
	"log"
)

func toJson(message interface{}) string {
	b, error := json.Marshal(message)
	if error != nil {
		log.Fatal(error)
	}
	buffer := new(bytes.Buffer)
	buffer.Write(b)
	buffer.Write([]byte("\n"))
	return buffer.String()
}
