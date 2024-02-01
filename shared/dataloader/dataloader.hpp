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

#ifndef DATALOADER_H
#define DATALOADER_H

#include <stdexcept>
#include <string>
#include <unordered_map>
#include <fstream>
#include <vector>
#include <shared_mutex>
#include <unordered_map>

using namespace std;

namespace dataloader {

	class MetaConfig {
	public:
		string GetStoragePath() {
			shared_lock<shared_mutex> lock(metaLock);
			return storagePath;
		};

		vector<string> GetClusterNodeAddr() {
			shared_lock<shared_mutex> lock(metaLock);
			return clusterNodeAddr;
		};

		void FetchData() {
			unique_lock<shared_mutex> lock(metaLock);
			// Implement parser
		};

	private:
		shared_mutex metaLock;
		string storagePath;
		vector<string> clusterNodeAddr;

		/**
		 * Parse configuration into a <string,string> map
		 *
		 * Uses a custom parser, example config:
		 *
		 * ```
		 *
		 * # I'm a comment until newline
		 * somekey="some.value;9?
		 * I can contain spaces, tabs, newlines
		 * "uglyplacedkey="I'm valid too"
		 *
		 * wellplacedkey=""
		 * / I'm also a comment until newline
		 * ```
		 *
		 * If a key is placed multiple times, only the first one is evaluated
		 */
		unorder_map<string, string> readConfiguration(string path) {
			unordered_map<string, string> mapBuffer;
			// Read config file
			ifstream file(path);
			if (!file.is_open()) {
				throw runtime_error("Failed to open config file at: " + path);
			}

			// Char buffer
			char c;
			// Keeps track of lines for debug messages
			int lineCount = 0;
			// Current key buffer
			string curKey;
			// Current value buffer
			string curVal;

			// Iterate over chars
			while (file.get(c)) {
				// Skip newline
				if (c=='\n') {
					lineCount++;
					continue;
				}
				// Skip space, tab
				if (isspace(c)) {
					continue;
				}
				// # | / indicate a comment
				if (c=='#'||c=='/') {
					// Skip til EOF or newline
					while (file.get(c)&&c!='\n');
					lineCount++;
					continue;
				}

				// Eat key
				curKey = "";
				do {
					curKey+=c;
					// EOF or newline in key is not allowed
					if (!file.get(c)||c=='\n') {
						throw runtime_error(
						  "Failed to parse config file at: "
						  + path + "\n"
						  + "Unexpected EOF or newline on line: " + to_string(lineCount)
					  );
					}
					// Read until '=' char
				} while (c!='=');

				// Read next char which is expected to be '"'
				if (!file.get(c)||c!='"') {
					throw runtime_error(
					  "Failed to parse config file at: "
					  + path + "\n"
					  + "Expected '\"' after '=' on line: " + to_string(lineCount)
				  );
				}
				
				// Eat value
				curVal = "";
				while (true) {
					// EOF is not expected in value, every other char can be used
					if (!file.get(c)) {
						throw runtime_error(
						  "Failed to parse config file at: "
						  + path + "\n"
						  + "Unexpected EOF on line: " + to_string(lineCount)
					  );
						// Read until '"' char 
					} else if (c=='"') break;
					// Add linecount
					if (c=='\n') lineCount++;
					curVal+=c;
				};
				// Use insert, first key inserted is valid, other same keys are invalidated
				mapBuffer.insert({curKey, curVal});
			}
			return mapBuffer;
		}
		
		// The functions below specify how to parse each value from string
		string parseStoragePath(string originalValue) {
			return originalValue;
		};
		vector<string> parseClusterNodeAddr(string originalValue) {

		};
	};

	
}
#endif
