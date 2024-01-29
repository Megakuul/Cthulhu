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

#include <string>
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
		string runtimeinfo;
		LOGLEVEL loglevel;
	}; 

	class Logger {
	public:
		Logger(LOGLEVEL logLevel, string logPath, bool logToStd, bool logDebug, int logQueueThreshold) {
			
		}
		virtual ~Logger() {

		}

	private:
		LOGLEVEL logLevel;
		bool logToStd;
		bool logDebug;
		int logChanThreshold;
		chan<LogMessage> logChan;
	};
};

#endif
