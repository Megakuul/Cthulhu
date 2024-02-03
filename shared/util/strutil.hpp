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

#ifndef STRUTIL_H
#define STRUTIL_H

#include <algorithm>
#include <sstream>
#include <vector>
#include <string>

using namespace std;

/**
 * Strutil provides a set of std::string utilities
 * that would otherwise produce a lot of unnecessary boilerplate code.
 */
namespace util::strutil {
	/**
	 * Split a string into a list based on a delimiter
	 *
	 * Empty elements are omitted (e.g. "some,stuff,," -> ["some", "stuff"])
	 */
	inline vector<string> split(const string& str, const char& delimiter) {
		vector<string> tokens;
		string tokBuf;
		istringstream tokStream(str);

		while (getline(tokStream, tokBuf, delimiter)) {
			if (!tokBuf.empty())
				tokens.push_back(tokBuf);
		}
		return tokens;
	}
	
	/**
	 * Unsplit a list into a string with delimiters
	 *
	 * Unsplit will end the string with a delimiter (e.g. ["some", "stuff"] -> "some,stuff,")
	 */
	inline string unsplit(const vector<string>& tokens, const char& delimiter) {
		string str;
		for (auto tokBuf : tokens) {
		  str+=tokBuf;
			str+=delimiter;
		}
		return str;
	}
	
	/**
	 * Compare two strings without case sensitivity
	 *
	 * SaLaD == sAlAd # true
	 * bowl == bread # false
	 */
	inline bool cmpIgnoreCase(const string& str1, const string& str2) {
		if (str1.length()!=str2.length()) return false;

		return equal(str1.begin(), str1.end(), str2.begin(), str2.end(),
								 [](char a, char b) {
									 return tolower(a) == tolower(b);
								 });
	}
}

#endif
