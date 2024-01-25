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
	MetaLock sync.RWMutex
	StoragePath string `mapstructure:"storagepath"`
	ClusterNodeAddr []string `mapstructure:"clusternodeaddr"`
}

func (m *MetaConfig) FetchData() error {
	m.MetaLock.Lock()
	defer m.MetaLock.Unlock();

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
