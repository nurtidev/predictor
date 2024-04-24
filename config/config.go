package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Mode      string   `yaml:"mode"`
	Filepath  string   `yaml:"filepath"`
	Market    []string `yaml:"market"`
	Timeframe []string `yaml:"timeframe"`
	Telegram  struct {
		Enable bool   `yaml:"enable"`
		Token  string `yaml:"token"`
		Users  []int  `yaml:"users"`
	} `yaml:"telegram"`
	Trade struct {
		BufferSize []int `yaml:"buffer_size"`
		Breakdown  struct {
			Percent float64 `yaml:"percent"`
			MinSize int     `yaml:"min_size"`
			MaxSize int     `yaml:"max_size"`
		} `yaml:"breakdown"`
		Motion struct {
			MinSize int `yaml:"min_size"`
			MaxSize int `yaml:"max_size"`
		} `yaml:"motion"`
	} `yaml:"trade"`
	Filters struct {
		Fractal struct {
			Size int `yaml:"size"`
		} `yaml:"fractal"`
	} `yaml:"filters"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
