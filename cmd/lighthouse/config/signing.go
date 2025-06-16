package config

import (
	"github.com/pkg/errors"
)

type signingConf struct {
	Algorithm            string       `yaml:"alg"`
	KeyFile              string       `yaml:"key_file"`
	KeyDir               string       `yaml:"key_dir"`
	AutomaticKeyRollover rolloverConf `yaml:"automatic_key_rollover"`
}

type rolloverConf struct {
	Enabled  bool  `yaml:"enabled"`
	Interval int64 `yaml:"interval"`
}

var defaultSigningConf = signingConf{
	Algorithm: "ES512",
	AutomaticKeyRollover: rolloverConf{
		Enabled:  false,
		Interval: 600000,
	},
}

func (c *signingConf) validate() error {
	if c.KeyFile == "" && c.KeyDir == "" {
		return errors.New("error in signing conf: either key_file or key_dir must be specified")
	}
	if c.AutomaticKeyRollover.Enabled && c.KeyDir == "" {
		return errors.New(
			"error in signing conf: if automatic_key_rollover" +
				" is enabled, key_dir must be specified, not key_file",
		)
	}
	return nil
}
