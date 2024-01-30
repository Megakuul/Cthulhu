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
	"time"
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
	debuginfo string
	loglevel LOGLEVEL
}

type Logger struct {
	logLevel LOGLEVEL
	logFile *os.File
	logToStd bool
	logDebug bool
	logChanThreshold int
	logChan chan *LogMessage
}

func InitLogger(logLevel LOGLEVEL, logPath string, logToStd bool, logDebug bool, logQueueSize int8) error {
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
	// Queue threshold is set to 50%. If it goes beyond, this is already very critical
	logger.logChanThreshold = int(logQueueSize) / 2
	logger.logChan = make(chan *LogMessage, logQueueSize)

	logger.startLogWorker()
	
	return err
}

func (l* Logger) CloseLogger() {
	l.closeLogWorker()
	l.logFile.Close()
}

func (l* Logger) LogError(msg string) {
	debuginfo := ""
	if l.logDebug {
		debuginfo = l.getDebugInfo(2)
	}
	l.logChan<-&LogMessage{msg, debuginfo, ERROR}
}

func (l* Logger) LogWarn(msg string) {
	if l.logLevel>ERROR {
		debuginfo := ""
		if l.logDebug {
			debuginfo = l.getDebugInfo(2)
		}
		l.logChan<-&LogMessage{msg, debuginfo, WARN}
	}
}

func (l* Logger) LogInfo(msg string) {
	if l.logLevel>WARN {
		debuginfo := ""
		if l.logDebug {
			debuginfo = l.getDebugInfo(2)
		}
		l.logChan<-&LogMessage{msg, debuginfo, INFO}
	}
}

func (l* Logger) getDebugInfo(stackdepth int) string {
	debuginfo := "[ RUNTIME INFORMATION ]:\n"
	// Get stack information from the callerstack + stackdepth
	_, file, line, ok := runtime.Caller(stackdepth+1)
	if ok {
		debuginfo += fmt.Sprintf("|-[ LOG CALLER STACK ]: Line (%d) File (%s)\n", line, file)
	}
	return debuginfo
}

func (l* Logger) log(msg *LogMessage) {
	outstr := time.Now().Format("\n[ 05:04:15 - 02.01.2006 ]\n")
	switch msg.loglevel {
	case ERROR:
		outstr += "[ ERROR ]:\n"
		outstr += msg.message
		outstr += "\n"
		outstr += msg.debuginfo
		outstr += "\n"
		l.logFile.Write([]byte(outstr))
		if l.logToStd {
			os.Stderr.Write([]byte(outstr))
		}
	case WARN:
		outstr += "[ WARNING ]:\n"
		outstr += msg.message
		outstr += "\n"
		outstr += msg.debuginfo
		outstr += "\n"
		l.logFile.Write([]byte(outstr))
		if l.logToStd {
			os.Stderr.Write([]byte(outstr))
		}
	case INFO:
		outstr += "[ INFORMATION ]:\n"
		outstr += msg.message
		outstr += "\n"
		outstr += msg.debuginfo
		outstr += "\n"
		l.logFile.Write([]byte(outstr))
		if l.logToStd {
			os.Stdout.Write([]byte(outstr))
		}
	}
}


func (l* Logger) startLogWorker() {
	for {
		select {
		case msg, ok := <-l.logChan:
			if ok {
				if len(l.logChan) > l.logChanThreshold {
					l.log(&LogMessage{
						"Log Queue is under high pressure!",
						l.getDebugInfo(1),
						WARN,
					})
				}
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
