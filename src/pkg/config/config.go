package config

import (
	_ "embed"
	"os"

	"github.com/ilrudie/bulk-transcode/src/pkg/ffmpeg"
	"sigs.k8s.io/yaml"
)

//go:embed default.yaml
var defaultConfig []byte

type Config struct {
	InputDir         string      `json:"input_dir"`
	OutputDir        string      `json:"output_dir"`
	OutputMark       string      `json:"output_mark"`
	Recursive        bool        `json:"recursive"`
	Exec             bool        `json:"exec"`
	CommandArguments ffmpeg.Args `json:"command_arguments"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg *Config
	yaml.Unmarshal(data, cfg)
	return cfg, nil
}

func DefaultConfig() *Config {
	c := &Config{}
	err := yaml.Unmarshal(defaultConfig, c)
	if err != nil {
		panic("Failed to load default configuration")
	}
	return c
}

func (c *Config) ArgOverrides(inputDir, outputDir, mark string, exec, eset, recursive, rset bool) {
	if inputDir != "" {
		c.InputDir = inputDir
	}
	if outputDir != "" {
		c.OutputDir = outputDir
	}
	if mark != "" {
		c.OutputMark = mark
	}
	if rset {
		c.Recursive = recursive
	}
	if eset {
		c.Exec = exec
	}
}
