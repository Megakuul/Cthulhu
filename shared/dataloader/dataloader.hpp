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

#include <string>
#include <vector>
#include <shared_mutex>

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
	};

	
}
#endif
