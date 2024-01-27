/**
 * Cthulhu System
 *
 * Copyright (C) 2024  Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type LOGLEVEL int
const (
	ERROR LOGLEVEL = iota
	WARN
	INFO
)

type LogMessage struct {
	message string
	runtimeinfo string
	loglevel LOGLEVEL
}

type Logger struct {
	logLevel LOGLEVEL
	logFile *os.File
	logToStd bool
	logDebug bool
	logChan chan *LogMessage
}

var GlobLogger Logger

func InitLogger(logLevel LOGLEVEL, logPath string, logToStd bool, logDebug bool, logQueueThreshold int) error {
	// Create Logfile path if not existent
	logPathParent, _ := filepath.Split(logPath)
	if err := os.MkdirAll(logPathParent, 0755); err!=nil {
		return err
	}
	
	var logger Logger
	var err error
	logger.logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err!=nil {
		return err
	}

	logger.logToStd = logToStd
	logger.logDebug = logDebug
	logger.logLevel = logLevel
	logger.logChan = make(chan *LogMessage, logQueueThreshold)

	logger.startLogWorker()
	
	return err
}

func (l* Logger) CloseLogger() {
	l.closeLogWorker()
	l.logFile.Close()
}

func (l* Logger) LogError(msg string) {
	runtimeinfo := ""
	if l.logDebug {
		runtimeinfo = l.getRuntimeInfo()
	}
	l.logChan<-&LogMessage{msg, runtimeinfo, ERROR}
}

func (l* Logger) LogWarn(msg string) {
	if l.logLevel>ERROR {
		runtimeinfo := ""
		if l.logDebug {
			runtimeinfo = l.getRuntimeInfo()
		}
		l.logChan<-&LogMessage{msg, runtimeinfo, WARN}
	}
}

func (l* Logger) LogInfo(msg string) {
	if l.logLevel>WARN {
		runtimeinfo := ""
		if l.logDebug {
			runtimeinfo = l.getRuntimeInfo()
		}
		l.logChan<-&LogMessage{msg, runtimeinfo, INFO}
	}
}

func (l* Logger) getRuntimeInfo() string {
	runtimeinfo := "[ RUNTIME INFORMATION ]:\n"
	_, file, line, ok := runtime.Caller(3)
	if ok {
		runtimeinfo += fmt.Sprintf("|-[ LOG CALLER STACK ]: Line (%d) File (%s)\n", line, file)
	}
	return runtimeinfo
}

func (l* Logger) log(msg *LogMessage) {
	switch msg.loglevel {
	// TODO: Parse and write message to the output
	case ERROR:

	case WARN:
	case INFO:
	}
}


func (l* Logger) startLogWorker() {
	for {
		select {
		case msg, ok := <-l.logChan:
			if ok {
				l.log(msg)
			} else {
				// Exit if channel was closed
				return
			}
		}
	}
}

func (l* Logger) closeLogWorker() {
	// Close channel which will cause the logworker to exit
	close(l.logChan)
}
