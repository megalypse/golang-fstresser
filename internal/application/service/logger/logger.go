package logger

import (
	"bytes"
	"fmt"
	"log"
)

type Logger struct {
	buffer bytes.Buffer
}

func (l Logger) Log(message string) {
	finalMessage := fmt.Sprintf("%s\n", message)
	log.Println(finalMessage)

	l.buffer.Write([]byte(finalMessage))
}

func (l Logger) GetBuffer() []byte {
	return l.buffer.Bytes()
}
