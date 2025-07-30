package main

import (
	"github.com/go-oidfed/lib/jwx"

	"github.com/go-oidfed/lighthouse/cmd/lighthouse/config"
)

var keys *jwx.KeyStorage

func initKey() (err error) {
	c := config.Get().Signing

	keys, err = jwx.NewKeyStorage(
		c.KeyDir, map[string]jwx.KeyStorageConfig{
			jwx.KeyStorageTypeFederation: {
				Algorithm:    c.Alg,
				RSAKeyLen:    c.RSAKeyLen,
				RolloverConf: c.AutomaticKeyRollover,
			},
		},
	)
	if err != nil {
		return
	}
	return
}
