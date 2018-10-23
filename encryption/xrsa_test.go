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

func TestEncryptViaPrivateKeyAndDencryptViaPublicKey(t *testing.T) {
	privKeyFile := "samplekeys"
	pubKeyFile := "samplekeys.pem"

	privk, _ := os.Create(privKeyFile)
	pubk, _ := os.Create(pubKeyFile)

	err := CreateKeys(privk, pubk)
	if err != nil {
		t.Error("Failed to create keys", err)
		return
	}

	data := "hello"

	encData, err := EncryptViaPrivateKey(privKeyFile, data)
	if err != nil {
		t.Error("Failed to encrypt data", err)
		return
	}

	t.Logf("Data : %s --> %s", data, encData)

	decData, err := DencryptViaPublicKey(pubKeyFile, encData)
	if err != nil {
		t.Error("Failed to dencrypt data", err)
		return
	}

	if decData != data {
		t.Error("Decrypted data didn't match the original data", err)
		return
	}

	os.Remove("samplekeys")
	os.Remove("samplekeys.pem")
}
