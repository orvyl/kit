package encryption

/*
REFs:
https://medium.com/@raul_11817/export-import-pem-files-in-go-67614624adc7
https://medium.com/@raul_11817/golang-cryptography-rsa-asymmetric-algorithm-e91363a2f7b3
https://github.com/liamylian/x-rsa/blob/master/golang/xrsa/xrsa.go
https://stackoverflow.com/questions/18011708/encrypt-message-with-rsa-private-key-as-in-openssls-rsa-private-encrypt
*/

import (
	"io"
	"io/ioutil"

	"github.com/liamylian/x-rsa/golang/xrsa"
)

const (
	defaultKeyLength = 2048
)

//CreateKeys create an RSA type keys in the provided `io.Writer`s
func CreateKeys(privateKeyWriter, publicKeyWriter io.Writer) error {
	return xrsa.CreateKeys(publicKeyWriter, privateKeyWriter, defaultKeyLength)
}

//LoadKeys will load provided key files and return a `XRsa` value for encryption-decryption
func LoadKeys(privKeyFile, pubKeyFile string) (*xrsa.XRsa, error) {
	privKeyFByte, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, err
	}

	pubKeyFByte, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, err
	}
	return xrsa.NewXRsa(pubKeyFByte, privKeyFByte)
}
