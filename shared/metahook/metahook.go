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

package metahook

import (
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/megakuul/cthulhu/shared/metaconfig"
)

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
type MetaHook struct {
	metaConfig *metaconfig.MetaConfig
	updateHooks UpdateHooks
	socketPath string
	socketPerm fs.FileMode
	socketServer *http.Server
	socketServerMux *http.ServeMux
}

/**
 * Structure which holds function definitions for specific MetaConfig fields
 *
 * The hook function callback is called when the API is called to change the specified MetaConfig field.
 *
 * Hooks are expected to bring the system into a state where it operates like
 * the field was set at application start!
 *
 * Hooks are also expected to not return until the system is in the updated system.
 */
type UpdateHooks struct {
	// Hooks for string fields
	StringFieldHooks map[string]func(string, string) error
	// Hooks for bool fields
	BoolFieldHooks map[string]func(string, bool) error
	// Hooks for double fields
	DoubleFieldHooks map[string]func(string, float64) error
	// Hooks for list fields
	ListFieldHooks map[string]func(string, []string) error
}

/**
 * Initialize MetaHook API
 */
func CreateMetaHook(
	socketpath string,
	socketperm fs.FileMode,
	updatehooks UpdateHooks,
	config *metaconfig.MetaConfig) (*MetaHook, error) {
	
	// Create path recursively
	parentpath := filepath.Dir(socketpath)
	if err:=os.MkdirAll(parentpath, 0755); err!=nil {
		return nil, err
	}
	// Cleanup old socket
	if err:=os.Remove(socketpath); err!=nil&&!os.IsNotExist(err) {
		return nil, err
	}
	
	// Create ServeMux
	sockMux := http.NewServeMux()
	
	// Create HTTP Server
	sockSrv := &http.Server{
		Handler: sockMux,
	}

	metaHook := &MetaHook{
		config,
		updatehooks,
		socketpath,
		socketperm,
		sockSrv,
		sockMux,
	}

	// Register handlers
	sockMux.HandleFunc("/update", metaHook.updateHandler)

	return metaHook, nil
}

/**
 * Create unix socket / listener and start HTTP server
 *
 * Serve() will block execution, you can safely push it to a goroutine
 */
func (m* MetaHook) Serve() error {
	// Remove socket if already existent
	if err:=os.Remove(m.socketPath); err!=nil && !os.IsNotExist(err) {
		return err
	}
	// Create socket and open listener
	unixListener, err := net.Listen("unix", m.socketPath)
	if err!=nil {
		return err
	}
	defer unixListener.Close()
	defer os.Remove(m.socketPath)
	
	// Change socket permissions
	if err:=os.Chmod(m.socketPath, m.socketPerm); err!=nil {
		return err
	}
	// Start HTTP server
	if err:=m.socketServer.Serve(unixListener); err!=nil {
		return err
	}
	return nil
}

// Meta Handlers

type metaStringField struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type metaBoolField struct {
	Key string `json:"key"`
	Value bool `json:"value"`
}

type metaDoubleField struct {
	Key string `json:"key"`
	Value float64 `json:"value"`
}

type metaListField struct {
	Key string `json:"key"`
	Value []string `json:"value"`
}

type updateRequest struct {
	StringFields []metaStringField `json:"string_fields"`
	BoolFields []metaBoolField `json:"bool_fields"`
	DoubleFields []metaDoubleField `json:"double_fields"`
	ListFields []metaListField `json:"list_fields"`
}

type updateResponse struct {
	Err []error `json:"err"`
}

/**
 * Handler update requests
 *
 * Updates a value in the associated MetaConfig
 * and calls the updateHook for it (if defined)
 */
func (m* MetaHook) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method, expected POST!", http.StatusMethodNotAllowed)
		return
	}

	var req updateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res updateResponse
	var resMutex sync.Mutex
	var wg sync.WaitGroup
	
	// String fields
	for _,kv := range req.StringFields {
		wg.Add(1)
		go func(field metaStringField) {
			defer wg.Done()
			m.metaConfig.SetString(&field.Key, &field.Value)
			hook, exists := m.updateHooks.StringFieldHooks[field.Key]
			if exists {
				err := hook(field.Key, field.Value)
				if err!=nil {
					resMutex.Lock()
					res.Err = append(res.Err, err)
					resMutex.Unlock()
				}
			}
		}(kv)
	}

	// Bool fields
	for _,kv := range req.BoolFields {
		wg.Add(1)
		go func(field metaBoolField) {
			defer wg.Done()
			m.metaConfig.SetBool(&field.Key, &field.Value)
			hook, exists := m.updateHooks.BoolFieldHooks[field.Key]
			if exists {
				err := hook(field.Key, field.Value)
				if err!=nil {
					resMutex.Lock()
					res.Err = append(res.Err, err)
					resMutex.Unlock()
				}
			}
		}(kv)
	}

	// Double fields
	for _,kv := range req.DoubleFields {
		wg.Add(1)
		go func(field metaDoubleField) {
			defer wg.Done()
			m.metaConfig.SetDouble(&field.Key, &field.Value)
			hook, exists := m.updateHooks.DoubleFieldHooks[field.Key]
			if exists {
				err := hook(field.Key, field.Value)
				if err!=nil {
					resMutex.Lock()
					res.Err = append(res.Err, err)
					resMutex.Unlock()
				}
			}
		}(kv)
	}

	// List fields
	for _,kv := range req.ListFields {
		wg.Add(1)
		go func(field metaListField) {
			defer wg.Done()
			m.metaConfig.SetList(&field.Key, &field.Value)
			hook, exists := m.updateHooks.ListFieldHooks[field.Key]
			if exists {
				err := hook(field.Key, field.Value)
				if err!=nil {
					resMutex.Lock()
					res.Err = append(res.Err, err)
					resMutex.Unlock()
				}
			}
		}(kv)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
