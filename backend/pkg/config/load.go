package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(c *Config, cfgFile string) error {
	if cfgFile == "" {
		return nil
	}

	k := koanf.New("!")
	if _, err := os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("[INFO] No config file found")
		return nil
	}

	if err := k.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		fmt.Printf("[ERROR] failed to load config file: %v\n", err)
	}

	err := k.Unmarshal("", c)
	if err != nil {
		return err
	}

	return err
}
