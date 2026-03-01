package jwtkeys

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"

	"be2/services/authsvc/internal/config"
)

type RSAKey struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
	KID     string
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewRSAKey(secrets *config.Secrets) (*RSAKey, error) {
	priv, err := parsePrivateKey([]byte(secrets.JWTPrivateKey))
	if err != nil {
		return nil, err
	}
	kid, err := keyID(&priv.PublicKey)
	if err != nil {
		return nil, err
	}
	return &RSAKey{Private: priv, Public: &priv.PublicKey, KID: kid}, nil
}

func (k *RSAKey) JWK() (JWK, error) {
	if k.Public == nil {
		return JWK{}, errors.New("missing public key")
	}
	n := base64.RawURLEncoding.EncodeToString(k.Public.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.Public.E)).Bytes())
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: k.KID,
		N:   n,
		E:   e,
	}, nil
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM data")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("unsupported private key type")
		}
		return priv, nil
	default:
		return nil, errors.New("unsupported PEM block type")
	}
}

func keyID(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(der)
	return base64.RawURLEncoding.EncodeToString(sum[:]), nil
}
