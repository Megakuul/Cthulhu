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

namespace util {
	/**
	 * Chan is a simple wrapper around std::queue
	 * that allows asynchron access and waiting for new value
	 *
	 * Chan behaves really simular to a Go chan.
	 * 
	 * Every operation / method can be used asynchron without any synchronisation mechanism.
	 *
	 * Close the channel with close() or call the destructor.
	 */
	template <typename T>
	class chan {
	public:
		virtual ~chan() {
			// If channel was not shut, shut it now. Note that this is more like a preventFootGun() function
			// its recommended to close the channel in a controlled manner with close();
			close();
		};
	
		/**
		 * Push a value to the chan
		 *
		 * If the channel is already closed, it will do nothing.
		 */ 
		void push(T val) {
			lock_guard<mutex> lock(chanMutex);
			// If channel is shut don't allow anything to be pushed to the queue
			if (isChanShut) return;
			// Push value to the queue and notify the chanCond to update one random reader thread
			// This is the same behavior as you will see in Go channels.
			chanQueue.push(val);
			chanCond.notify_one();
		};

		/**
		 * Get next value from the chan
		 *
		 * If no value is in chan, this will suspend the thread
		 * and block execution until the next value is pushed.
		 *
		 * This will return a pair, which always returns the `value` and a `ok` parameter.
		 * If everything is fine, the `value` is filled and `ok` is `true`.
		 * If the channel was closed, `value` is `T()` and `ok` is `false`.
		 *
		 * Important: If you use multiple get() that wait at the same time,
		 * the thread which gets informed is "randomly" determined by the OS thread handler.
		 */
		pair<T, bool> get() {
			unique_lock<mutex> lock(chanMutex);
			// If get is called while ChanShut, it must be catched here
			// because if not wait will wait forever (as notify_all() was already called at this point)
			if (isChanShut) return make_pair(T(), false);

			// Increment reader count
			readerThreadCount++;
			// Wait for a chanCond notification, this happens in 2 scenarios, 1. Something is pushed 2. channel is shut
			chanCond.wait(lock, [this]{ return isChanShut || !chanQueue.empty(); });
			// Decrement reader count
			readerThreadCount--;
			// If channel is shut, notify closer to check the readerCount
			if (isChanShut) {
				readerCond.notify_one();
				return make_pair(T(), false);
			}
			// If something is pushed, pop it from queue and return it
			T el = move(chanQueue.front());
			chanQueue.pop();
			return make_pair(el, true);
		};

		/**
		 * Close the channel
		 *
		 * Closing the channel will send a notification to all threads where a reader is waiting (with *get()*).
		 * All readers will then return <T(), false> to indicate the channel has closed (simular to a go chan).
		 *
		 * The channel is also closed if the object is destructed.
		 *
		 * Closing the channel and keeping the object alive will not blow your leg off. This is a controlled operation.
		 */
		void close() {
			unique_lock<mutex> lock(chanMutex);
			if (isChanShut) return;
			isChanShut = true;
			chanCond.notify_all();
			readerCond.wait(lock, [this]{ return readerThreadCount<=0; });
		};

		/**
		 * Returns the state of the channel
		 */
		bool isclosed() {
			lock_guard<mutex> lock(chanMutex);
			return isChanShut;
		};

		/**
		 * Get size of the chan
		 *
		 * Important: This operation is not constant and has a small overhead due to a mutex locking
		 */
		int size() {
			lock_guard<mutex> lock(chanMutex);
			if (isChanShut) return 0;
			return chanQueue.size();
		};
	
	private:
		// Determines the state of the channel
		bool isChanShut = false;

		// Underlying FIFO datastructur
		queue<T> chanQueue;
		// Lock for any operation in chan
		mutex chanMutex;
		// Variable for notifying readers if state changed or something is pushed to the structure
		condition_variable chanCond;

		// Count of waiting readers (suspended threads waiting to be notified)
		int readerThreadCount = 0;
		// Variable for notifying the channel after they have been shut
		condition_variable readerCond;
	};
}

#endif
