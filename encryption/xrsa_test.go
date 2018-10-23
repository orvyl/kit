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

func TestPrivateEncryptPublicDecrypt(t *testing.T) {
	privKeyFile := "samplekeys"
	pubKeyFile := "samplekeys.pem"

	privk, _ := os.Create(privKeyFile)
	pubk, _ := os.Create(pubKeyFile)

	err := CreateKeys(privk, pubk)
	if err != nil {
		t.Error("Failed to create keys", err)
		return
	}

	xrsa, err := LoadKeys(privKeyFile, pubKeyFile)
	if err != nil {
		t.Error("Failed to load keys", err)
		return
	}

	data := "hello"

	encData, err := xrsa.PrivateEncrypt(data)
	if err != nil {
		t.Error("Failed to encrypt using private key", err)
		return
	}

	t.Logf("Encryption result : %s -> %s\n", data, encData)

	decData, err := xrsa.PublicDecrypt(encData)
	if err != nil {
		t.Error("Failed to decrypt using pub key", err)
		return
	}

	if decData != data {
		t.Error("Data and decnrypted data didn't match")
		return
	}

	os.Remove("samplekeys")
	os.Remove("samplekeys.pem")
}
