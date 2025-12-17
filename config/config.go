package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GroupName string `env:"GROUP_NAME, overwrite, default=acme.mycompany.com" yaml:"groupName" json:"groupName"`
	AliDNS    struct {
		Region string `env:"REGION, overwrite, default=cn-hangzhou" yaml:"region" json:"region"`
	} `env:", prefix=ALIDNS_" yaml:"alidns" json:"alidns"`
	Log struct {
		Color  bool   `env:"COLOR, overwrite, default=true" yaml:"color" json:"color"`
		Format string `env:"FORMAT, overwrite, default=console" yaml:"format" json:"format"`
		Level  string `env:"LEVEL, overwrite, default=debug" yaml:"level" json:"level"`
	} `env:", prefix=LOG_" yaml:"log" json:"log"`
}

var GlobalConfig Config

func LoadGlobalConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if yamlFile != nil {
		if err := yaml.Unmarshal(yamlFile, &GlobalConfig); err != nil {
			return err
		}
	}

	if err := envconfig.Process(context.Background(), &GlobalConfig); err != nil {
		return err
	}

	l, _ := zerolog.ParseLevel(GlobalConfig.Log.Level)
	if l == zerolog.NoLevel {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(l)
	}
	if GlobalConfig.Log.Format != "json" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: !GlobalConfig.Log.Color})
	}

	b, err := json.Marshal(&GlobalConfig)
	if err != nil {
		return err
	}
	log.Debug().Msgf("loaded config: %s", string(b))

	return nil
}
