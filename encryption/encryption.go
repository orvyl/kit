package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	defaultKeyLength   = 2048
	defaultPrivKeyName = "enc_key"
	defaultPubKeyName  = "enc_key.pem"
)

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultKeyLength)
	if err != nil {
		return fmt.Errorf("Failed to generate key: %s", err)
	}

	if err = createKey(privateKeyWriter, "PRIVATE KEY", privateKey); err != nil {
		return err
	}
	if err = createKey(publicKeyWriter, "PUBLIC KEY", &privateKey.PublicKey); err != nil {
		return err
	}

	return nil
}

func createKey(w io.Writer, keyType string, v interface{}) error {
	var content []byte
	var err error
	switch keyType {
	case "PRIVATE KEY":
		content, err = x509.MarshalPKCS8PrivateKey(v)
	case "PUBLIC KEY":
		content, err = x509.MarshalPKIXPublicKey(v)
	default:
		return fmt.Errorf("Key type %s not supported", keyType)
	}

	if err != nil {
		return fmt.Errorf("Failed to marshall generated key: %s", err)
	}

	block := &pem.Block{
		Type:  keyType,
		Bytes: content,
	}

	err = pem.Encode(w, block)
	if err != nil {
		return fmt.Errorf("Failed to encode generated key: %s", err)
	}

	return nil

}

func LoadKeys(privKeyFile, pubKeyFile string) (*XRsa, error) {
	privKeyIntr, err := loadKey(privKeyFile, "PRIVATE KEY")
	if err != nil {
		return nil, err
	}
	privateKey, ok := privKeyIntr.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Failed to convert parsed private key %s: %s", privKeyFile, err)
	}

	pubKeyIntr, err := loadKey(pubKeyFile, "PUBLIC KEY")
	if err != nil {
		return nil, err
	}
	publicKey, ok := pubKeyIntr.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to convert parsed public key %s: %s", pubKeyFile, err)
	}

	return &XRsa{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func loadKey(f, keyType string) (interface{}, error) {
	file, err := ioutil.ReadFile(f)
	keyBlock, _ := pem.Decode(file)

	var keyIntr interface{}
	switch keyType {
	case "PRIVATE KEY":
		keyIntr, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	case "PUBLIC KEY":
		keyIntr, err = x509.ParsePKIXPublicKey(keyBlock.Bytes)
	default:
		return nil, fmt.Errorf("Key type %s not supported", keyType)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse key %s: %s", f, err)
	}

	return keyIntr, nil
}

/*
REFs:
https://medium.com/@raul_11817/export-import-pem-files-in-go-67614624adc7
https://medium.com/@raul_11817/golang-cryptography-rsa-asymmetric-algorithm-e91363a2f7b3
https://github.com/liamylian/x-rsa/blob/master/golang/xrsa/xrsa.go
https://stackoverflow.com/questions/18011708/encrypt-message-with-rsa-private-key-as-in-openssls-rsa-private-encrypt
*/
