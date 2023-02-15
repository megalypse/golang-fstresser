package common

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var lgrMap map[string]*Logger

func init() {
	lgrMap = make(map[string]*Logger)
}

func GetLogger(ctx context.Context) *Logger {

	profileNameRaw := ctx.Value(GetCtxKey("profile-name"))

	if profileNameRaw != nil {
		return getLogger(profileNameRaw.(string))
	}

	return getLogger("default")
}

func getLogger(lgrName string) *Logger {
	lgr, ok := lgrMap[lgrName]

	if !ok {
		lgr = getNewLogger(lgrName)
		lgrMap[lgrName] = lgr
	}

	return lgr
}

func getNewLogger(lgrName string) *Logger {
	return &Logger{
		buffer:  bytes.NewBuffer([]byte{}),
		lgrName: lgrName,
	}
}

type Logger struct {
	buffer  *bytes.Buffer
	lgrName string
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

	log.Println(finalMessage + "\n")
	l.buffer.Write([]byte(finalMessage))
}

func (l Logger) SilentLog(message string) {
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

	l.buffer.Write([]byte(finalMessage))
}

func (l Logger) RegisterLogs() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	logsDir := basepath + "/../../../logs"

	profileName := strings.ReplaceAll(l.lgrName, " ", "_")
	fileName := fmt.Sprintf(logsDir+"/%s_%d.txt", profileName, time.Now().UnixMilli())

	os.Mkdir(logsDir, 0777)
	err := os.WriteFile(fileName, l.GetBuffer(), 0644)
	if err != nil {
		l.Log(err.Error())
	}

	log.Println("Logs saved successfully.")
}

func (l Logger) GetBuffer() []byte {
	return l.buffer.Bytes()
}
