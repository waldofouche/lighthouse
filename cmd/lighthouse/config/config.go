package config

import (
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/go-oidfed/lighthouse"
)

// Config holds configuration for the entity
type Config struct {
	Server     lighthouse.ServerConf `yaml:"server"`
	Logging    loggingConf           `yaml:"logging"`
	Storage    storageConf           `yaml:"storage"`
	Signing    signingConf           `yaml:"signing"`
	Endpoints  Endpoints             `yaml:"endpoints"`
	Federation federationConf        `yaml:"federation_data"`
}

type configValidator interface {
	validate() error
}

// Validate checks all fields of Config that implement configValidator
func (c *Config) Validate() error {
	v := reflect.ValueOf(c).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldIface := fieldVal.Interface()

		validator, ok := fieldIface.(configValidator)
		if ok {
			if err := validator.validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

var c = Config{
	Server:     defaultServerConf,
	Logging:    defaultLoggingConf,
	Storage:    storageConf{},
	Signing:    defaultSigningConf,
	Endpoints:  Endpoints{},
	Federation: defaultFederationConf,
}

// Get returns the Config
func Get() Config {
	return c
}

// Load loads the config from the given file
func Load(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(content, &c); err != nil {
		log.Fatal(err)
	}
	if err = c.Validate(); err != nil {
		log.Fatal(err)
	}

}
