package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

//後で解読する

var (
	certsPath      = "certs"
	certName       = "serverCert.pem"
	keyName        = "serverCert.key"
	clientCertName = "client.pem"
	clientKeyName  = "client.key"
)

//初期化でディレクトリがない場合は作成
func init() {
	//create certificates directory if it doesn't exist
	if _, err := os.Stat(certsPath); os.IsNotExist(err) {
		if err2 := os.MkdirAll(certsPath, 0755); err2 != nil {
			log.Fatalf("Could not create the path %s", certsPath)
		}
	}
}

//root証明書を作成
func genRootCert() (x509.Certificate, *ecdsa.PrivateKey, error) {
	notBefore := time.Now()
	notAfter := notBefore.AddDate(10, 0, 0)
	//左シフト
	serialNo, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	//鍵を作成
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//鍵を保存
	_, err = saveKey("rootCAKey.key", key)
	if err != nil {
		return x509.Certificate{}, key, err
	}
	//証明書を作成
	cert := x509.Certificate{
		SerialNumber: serialNo,
		//有効期限
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Subject: pkix.Name{
			Organization: []string{"Reverse Shell Root"},
			CommonName:   "Reverse Shell Root CA",
		},
	}
	der, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &key.PublicKey, key)
	if err != nil {
		return x509.Certificate{}, key, err
	}

	//pemを保存
	if _, err = saveCert("rootCACert.pem", der); err != nil {
		return x509.Certificate{}, key, err
	}
	return cert, key, nil
}

//GenCerts generates self-signed certificates
func GenServerCerts() (certFile, keyFile string, err error) {
	//鍵が存在する場合は、そのまま返す
	suppliedCert := filepath.Join(certsPath, certName)
	suppliedKey := filepath.Join(certsPath, keyName)
	//if certificate and key already exist attempt to use them, otherwise generate self-signed ones
	if _, err := os.Stat(suppliedCert); !os.IsNotExist(err) {
		if _, err := os.Stat(suppliedKey); !os.IsNotExist(err) {
			return suppliedCert, suppliedKey, nil
		}
	}
	//ルート証明書、鍵を作成
	rootCert, rootKey, err := genRootCert()
	if err != nil {
		return certFile, keyFile, err
	}
	//サーバー鍵、証明書を作成
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return certFile, keyFile, err
	}
	keyFile, err = saveKey(keyName, key)
	if err != nil {
		return certFile, keyFile, err
	}
	serialNo, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return certFile, keyFile, err
	}
	notBefore := time.Now()
	notAfter := notBefore.AddDate(10, 0, 0) //10 years
	cert := x509.Certificate{
		SerialNumber:          serialNo,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  false,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Subject: pkix.Name{
			Organization: []string{"Reverse Shell organisation"},
			CommonName:   "Reverse Shell Organisation certificate",
		},
		DNSNames: []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, &cert, &rootCert, &key.PublicKey, rootKey)
	if err != nil {
		return certFile, keyFile, err
	}
	certFile, err = saveCert(certName, der)
	if err != nil {
		return certFile, keyFile, err
	}
	return certFile, keyFile, nil
}

//GenCerts generates self-signed certificates
func GenClientCerts() (certFile, keyFile string, err error) {
	//鍵が存在する場合は、そのまま返す
	suppliedCert := filepath.Join(certsPath, clientCertName)
	suppliedKey := filepath.Join(certsPath, clientKeyName)
	//if certificate and key already exist attempt to use them, otherwise generate self-signed ones
	if _, err := os.Stat(suppliedCert); !os.IsNotExist(err) {
		if _, err := os.Stat(suppliedKey); !os.IsNotExist(err) {
			return suppliedCert, suppliedKey, nil
		}
	}
	//ルート証明書、鍵を作成
	rootCert, rootKey, err := genRootCert()
	if err != nil {
		return certFile, keyFile, err
	}
	//サーバー鍵、証明書を作成
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return certFile, keyFile, err
	}
	keyFile, err = saveKey(clientKeyName, key)
	if err != nil {
		return certFile, keyFile, err
	}
	serialNo, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return certFile, keyFile, err
	}
	notBefore := time.Now()
	notAfter := notBefore.AddDate(10, 0, 0) //10 years
	cert := x509.Certificate{
		SerialNumber:          serialNo,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  false,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Subject: pkix.Name{
			Organization: []string{"Reverse Shell organisation"},
			CommonName:   "Reverse Shell Organisation certificate",
		},
		DNSNames: []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, &cert, &rootCert, &key.PublicKey, rootKey)
	if err != nil {
		return certFile, keyFile, err
	}
	certFile, err = saveCert(clientCertName, der)
	if err != nil {
		return certFile, keyFile, err
	}
	return certFile, keyFile, nil
}

func saveKey(fileName string, key *ecdsa.PrivateKey) (string, error) {
	fileName = filepath.Join(certsPath, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return fileName, err
	}
	defer file.Close()
	kb, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return fileName, err
	}
	if err = pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: kb}); err != nil {
		return fileName, err
	}
	return fileName, nil
}

func saveCert(fileName string, der []byte) (string, error) {
	fileName = filepath.Join(certsPath, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return fileName, err
	}
	defer file.Close()
	if err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return fileName, err
	}
	return fileName, nil
}
