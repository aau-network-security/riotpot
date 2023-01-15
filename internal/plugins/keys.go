package plugins

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/riotpot/internal/logger"
)

type (
	KeyType string
	KeySize int
)

const (
	Public  KeyType = "public"
	Private KeyType = "private"

	InsecureKey KeySize = 1024
	LiteKey     KeySize = 2048
	DefaultKey  KeySize = 4096
)

type CKey interface {
	Generate() []byte
	GetPEM() []byte
	SetPEM(pem []byte)
}

type AbstractKey struct {
	CKey
	pem []byte
}

func (k *AbstractKey) GetPEM() []byte {
	return k.pem
}

func (k *AbstractKey) SetPEM(pem []byte) {
	k.pem = pem
}

type PrivateKey struct {
	key  AbstractKey
	priv *rsa.PrivateKey
}

func (k *PrivateKey) GetPEM() []byte {
	return k.key.GetPEM()
}

func (k *PrivateKey) SetPEM(pem []byte) {
	k.key.SetPEM(pem)
}

func (k *PrivateKey) SetKey(key *rsa.PrivateKey) {
	k.priv = key
}

// Function to Generate and store a private RSA key and PEM
func (k *PrivateKey) Generate(size KeySize) (cert []byte) {
	reader := rand.Reader
	priv, err := rsa.GenerateKey(reader, int(size))
	if err != nil {
		logger.Log.Fatal().Err(err)
	}

	err = priv.Validate()
	if err != nil {
		logger.Log.Fatal().Err(err)
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	cert = pem.EncodeToMemory(block)

	k.SetKey(priv)
	k.SetPEM(cert)
	return
}

func NewPrivateKey(size KeySize) *PrivateKey {
	k := &PrivateKey{}

	k.Generate(size)
	return k
}
