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

// lgrMap holds a logger instance reference for each profile
var lgrMap map[string]*Logger

func init() {
	lgrMap = make(map[string]*Logger)
}

// GetLogger provides the logger for the profile that requested it.
// If the context does not contain any profile name, then the default logger will be provided.
func GetLogger(ctx context.Context) *Logger {

	profileNameRaw := ctx.Value(GetCtxKey("profile-name"))

	if profileNameRaw != nil {
		return getLogger(profileNameRaw.(string))
	}

	return getLogger("default")
}

// getLogger searches for the logger by the provided key,
// if no key is provided, a new logger is created, saved and returned.
func getLogger(lgrName string) *Logger {
	lgr, ok := lgrMap[lgrName]

	if !ok {
		lgr = makeNewLogger(lgrName)
		lgrMap[lgrName] = lgr
	}

	return lgr
}

// makeNewLogger creates a new Logger instance with the provided lgrName
func makeNewLogger(lgrName string) *Logger {
	return &Logger{
		buffer:  bytes.NewBuffer([]byte{}),
		lgrName: lgrName,
	}
}

// Logger to be used to provide visual feedback for the users.
type Logger struct {
	// buffer holds all the logged information so it can be registered later
	buffer *bytes.Buffer

	// lgrName the name provided to an instance of Logger
	lgrName string
}

// Log logs the given message along with a timestamp and also saves
// the final message inside `l.buffer` for later registration.
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

// SilentLog does not give feedback to the final user and only saves the final message
// inside `l.logger`.
//
// This may be useful when saving error messages without the need to spam the console.
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

/*
RegisterLogs uses `l.buffer` to save all the messages logged inside a .txt file
*/
func (l Logger) RegisterLogs() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// TODO: get logs path by env variable
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
