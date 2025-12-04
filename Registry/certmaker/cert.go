package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

// ---------------------- UTILITIES ----------------------

func writePem(filename string, blockType string, bytes []byte) error {
	// Check overwrite
	if _, err := os.Stat(filename); err == nil {
		ok, err := promptOverwrite(filename)
		if err != nil {
			return err
		}
		if !ok {
			log.Printf("Skipped writing %s\n", filename)
			return nil
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: blockType, Bytes: bytes})
}

func promptOverwrite(name string) (bool, error) {
	fmt.Printf("File %s already exists. Overwrite? (y/N): ", name)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	return input == "y\n" || input == "Y\n", nil
}

func newKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ---------------------- CERT GENERATION ----------------------

func generateCA(caName string) (*x509.Certificate, *rsa.PrivateKey, error) {
	log.Println("Generating CA certificate...")

	caKey, err := newKey()
	if err != nil {
		return nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   caName,
			Organization: []string{"MyOrg"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	must(writePem("ca.pem", "CERTIFICATE", der))
	must(writePem("ca-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caKey)))

	log.Println("✔ CA certificate generated.")
	return ca, caKey, nil
}

func generateServerCert(ca *x509.Certificate, caKey *rsa.PrivateKey, serverName, ip string) error {
	log.Println("Generating server certificate...")

	serverKey, err := newKey()
	if err != nil {
		return err
	}

	server := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"MyOrg"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(5, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP(ip)},
	}

	if server.IPAddresses[0] == nil {
		return errors.New("invalid server IP")
	}

	der, err := x509.CreateCertificate(rand.Reader, server, ca, &serverKey.PublicKey, caKey)
	if err != nil {
		return err
	}

	must(writePem("server.pem", "CERTIFICATE", der))
	must(writePem("server-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(serverKey)))

	log.Println("✔ Server certificate generated.")
	return nil
}

func generateClientCert(ca *x509.Certificate, caKey *rsa.PrivateKey) error {
	log.Println("Generating client certificate...")

	key, err := newKey()
	if err != nil {
		return err
	}

	client := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName:   "grpc-client",
			Organization: []string{"MyOrg"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(5, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	der, err := x509.CreateCertificate(rand.Reader, client, ca, &key.PublicKey, caKey)
	if err != nil {
		return err
	}

	must(writePem("client.pem", "CERTIFICATE", der))
	must(writePem("client-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key)))

	log.Println("✔ Client certificate generated.")
	return nil
}

// ---------------------- MAIN (CLI) ----------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// CLI FLAGS
	serverIP := flag.String("ip", "127.0.0.1", "Server IP for certificate SAN")
	serverName := flag.String("name", "grpc-server", "Common Name for server certificate")
	caName := flag.String("ca", "My Auto CA", "Common Name of the generated CA")
	withClient := flag.Bool("client", false, "Generate client certificate")

	flag.Parse()

	log.Println("Starting certificate generation...")

	// 1 — CA
	ca, caKey, err := generateCA(*caName)
	must(err)

	// 2 — Server cert
	must(generateServerCert(ca, caKey, *serverName, *serverIP))

	// 3 — Client cert (optional)
	if *withClient {
		must(generateClientCert(ca, caKey))
	}

	log.Println("All operations completed successfully.")
}
