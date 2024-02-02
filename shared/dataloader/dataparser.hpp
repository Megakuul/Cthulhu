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

#ifndef DATAPARSER_H
#define DATAPARSER_H

#define TMP_FILE_EXTENSION ".tmp"

#include <string>
#include <fstream>
#include <filesystem>
#include <unordered_map>
#include <shared_mutex>
#include <mutex>

using namespace std;

namespace dataloader {

	/**
	 * Object holding a inmem configuration
	 *
	 * Configuration can be written and read from/to disk
	 *
	 * All operations that are fully thread-safe (synchronized).
	 *
	 * Uses a custom parser, that parses a kind of a key-value config file (example):
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
	 */
	class ConfigParser {
	public:
		ConfigParser(string path) {
			configPath = path;
		};

		/**
		 * Get value to specific key
		 *
		 * This operation does not read / parse anything from disk!
		 */
		string GetValue(string &key) {
			shared_lock<shared_mutex> confLock(configLock);
			auto it = config.find(key);
			if (it != config.end()) return it->second;
			else return "";
		};

		/**
		 * Get full parsed configuration object
		 *
		 * This operation does not read / parse anything from disk!
		 */
		unordered_map<string, string> GetConfig() {
			shared_lock<shared_mutex> confLock(configLock);
			return config;
		};

		/**
		 * Set value to specific key
		 *
		 * This operation does not write anything to disk!
		 */
		void SetValue(string &key, string &value) {
			unique_lock<shared_mutex> confLock(configLock);
			config[key] = value;
		}

		/**
		 * Set configuration
		 *
		 * This operation does not write anything to disk!
		 */
		void SetConfig(unordered_map<string, string>& map) {
			unique_lock<shared_mutex> confLock(configLock);
			config = map;
		};
		
		/**
		 * Read and Parse configuration directly from disk to inmem config
		 *
		 * If a key is placed multiple times, only the first one is evaluated
		 *
		 * Function will throw a runtime error if it fails
		 */
	  void ReadFromDisk() {
			// Read lock the file config lock
			shared_lock<shared_mutex> fileLock(configFileLock);
			
			unordered_map<string, string> mapBuffer;
			// Read config file
			ifstream file(configPath);
			if (!file.is_open()) {
				throw runtime_error("Failed to open config file at: " + configPath);
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
							+ configPath + "\n"
							+ "Unexpected EOF or newline on line: " + to_string(lineCount)
						);
					}
					// Read until '=' char
				} while (c!='=');

				// Read next char which is expected to be '"'
				if (!file.get(c)||c!='"') {
					throw runtime_error(
					 	"Failed to parse config file at: "
					  + configPath + "\n"
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
							+ configPath + "\n"
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
			// Write lock the inmem config lock
			unique_lock<shared_mutex> confLock(configLock);
			config = mapBuffer;
		}
		
		/**
		 * Writes inmem configuration directly to disk
		 *
		 * Function will throw a runtime error if it fails
		 */ 
		void WriteToDisk() {
			// Write lock the file config lock
			unique_lock<shared_mutex> fileLock(configFileLock);
			// Read lock the inmem config lock
			shared_lock<shared_mutex> confLock(configLock);
			
			// Open tmp config file
			ofstream file(configPath + TMP_FILE_EXTENSION, ofstream::out | ofstream::trunc);
			if (!file.is_open()) {
				throw runtime_error("Failed to open config file at: " + configPath);
			}

			// Insert deparsed configuration
			file << "# Manual changes to configuration may be overwritten\n";
			file << "# Consider using Meta Hook from the Cthulhu component\n";
			for (auto kv : config) {
				file << kv.first << "=" << "\"" << kv.second << "\"" << "\n";
			}
			file << "# End of config" << endl;

			// Close stream manually, because it will be accessed by the fs::rename below
			file.close();

			// Move tmp config to config
			// This prevents file corruption on unexpected application crashes (e.g. shutdown while writing).
			filesystem::rename(configPath + TMP_FILE_EXTENSION, configPath);
		}

	private:
		// Mutex lock for the configuration file
		shared_mutex configFileLock;
		// Mutex lock for the inmem configuration
		shared_mutex configLock;
		// Path of the configuration
		string configPath;
		// In memory configuration object
		unordered_map<string, string> config;
	};
}

#endif