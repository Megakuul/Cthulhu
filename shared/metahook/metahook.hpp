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

#ifndef METAHOOK_H
#define METAHOOK_H

#include <string>
#include <unordered_map>
#include <functional>

#include "shared/metaconfig/metaconfig.hpp"

using namespace std;

namespace metahook {

	/**
	 * Structure which holds function definitions for specific MetaConfig fields
	 *
	 * The hook function callback is called when the API is called to change the specified MetaConfig field.
	 *
	 * Every hook is executed synchroniously, make sure they do not use cost-intensive IO operations.
	 *
	 * Hooks are expected return after the system
	 * is in a state where the updated field is fully operational.
	 */
	struct UpdateHooks {
		// Hooks for string fields
		unordered_map<string, function<void(string, string)>> StringFieldHooks;
		// Hooks for bool fields
		unordered_map<string, function<void(string, bool)>> BoolFieldHooks;
		// Hooks for double fields
		unordered_map<string, function<void(string, double)>> DoubleFieldHooks;
		// Hooks for list fields
		unordered_map<string, function<void(string, vector<string>)>> ListFieldHooks;
	};

	/**
	 * MetaHook is a component to update the MetaConfiguration
	 * over a controlled HTTP API
	 *
	 * It uses a updateHookMap to specify callback functions
	 * for specific MetaConfig keys, those callbacks can be used
	 * to live-update the configuration in the components.
	 *
	 * MetaHook launches a HTTP API over a UNIX socket on the specified location
	 *
	 * Main purpose for this API is that infrastructure controllers like juju
	 * can manage the MetaConfig at runtime.
	 */
	class MetaHook {

	private:
		metaconfig::MetaConfig* metaConfig;
		UpdateHooks updateHooks;
		string socketPath;
		FileMod socketPerm;
		
	};
}

#endif
