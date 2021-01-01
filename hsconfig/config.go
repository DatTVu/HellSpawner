package hsconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

type Config struct {
	RecentProjects   []string
	AuxiliaryMpqPath string
}

func getConfigPath() string {
	configPath := configdir.LocalConfig("hellspawner")
	if err := configdir.MakePath(configPath); err != nil {
		log.Fatal(err)
	}
	return filepath.Join(configPath, "environment.json")
}

func generateDefaultConfig() *Config {
	result := &Config{
		RecentProjects: []string{},
	}
	result.Save()

	return result
}

func Load() *Config {
	configFile := getConfigPath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return generateDefaultConfig()
	}

	var err error
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err != nil {
		return generateDefaultConfig()
	}

	var result *Config
	if err = json.Unmarshal(data, &result); err != nil {
		return generateDefaultConfig()
	}

	return result
}

func (c *Config) Save() error {
	var err error
	var data []byte

	if data, err = json.MarshalIndent(c, "", "   "); err != nil {
		return err
	}

	if err = ioutil.WriteFile(getConfigPath(), data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}

func (c *Config) AddToRecentProjects(filePath string) {
	for idx := range c.RecentProjects {
		if c.RecentProjects[idx] == filePath {
			if idx != 0 {
				old := c.RecentProjects[0]
				c.RecentProjects[0] = filePath
				c.RecentProjects[idx] = old
			}
			return
		}
	}

	recent := []string{filePath}
	for idx := range c.RecentProjects {
		if idx == 5 {
			break
		}

		recent = append(recent, c.RecentProjects[idx])
	}

	c.RecentProjects = recent
	c.Save()
}
