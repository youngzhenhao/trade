package config

import (
	"gopkg.in/yaml.v3"
	"path/filepath"
	"trade/utils"
)

func generateConfig() string {
	conf, _ := yaml.Marshal(&Config{})
	return string(conf)
}

func WriteConfigExample(dir string) {
	utils.CreateFile(filepath.Join(dir, "config.yaml.example"), generateConfig())
}
