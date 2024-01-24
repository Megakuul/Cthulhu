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
