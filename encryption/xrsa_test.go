package encryption

import (
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateKey(t *testing.T) {
	privKeyFile := "samplekeys"
	pubKeyFile := "samplekeys.pem"

	privk, _ := os.Create(privKeyFile)
	pubk, _ := os.Create(pubKeyFile)

	err := CreateKeys(privk, pubk)
	if err != nil {
		t.Errorf("Failed creating samplekeys %v\n", err)
	}

	privBytes, _ := ioutil.ReadFile(privKeyFile)
	privBlock, _ := pem.Decode(privBytes)
	if privBlock.Type != "PRIVATE KEY" {
		t.Errorf("`%s` should be a private key", privKeyFile)
	}

	pubBytes, _ := ioutil.ReadFile(pubKeyFile)
	pubBlock, _ := pem.Decode(pubBytes)
	if pubBlock.Type != "PUBLIC KEY" {
		t.Errorf("`%s` should be a public key", pubKeyFile)
	}

	os.Remove("samplekeys")
	os.Remove("samplekeys.pem")
}

func TestLoadKeys(t *testing.T) {
	privKeyFile := "samplekeys"
	pubKeyFile := "samplekeys.pem"

	privk, _ := os.Create(privKeyFile)
	pubk, _ := os.Create(pubKeyFile)

	CreateKeys(privk, pubk)

	xrsa, err := LoadKeys(privKeyFile, pubKeyFile)
	if err != nil {
		t.Error("Failed to load keys", err)
	}

	if xrsa.privateKey.PublicKey.E != xrsa.publicKey.E {
		t.Error("Loaded keys not aligned")
	}

	os.Remove("samplekeys")
	os.Remove("samplekeys.pem")
}