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

#ifndef LOGGER_H
#define LOGGER_H

#include <chrono>
#include <ctime>
#include <filesystem>
#include <iomanip>
#include <stdexcept>
#include <fstream>
#include <string>
#include <format>
#include <iostream>
#include <thread>

#include "shared/util/chan.hpp"

using namespace std;

namespace logger {

	enum LOGLEVEL {
		ERROR = 1,
		WARN = 2,
		INFO = 3,
	};

	struct LogMessage {
		string message;
		string debuginfo;
		LOGLEVEL loglevel;
	}; 

	class Logger {
	public:
		Logger(LOGLEVEL logLevel, string logPath, bool logToStd, bool logDebug, int logQueueThreshold) {
			// Create Logfile path if not existent
			filesystem::create_directories(filesystem::path(logPath).parent_path());
			// Before logger is initalized, errors are just thrown to top level
			logFile.open(logPath, ios::out | ios::app);
			if (!logFile.is_open()) {
				throw runtime_error("Failed to open logfile at: " + logPath);
			}

			this->logToStd = logToStd;
			this->logDebug = logDebug;
			this->logLevel = logLevel;
			// Queue threshold is halved to mimic the Go logger behavior, the C++ channel can grow indefinitely.
			// The threshold is only for defining a warn-threshold after which a warning is emitted.
			this->logChanThreshold = logQueueThreshold / 2;

			startLogWorker();
		}
		virtual ~Logger() {
			closeLogWorker();
			logFile.close();
		}

		/**
		 * Log an error
		 *
		 * FILE and LINE parameters are debuginfos and must be set to the respective C++ macros
		 * at the stack where you want to get the information from.
		 *
		 * FILE=__FILE__
		 * LINE=__LINE__
		 */
		void LogError(string msg, const string FILE, const int LINE) {
			string debuginfo;
			if (logDebug) {
				debuginfo = getDebugInfo(FILE, LINE);
			}
			logChan.push({msg, debuginfo, ERROR});
		}

		/**
		 * Log a warning
		 *
		 * FILE and LINE parameters are debuginfos and must be set to the respective C++ macros
		 * at the stack where you want to get the information from.
		 *
		 * FILE=__FILE__
		 * LINE=__LINE__
		 */
		void LogWarn(string msg, const string FILE, const int LINE) {
			if (logLevel>ERROR) {
				string debuginfo;
				if (logDebug) {
					debuginfo = getDebugInfo(FILE, LINE);
				}
				logChan.push({msg, debuginfo, WARN});
			}
		}

		/**
		 * Log a information
		 *
		 * FILE and LINE parameters are debuginfos and must be set to the respective C++ macros
		 * at the stack where you want to get the information from.
		 *
		 * FILE=__FILE__
		 * LINE=__LINE__
		 */
		void LogInfo(string msg, const string FILE, const int LINE) {
			if (logLevel>WARN) {
				string debuginfo;
				if (logDebug) {
					debuginfo = getDebugInfo(FILE, LINE);
				}
				logChan.push({msg, debuginfo, INFO});
			}
		}

	private:
		LOGLEVEL logLevel;
		ofstream logFile;
		bool logToStd;
		bool logDebug;
		int logChanThreshold;
		util::chan<LogMessage> logChan;
		// C++ IO operations like writes to ofstream are not thread-safe
		// This lock synchronizes every io operation (writes to stdout / disk).
		mutex ioMutex;
		
		/**
		 * Get and format debuginformation
		 *
		 * FILE and LINE parameters are debuginfos and must be set to the respective C++ macros
		 * at the stack where you want to get the information from.
		 *
		 * FILE=__FILE__
		 * LINE=__LINE__
		 */
		string getDebugInfo(const string FILE, const int LINE) {
			string debuginfo = "[ RUNTIME INFORMATION ]:\n";
			debuginfo += format("|-[ LOG CALLER STACK ]: Line ({}) File ({})\n", LINE, FILE);
			return debuginfo;
		}

		/**
		 * Writes a message directly to the log (and optionally to std)
		 *
		 * Synchronisation is handled through the ioMutex
		 */
		void log(const LogMessage &msg) {
			stringstream outstream;
			string outstr;
			
			time_t now = chrono::system_clock::to_time_t(chrono::system_clock::now());
			outstream << put_time(localtime(&now), "\n[ %H:%M:%S - %d.%m.%Y ]\n");

			lock_guard<mutex> lock(ioMutex);
			switch (msg.loglevel) {
			case ERROR:
				outstream << "[ ERROR ]:\n";
				outstream << msg.message << "\n";
				outstream << msg.debuginfo;
				
				outstr = outstream.str();
				logFile << outstr << endl;
				if (logToStd) {
					cerr << outstr << endl;
				}
				break;
			case WARN:
				outstream << "[ WARN ]:\n";
				outstream << msg.message << "\n";
				outstream << msg.debuginfo;
				
				outstr = outstream.str();
				logFile << outstr << endl;
				if (logToStd) {
					cerr << outstr << endl;
				}
				break;
			case INFO:
				outstream << "[ INFO ]:\n";
				outstream << msg.message << "\n";
				outstream << msg.debuginfo;
				
				outstr = outstream.str();
				logFile << outstr << endl;
				if (logToStd) {
					cout << outstr << endl;
				}
				break;
			}			
		}

		/**
		 * Start seperate thread to write logs
		 *
		 * Reads messages from the logChan
		 */
		void startLogWorker() {
			thread worker([this]() {
				pair<LogMessage, bool> res;
				while (true) {
					res = logChan.get();
					if (res.second) {
						if (logChan.size() > logChanThreshold) {
							log({
									"Log Queue is under high pressure!",
									getDebugInfo(__FILE__, __LINE__),
									WARN,
								});
						}
						log(res.first);
					} else {
						// Exit if channel was closed
						return;
					}
				}
			});

			worker.detach();
		}

		/**
		 * Stop logworker
		 *
		 * This function waits until the worker has fully exited
		 * so it can be savely used in destructor.
		 */
		void closeLogWorker() {
			logChan.close();
		}
	};
};

#endif
