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
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const TMP_FILE_EXTENSION string = ".tmp"

/**
 * Object holding a inmem configuration
 *
 * Configuration can be written and read from/to disk
 *
 * All operations that are fully thread-safe (synchronized).
 *
 * Uses a custom parser, that parses a kind of a key-value config file (example):
 *
 * ```
 *
 * # I'm a comment until newline
 * somekey="some.value;9?
 * I can contain spaces, tabs, newlines
 * "uglyplacedkey="I'm valid too"
 *
 * wellplacedkey=""
 * / I'm also a comment until newline
 * ```
 */
type MetaConfig struct {
	// Mutex lock for the configuration file
	configFileLock sync.RWMutex
	// Mutex lock for the inmem configuratio
	configLock sync.RWMutex
	// Path of the configuration
	configPath string
	// In memory configuration object
	config map[string]string
}

/**
 * Initializes MetaConfig and creates the config file if not existent
 */
func CreateMetaConfig(path string) (*MetaConfig, error) {
	config := &MetaConfig{}
	config.configPath = path
	// Generate file path recursively
	parentpath := filepath.Dir(config.configPath)
	if err := os.MkdirAll(parentpath, 0755); err!=nil {
		return config, err
	}
	// Generate file
	_, err := os.Create(config.configPath);
	return config, err
}

/**
 * Returns true if the key exists and false if it doesn't
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) Exists(key *string) bool {
	m.configLock.RLock()
	defer m.configLock.RUnlock()
	_, exists := m.config[*key]
	return exists
}

/**
 * Get full parsed configuration object
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) GetConfig(key *string) map[string]string {
	m.configLock.RLock()
	defer m.configLock.RUnlock()

	mapBuf := make(map[string]string)

	for k,v := range m.config {
		mapBuf[k] = v
	}
	return mapBuf
}

/**
 * Get string value of specific key
 *
 * If key is not found, it will return an empty string
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) GetString(key *string) string {
	m.configLock.RLock()
	defer m.configLock.RUnlock()	
	val, _ := m.config[*key]
	return val
}

/**
 * Get bool value of specific key
 *
 * Underlying string is evaluated true if it is set to "true" or "YES"
 *
 * If key is not found, it will return false
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) GetBool(key *string) bool {
	m.configLock.RLock()
	defer m.configLock.RUnlock()
	
	val, exists := m.config[*key]
	if exists {
		return strings.ToLower(val)=="true"||strings.ToLower(val)=="yes"
	} else {
		return false
	}
}

/**
 * Get double value of specific key
 *
 * If the conversion fails (invalid double in config) it will return 0
 *
 * If key is not found, it will return 0 aswell
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) GetDouble(key *string) float64 {
	m.configLock.RLock()
	defer m.configLock.RUnlock()
	
	val, exists := m.config[*key]
	if exists {
		numval, err := strconv.ParseFloat(val, 64)
		if err!=nil {
			return 0.0
		}
		return numval
	} else {
		return 0.0
	}
}

/**
 * Get list value of specific key
 *
 * Underlying string is splitted based on ','
 * empty fields ("") are omitted
 *
 * If key is not found, it will return a empty list
 *
 * This operation does not read / parse anything from disk!
 */
func (m* MetaConfig) GetList(key *string) []string {
	m.configLock.RLock()
	defer m.configLock.RUnlock()
	
	val, exists := m.config[*key]
	if exists {
		// Split tokens
		listval := strings.Split(val, ",")
		// Remove empty fields
		var tokens []string
		for _, tokBuf := range listval {
			if tokBuf!="" {
				tokens = append(tokens, tokBuf)
			}
		}
		return tokens
	} else {
		return []string{}
	}
}

/**
 * Set string value to specific key
 *
 * This operation does not write anything to disk!
 */
func (m* MetaConfig) SetString(key *string, value *string) {
	m.configLock.Lock()
	defer m.configLock.Unlock()

	m.config[*key] = *value
}

/**
 * Set bool value to specific key
 *
 * This operation does not write anything to disk!
 */
func (m* MetaConfig) SetBool(key *string, value *bool) {
	m.configLock.Lock()
	defer m.configLock.Unlock()

	if *value {
		m.config[*key] = "true"
	} else {
		m.config[*key] = "false"
	}
}

/**
 * Set double value to specific key
 *
 * This operation does not write anything to disk!
 */
func (m* MetaConfig) SetDouble(key *string, value *float64) {
	m.configLock.Lock()
	defer m.configLock.Unlock()

	m.config[*key] = strconv.FormatFloat(*value, 'f', -1, 64)
}


/**
 * Set list value to specific key
 *
 * This operation does not write anything to disk!
 */
