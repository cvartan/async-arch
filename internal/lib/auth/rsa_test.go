package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestRSARandom(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	b := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	block := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: b,
	}
	out := pem.EncodeToMemory(&block)

	nblock, _ := pem.Decode(out)

	if nblock == nil {
		t.Fatal("FUCK")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(nblock.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(publicKey)
}
