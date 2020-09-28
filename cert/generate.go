package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func GeneratePki(path, domain string) error {
	privateCaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicCaKey := privateCaKey.Public()
	caName := "ca." + domain
	subjectCa := pkix.Name{
		CommonName:         caName,
		OrganizationalUnit: []string{"kakoi"},
		Organization:       []string{"kakoi"},
		Country:            []string{"JP"},
	}

	caTpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               subjectCa,
		NotAfter:              time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		NotBefore:             time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	caCertificate, err := x509.CreateCertificate(rand.Reader, caTpl, caTpl, publicCaKey, privateCaKey)

	caCertFile, err := os.Create(filepath.Join(path, caName+".crt"))
	if err != nil {
		return err
	}
	defer caCertFile.Close()
	if err := pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertificate}); err != nil {
		return err
	}

	caKeyFile, err := os.Create(filepath.Join(path, caName+".key"))
	if err != nil {
		return err
	}
	defer caKeyFile.Close()
	if err := pem.Encode(caKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateCaKey)}); err != nil {
		return err
	}

	privateSslKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicSslKey := privateSslKey.Public()
	serverName := "server." + domain
	subjectSsl := pkix.Name{
		CommonName:         serverName,
		OrganizationalUnit: []string{"kakoi"},
		Organization:       []string{"kakoi"},
		Country:            []string{"JP"},
	}

	sslTpl := &x509.Certificate{
		SerialNumber: big.NewInt(123),
		Subject:      subjectSsl,
		NotAfter:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		NotBefore:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{serverName},
	}

	derSslCertificate, err := x509.CreateCertificate(rand.Reader, sslTpl, caTpl, publicSslKey, privateCaKey)
	if err != nil {
		return err
	}
	sslServerCrtFile, err := os.Create(filepath.Join(path, serverName+".crt"))
	if err != nil {
		return err
	}
	defer sslServerCrtFile.Close()
	if err := pem.Encode(sslServerCrtFile, &pem.Block{Type: "CERTIFICATE", Bytes: derSslCertificate}); err != nil {
		return err
	}

	sslServerKeyFile, err := os.Create(filepath.Join(path, serverName+".key"))
	if err != nil {
		return err
	}
	defer sslServerKeyFile.Close()
	if err := pem.Encode(sslServerKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateSslKey)}); err != nil {
		return err
	}

	privateClientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicClientKey := privateClientKey.Public()

	//Client Certificate
	clientName := "client." + domain
	subjectClient := pkix.Name{
		CommonName:         clientName,
		OrganizationalUnit: []string{"kakoi"},
		Organization:       []string{"kakoi"},
		Country:            []string{"JP"},
	}

	cliTpl := &x509.Certificate{
		SerialNumber: big.NewInt(456),
		Subject:      subjectClient,
		NotAfter:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		NotBefore:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	derClientCertificate, err := x509.CreateCertificate(rand.Reader, cliTpl, caTpl, publicClientKey, privateCaKey)
	if err != nil {
		return err
	}

	clientCertFile, err := os.Create(filepath.Join(path, clientName+".crt"))
	if err != nil {
		return err
	}
	defer clientCertFile.Close()
	if err := pem.Encode(clientCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: derClientCertificate}); err != nil {
		return err
	}

	clientKeyFile, err := os.Create(filepath.Join(path, clientName+".key"))
	if err != nil {
		return err
	}
	defer clientKeyFile.Close()
	if err := pem.Encode(clientKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateClientKey)}); err != nil {
		return err
	}
	return nil
}

func GenerateKeyPair(name, path string) error {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return err
	}
	publicKey, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	// encode
	privateKeyFile, err := os.Create(filepath.Join(path, name) + ".pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	if err := pem.Encode(privateKeyFile, &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	publicKeyFile, err := os.Create(filepath.Join(path, name) + ".pub")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicKey)
	if _, err := publicKeyFile.Write(pubKeyBytes); err != nil {
		return err
	}
	return nil
}

func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func generatePublicKey(publicKey *rsa.PublicKey) (ssh.PublicKey, error) {
	return ssh.NewPublicKey(publicKey)
}
