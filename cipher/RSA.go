package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//RSADecrypt PKCS8解密
func RSADecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//X509解码
	prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pKey := prkI.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, pKey, ciphertext) //RSA算法解密
}

//RSAEncrypt PKCS8加密
func RSAEncrypt(origData, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey) //将密钥解析成公钥实例
	if block == nil {
		return nil, errors.New("public key error")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	//类型断言
	pKey := publicKeyInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pKey, origData) //RSA算法加密
}
