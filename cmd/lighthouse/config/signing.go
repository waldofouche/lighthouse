package config

import (
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/pkg/errors"
)

type signingConf struct {
	Alg                  string                 `yaml:"alg"`
	Algorithm            jwa.SignatureAlgorithm `yaml:"-"`
	RSAKeyLen            int                    `yaml:"rsa_keylen"`
	KeyFile              string                 `yaml:"key_file"`
	KeyDir               string                 `yaml:"key_dir"`
	AutomaticKeyRollover rolloverConf           `yaml:"automatic_key_rollover"`
}

type rolloverConf struct {
	Enabled  bool  `yaml:"enabled"`
	Interval int64 `yaml:"interval"`
}

var defaultSigningConf = signingConf{
	Alg:       "ES512",
	RSAKeyLen: 2048,
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
	var ok bool
	c.Algorithm, ok = jwa.LookupSignatureAlgorithm(c.Alg)
	if !ok {
		return errors.New("error in signing conf: unknown algorithm: " + c.Alg)
	}
	return nil
}
