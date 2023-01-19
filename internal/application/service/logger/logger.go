package logger

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

func NewLogger() Logger {
	return Logger{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

type Logger struct {
	buffer *bytes.Buffer
}

func (l Logger) Log(message string) {
	pst, err := time.LoadLocation("America/Los_Angeles")

	if err != nil {
		log.Fatal(err.Error())
	}

	finalMessage := fmt.Sprintf(
		"\nBrazil: %s\nPST: %s%s",
		time.Now().Format(time.RFC3339),
		time.Now().In(pst).Format(time.RFC3339),
		message,
	)

	log.Println(finalMessage)
	l.buffer.Write([]byte(finalMessage))
}

func (l Logger) GetBuffer() []byte {
	return l.buffer.Bytes()
}
