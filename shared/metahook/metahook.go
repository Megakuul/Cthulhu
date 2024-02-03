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
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/megakuul/cthulhu/shared/dataloader"
)

type MetaHook struct {
	metaConfig *dataloader.MetaConfig
	socketPath string
	socketServer *http.Server
	socketServerMux *http.ServeMux
}

func CreateMetaHook(path string, config *dataloader.MetaConfig) (*MetaHook, error) {
	// Create path recursively
	parentpath := filepath.Dir(path)
	if err := os.MkdirAll(parentpath, 0755); err!=nil {
		return nil, err
	}
	
	// Create ServeMux and register handlers
	sockMux := http.NewServeMux()
	// TODO: define handlers here
	
	// Create HTTP Server
	sockSrv := &http.Server{
		Handler: sockMux,
	}
	
	return &MetaHook{
		config,
		path,
		sockSrv,
		sockMux,
	}, nil
}

func (m* MetaHook) Serve() error {
	unixListener, err := net.Listen("unix", m.socketPath)
	if err!=nil {
		return err
	}
	defer unixListener.Close()
	defer os.Remove(m.socketPath)

	if err:=m.socketServer.Serve(unixListener); err!=nil && err!=http.ErrServerClosed {
		return err
	}
	return nil
}

func (m* MetaHook) Exit() error {
	return m.socketServer.Shutdown(context.Background())
}
