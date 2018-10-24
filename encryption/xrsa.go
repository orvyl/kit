package encryption

/*
REFs:
https://medium.com/@raul_11817/export-import-pem-files-in-go-67614624adc7
https://medium.com/@raul_11817/golang-cryptography-rsa-asymmetric-algorithm-e91363a2f7b3
https://github.com/liamylian/x-rsa/blob/master/golang/xrsa/xrsa.go
https://stackoverflow.com/questions/18011708/encrypt-message-with-rsa-private-key-as-in-openssls-rsa-private-encrypt
*/

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/liamylian/x-rsa/golang/xrsa"
)

const (
	defaultKeyLength = 1024
)

//CreateKeys create an RSA type keys in the provided `io.Writer`s
func CreateKeys(privateKeyWriter, publicKeyWriter io.Writer) error {
	return xrsa.CreateKeys(publicKeyWriter, privateKeyWriter, defaultKeyLength)
}

//EncryptViaPrivateKey will encrypt given `data` using a provided private key
func EncryptViaPrivateKey(privateKeyF, data string) (string, error) {
	block, err := getKeyBlock(privateKeyF)
	if err != nil {
		return "", err
	}
	privIntr, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pri, ok := privIntr.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("Failed to cast value to `rsa.PrivateKey`")
	}

	partLen := pri.PublicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bts, err := xrsa.PrivateEncrypt(pri, chunk)
		if err != nil {
			return "", err
		}

		buffer.Write(bts)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

//DecryptViaPublicKey will decrypt given `encData` using a provided public key
func DecryptViaPublicKey(publicKeyF, encData string) (string, error) {
	block, err := getKeyBlock(publicKeyF)
	if err != nil {
		return "", err
	}
	pubIntr, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pub, ok := pubIntr.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("Failed to cast value to `rsa.PublicKey`")
	}

	partLen := pub.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encData)
	chunks := split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := xrsa.PublicDecrypt(pub, chunk)

		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), nil
}

func getKeyBlock(keyFile string) (*pem.Block, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("Failed to decode key file %s", keyFile)
	}

	return block, nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
