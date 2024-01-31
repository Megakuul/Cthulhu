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

package dataloader

import (
	"sync"
	"github.com/spf13/viper"
)

type MetaConfig struct {
	metaLock sync.RWMutex
	storagePath string `mapstructure:"storagepath"`
	clusterNodeAddr []string `mapstructure:"clusternodeaddr"`
}

func (m *MetaConfig) GetStoragePath() string {
	m.metaLock.Lock()
	defer m.metaLock.Unlock()
	return m.storagePath
}

func (m *MetaConfig) GetClusterNodeAddr() []string {
	m.metaLock.Lock()
	defer m.metaLock.Unlock()
	return m.clusterNodeAddr
}

func (m *MetaConfig) FetchData() error {
	m.metaLock.Lock()
	defer m.metaLock.Unlock();

	if err := viper.ReadInConfig(); err!=nil {
		return err
	}

	if err := viper.Unmarshal(m); err!=nil {
		return err
	}
	return nil
}

func CreateMetaConfig(path string) *MetaConfig {
	viper.SetConfigType("properties")
	viper.SetConfigFile(path)
	viper.SetTypeByDefaultValue(true)
	return &MetaConfig{}
}
