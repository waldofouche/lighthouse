package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/lighthouse/cmd/lighthouse/config"
	"github.com/go-oidfed/lighthouse/internal/utils/fileutils"
)

func mustLoadKey() crypto.Signer {
	conf := config.Get().Signing
	if !fileutils.FileExists(conf.KeyFile) {
		sk, err := generateKey(conf.Algorithm, conf.RSAKeyLen)
		if err != nil {
			log.Fatal(err)
		}
		if err = os.WriteFile(config.Get().Signing.KeyFile, exportPrivateKeyAsPemStr(sk), 0600); err != nil {
			log.Fatal(err)
		}
		return sk
	}
	sk, err := loadKey(conf.KeyFile, conf.Algorithm)
	if err != nil {
		log.Fatal(err)
	}
	return sk
}

// loadKey loads the private and public key from the passed keyfile
func loadKey(keyfile string, alg jwa.SignatureAlgorithm) (crypto.Signer, error) {
	keyFileContent, err := fileutils.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyFileContent)
	var sk crypto.Signer
	switch alg {
	case jwa.RS256(), jwa.RS384(), jwa.RS512(), jwa.PS256(), jwa.PS384(), jwa.PS512():
		sk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case jwa.ES256(), jwa.ES384(), jwa.ES512():
		sk, err = x509.ParseECPrivateKey(block.Bytes)
	case jwa.EdDSA():
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var ok bool
		sk, ok = key.(ed25519.PrivateKey)
		if !ok {
			return nil, errors.New("not an Ed25519 Private Key")
		}
	default:
		return nil, errors.New("unknown signing algorithm: " + alg.String())
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sk, nil
}

var signingKey crypto.Signer
var signingJWKS jwk.Set

func initKey() {
	signingKey = mustLoadKey()

	key, err := jwk.PublicKeyOf(signingKey.Public())
	if err != nil {
		log.Fatal(err)
	}
	if err = jwk.AssignKeyID(key); err != nil {
		log.Fatal(err)
	}
	if err = key.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
		log.Fatal(err)
	}
	if err = key.Set(jwk.AlgorithmKey, config.Get().Signing.Algorithm); err != nil {
		log.Fatal(err)
	}
	signingJWKS = jwk.NewSet()
	if err = signingJWKS.AddKey(key); err != nil {
		log.Fatal(err)
	}
}

// generateKey generates a cryptographic private key with the passed properties
func generateKey(alg jwa.SignatureAlgorithm, rsaKeyLen int) (
	sk crypto.Signer, err error,
) {
	switch alg {
	case jwa.RS256(), jwa.RS384(), jwa.RS512(), jwa.PS256(), jwa.PS384(), jwa.PS512():
		if rsaKeyLen <= 0 {
			return nil, errors.Errorf("%s specified, but no valid RSA key len", alg)
		}
		sk, err = rsa.GenerateKey(rand.Reader, rsaKeyLen)
	case jwa.ES256():
		sk, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case jwa.ES384():
		sk, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case jwa.ES512():
		sk, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case jwa.EdDSA():
		_, sk, err = ed25519.GenerateKey(rand.Reader)
	default:
		err = errors.Errorf("unknown signing algorithm '%s'", alg)
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// exportPrivateKeyAsPemStr exports the private key
func exportPrivateKeyAsPemStr(sk crypto.Signer) []byte {
	switch sk := sk.(type) {
	case *rsa.PrivateKey:
		return exportRSAPrivateKeyAsPem(sk)
	case *ecdsa.PrivateKey:
		return exportECPrivateKeyAsPem(sk)
	case ed25519.PrivateKey:
		return exportEDDSAPrivateKeyAsPem(sk)
	default:
		return nil
	}
}

func exportECPrivateKeyAsPem(privkey *ecdsa.PrivateKey) []byte {
	privkeyBytes, _ := x509.MarshalECPrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}

func exportRSAPrivateKeyAsPem(privkey *rsa.PrivateKey) []byte {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}

func exportEDDSAPrivateKeyAsPem(privkey ed25519.PrivateKey) []byte {
	privkeyBytes, _ := x509.MarshalPKCS8PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}
