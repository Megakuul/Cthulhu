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

using namespace std;

template <typename T>
class chan {
public:
	void push(T val) {
		lock_guard<mutex> lock(chanMutex);
		chanQueue.push(val);
		chanCond.notify_one();
	};

	T read() {
		unique_lock<mutex> lock(chanMutex);
		chanCond.wait(lock, [this]{ return !chanQueue.empty(); });
		T el = move(chanQueue.front());
		chanQueue.pop();
		return el;
	};
	
private:
	queue<T> chanQueue;
	mutex chanMutex;
	condition_variable chanCond;
};

#endif
