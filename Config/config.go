package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Name          string   `yaml:"name"`
	Port          string   `yaml:"port"`
	Version       string   `yaml:"version"`
	Sources       []string `yaml:"sources"`
	NegativeWords []string `yaml:"negativeWords"`
	DataDir       string   `yaml:"datadir"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) CheckConfig() {
	defaultVal := false
	if c.Port == "" {
		log.Println("Port address has not been assigned ")
	}
	if len(c.Sources) == 0 {
		log.Println("No sources found")
	}
	if len(c.NegativeWords) == 0 {
		log.Println("No negative keywords set")
	}
	if defaultVal {
		log.Fatalf("Failed to start %s due to misconfiguation", c.Name)
	}
}

func (c *Config) Load() *Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	viper.AddConfigPath(dir + "/Config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	if os.Getenv("env") == "development" {
		viper.SetConfigName("config.dev")
	} else {
		viper.SetConfigName("config")
	}
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Println(readConfigErr)
	}
	unmarshalErr := viper.Unmarshal(&c)
	if unmarshalErr != nil {
		log.Println(unmarshalErr)
	}
	c.CheckConfig()
	return c
}
