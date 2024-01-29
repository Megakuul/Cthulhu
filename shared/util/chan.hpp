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

#ifndef CHAN_H
#define CHAN_H

#include <condition_variable>
#include <mutex>
#include <queue>
#include <utility>

using namespace std;

/**
 * Chan is a simple wrapper around std::queue
 * that allows asynchron access and waiting for new values
 * 
 * Every operation / method can be used asynchron without any synchronisation mechanism
 */
template <typename T>
class chan {
public:
	virtual ~chan() {
		unique_lock<mutex> lock(chanMutex);
		chanCond.notify_all();
		readerCond.wait(
	};
	/**
	 * Push a value to the chan
	 */ 
	void push(T val) {
		lock_guard<mutex> lock(chanMutex);
		chanQueue.push(val);
		chanCond.notify_one();
	};

	/**
	 * Get next value from the chan
	 *
	 * If no value is in chan, this will suspend the thread
	 * and block execution until the next value is pushed.
	 *
	 * Important: If you use multiple get() that wait at the same time,
	 * the thread which gets informed first is "randomly" determined by the OS thread handler.
	 */
	pair<T, bool> get() {
		unique_lock<mutex> lock(chanMutex);
		readerThreadCount++;
		chanCond.wait(lock, [this]{ return isChanShut || !chanQueue.empty(); });
		if (isChanShut) {
			readerThreadCount--;
		  readerCond.notify_one();
			return make_pair(nullptr, false);
		}
		T el = move(chanQueue.front());
		chanQueue.pop();
		return make_pair(el, true);
	};

	/**
	 * Get size of the chan
	 *
	 * Important: This operation is not constant and has a small overhead due to a mutex locking
	 */
	int size() {
		lock_guard<mutex> lock(chanMutex);
		return chanQueue.size();
	};
	
private:
	bool isChanShut = false;
	
	queue<T> chanQueue;
	mutex chanMutex;
	condition_variable chanCond;

	int readerThreadCount = 0;
	condition_variable readerCond;
};

#endif
