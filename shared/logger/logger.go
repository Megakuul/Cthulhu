package logger

import (
	"fmt"
)

var GlobLogger Logger

type LOGLEVEL int
const (
	ERROR LOGLEVEL = iota
	WARN
	INFO
	DEBUG
)

type Logger struct {
	logLevel LOGLEVEL
	logPath string
	logSize int
	logChan chan *logMessage
}

type logMessage struct {
	message string
	loglevel LOGLEVEL
}

func InitLogger(logLevel LOGLEVEL, logPath string, logSize int, logQueueThreshold int) {
	var logger Logger
	logger.logLevel = logLevel
	// TODO: Test path
	logger.logPath = logPath
	logger.logSize = logSize
	logger.logChan = make(chan *logMessage, logQueueThreshold)
}

func (l* Logger) LogError(msg string) {
	l.logChan<-&logMessage{msg, ERROR}
}

func (l* Logger) LogWarn(msg string) {
	l.logChan<-&logMessage{msg, WARN}
}

func (l* Logger) LogInfo(msg string) {
	l.logChan<-&logMessage{msg, INFO}
}

func (l* Logger) LogDebug(msg string) {
	l.logChan<-&logMessage{msg, ERROR}
}

func (l* Logger) StartLogWorker() {
	for {
		select {
		case msg, ok := <-l.logChan:
			if ok {
				// TODO: Build logic like forming the output (maybe also mutex locking btw)
				// + Adding timestamp
				fmt.Println(msg.message)
			} else {
				return
			}
		}
	}
}

func (l* Logger) CloseLogWorker() {
	close(l.logChan)
}
