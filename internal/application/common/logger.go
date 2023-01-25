package common

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var lgr *Logger

func GetLogger() *Logger {
	if lgr == nil {
		lgr = &Logger{
			buffer: bytes.NewBuffer([]byte{}),
		}
	}

	return lgr
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
		"\nBrazil: %s\nPST: %s\n%s",
		time.Now().Format(time.RFC3339),
		time.Now().In(pst).Format(time.RFC3339),
		message,
	)

	log.Println(finalMessage)
	l.buffer.Write([]byte(finalMessage))
}

func (l Logger) RegisterLogs() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	logsDir := basepath + "/../../../logs"
	fileName := fmt.Sprintf(logsDir+"/%d.txt", time.Now().UnixMilli())

	os.Mkdir(logsDir, 0777)
	err := os.WriteFile(fileName, lgr.GetBuffer(), 0644)
	if err != nil {
		GetLogger().Log(err.Error())
	}

	log.Println("Logs saved successfully.")
}

func (l Logger) GetBuffer() []byte {
	return l.buffer.Bytes()
}
