package config

import (
	"log"
	"os"
	"reflect"

	"github.com/pkg/errors"
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

// Validate checks all fields of Config that implement configValidator (pointer receivers)
func (c *Config) Validate() error {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)

		// Get addressable pointer to field if possible
		if fieldVal.CanAddr() {
			ptr := fieldVal.Addr().Interface()

			if validator, ok := ptr.(configValidator); ok {
				if err := validator.validate(); err != nil {
					return errors.Errorf("validation failed for field '%s': %w", t.Field(i).Name, err)
				}
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