func (m* MetaConfig) SetList(key *string, value *[]string) {
	m.configLock.Lock()
	defer m.configLock.Unlock()

	outstr := ""
	for _,val := range *value {
		outstr+=val
		outstr+=","
	}
	m.config[*key] = outstr
}

/**
 * Read and Parse configuration directly from disk to inmem config
 *
 * If a key is placed multiple times, only the first one is evaluated
 *
 * Function will throw a runtime error if it fails
 */
func (m* MetaConfig) ReadFromDisk() error {
	// Read lock the file config lock
	m.configFileLock.RLock()
	defer m.configFileLock.RUnlock()

	mapBuffer := make(map[string]string)
	// Read config file
	file, err := os.OpenFile(m.configPath, os.O_RDONLY, 0755)
	if err!=nil {
		return err
	}
	defer file.Close()

	// Create buffered reader
	reader := bufio.NewReader(file)
	
	// Char buffer
	var c byte
	// Keeps track of lines for debug messages
	var lineCount int = 0
	// Current key buffer
	var curKey strings.Builder 
	// Current value buffer
	var curVal strings.Builder

	// Unnamed helper function to read one char at a time
	getChar := func(char *byte) bool {
		var val byte
		// Read next byte from stream
		val, err = reader.ReadByte()
		if err!=nil {
			return false
		}
		*char = val
		return true
	}

	// Iterate over chars
	for {
		// Eat next char
		if !getChar(&c) {
			break
		}
		// Skip newline
		if c=='\n' {
			lineCount++
			continue
		}
		// Skip space, tab
		if c==' '||c=='\t'||c=='\r' {
			continue
		}
		// # | / indicate a comment
		if c=='#'||c=='/' {
			// Skip til EOF or newline
			for {
				if !(getChar(&c)&&c!='\n') {
					break
				}
			}
			lineCount++
			continue
		}

		// Eat key
		curKey.Reset()
		for {
			curKey.WriteByte(c)
			// EOF or newline in key is not allowed
			if !getChar(&c)||c=='\n' {
				return fmt.Errorf(
					"Failed to parse config file at: %s\nUnexpected EOF or newline on line: %d",
					m.configPath, lineCount,
				)
			}
			// Read until '=' char
			if c=='=' {
				break
			}
		}

		// Read next char which is expected to be '"'
		if !getChar(&c)||c!='"' {
			return fmt.Errorf(
				"Failed to parse config file at: %s\nExpected '\"' after '=' on line: %d",
				m.configPath, lineCount,
			)
		}

		// Eat value
		curVal.Reset()
		for {
			// EOF is not expected in value, every other char can be used
			if !getChar(&c) {
				return fmt.Errorf(
					"Failed to parse config file at: %s\nUnexpected EOF on line: %d",
					m.configPath, lineCount,
				)
			} else if c=='"' {
				// Read until '"' char
				break
			}
			// Add linecount
			if c=='\n' {
				lineCount++
			}
			curVal.WriteByte(c)
		}
		// Insert first pair, the later pairs with same key are ignored
		strKey, strVal := curKey.String(), curVal.String()
		if _, exists := mapBuffer[strKey]; !exists {
			mapBuffer[strKey] = strVal
		}
	}

	// Error is expected to be EOF, if not there was a reading failure
	if err!=io.EOF {
		return err
	}

	// Write lock the inmen config lock
	m.configLock.Lock()
	defer m.configLock.Unlock()
	m.config = mapBuffer
	return nil
}

/**
 * Writes inmem configuration directly to disk
 *
 * Function will throw a runtime error if it fails
 */
func (m* MetaConfig) WriteToDisk() error {
	// Write lock the file config lock
	m.configFileLock.Lock()
	defer m.configFileLock.Unlock()
	// Read lock the inmem config lock
	m.configLock.RLock()
	defer m.configLock.RUnlock()

	// Outstr buffer
	var outstr string
	// Open tmp config file
	file, err := os.OpenFile(m.configPath+TMP_FILE_EXTENSION, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err!=nil {
		return err
	}

	// Insert deparsed configuration
	outstr += "# Manual changes to configuration may be overwritten\n"
	outstr += "# Consider using Meta Hook from the Cthulhu component\n"
	for k,v := range m.config {
		outstr += k
		outstr += "="
		outstr += "\""
		outstr += v
		outstr += "\""
		outstr += "\n"
	}
	outstr += "# End of config\n"
	
	_, err = file.Write([]byte(outstr))
	if err!=nil {
		file.Close()
		return err
	}
	
	err = file.Close()
	if err!=nil {
		return err
	}

	// Move tmp config to config
	// This prevents file corruption on unexpected application crashes (e.g. shutdown while writing).
	return os.Rename(m.configPath + TMP_FILE_EXTENSION, m.configPath)
}
