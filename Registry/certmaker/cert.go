package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// ---------------------- GLOBALS ----------------------

var overwriteAll bool

// ---------------------- UTILITIES ----------------------

func writePem(filename string, blockType string, bytes []byte) error {
	// Check overwrite
	if _, err := os.Stat(filename); err == nil {
		if !overwriteAll {
			ok, err := promptOverwrite(filename)
			if err != nil {
				return err
			}
			if !ok {
				log.Printf("Skipped writing %s\n", filename)
				return nil
			}
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
	// 2048 is fine for now; bump to 3072/4096 if you want.
	return rsa.GenerateKey(rand.Reader, 2048)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func randomSerial() (*big.Int, error) {
	// 128-bit serial
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func fingerprintSHA256(der []byte) string {
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:])
}

// ---------------------- CERT GENERATION ----------------------

func generateCA(caName string) (*x509.Certificate, *rsa.PrivateKey, []byte, error) {
	log.Println("Generating CA certificate...")

	serial, err := randomSerial()
	if err != nil {
		return nil, nil, nil, err
	}

	caKey, err := newKey()
	if err != nil {
		return nil, nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   caName,
			Organization: []string{"MyOrg"},
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
	}

	der, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, err
	}

	must(writePem("ca.pem", "CERTIFICATE", der))
	must(writePem("ca-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caKey)))

	log.Println("✔ CA certificate generated.")
	log.Printf("CA SHA-256 fingerprint: %s\n", fingerprintSHA256(der))

	return ca, caKey, der, nil
}

func generateServerCert(
	ca *x509.Certificate,
	caKey *rsa.PrivateKey,
	serverName string,
	ipList []string,
	dnsList []string,
) ([]byte, error) {
	log.Println("Generating server certificate...")

	serial, err := randomSerial()
	if err != nil {
		return nil, err
	}

	serverKey, err := newKey()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, raw := range ipList {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		ip := net.ParseIP(raw)
		if ip == nil {
			return nil, fmt.Errorf("invalid server IP: %q", raw)
		}
		ips = append(ips, ip)
	}

	var dnsNames []string
	for _, raw := range dnsList {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		dnsNames = append(dnsNames, raw)
	}

	if len(ips) == 0 && len(dnsNames) == 0 {
		return nil, errors.New("at least one IP or DNS name must be provided for server SAN")
	}

	server := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"MyOrg"},
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		// For TLS server: DigitalSignature + KeyEncipherment are typical.
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		IPAddresses: ips,
		DNSNames:    dnsNames,
	}

	der, err := x509.CreateCertificate(rand.Reader, server, ca, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	must(writePem("server.pem", "CERTIFICATE", der))
	must(writePem("server-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(serverKey)))

	log.Println("✔ Server certificate generated.")
	log.Printf("Server SHA-256 fingerprint (pin this in agent): %s\n", fingerprintSHA256(der))

	return der, nil
}

func generateClientCert(ca *x509.Certificate, caKey *rsa.PrivateKey) error {
	log.Println("Generating client certificate...")

	serial, err := randomSerial()
	if err != nil {
		return err
	}

	key, err := newKey()
	if err != nil {
		return err
	}

	client := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "grpc-client",
			Organization: []string{"MyOrg"},
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(5, 0, 0),

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
	log.Printf("Client SHA-256 fingerprint: %s\n", fingerprintSHA256(der))

	return nil
}

// ---------------------- MAIN (CLI) ----------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// CLI FLAGS
	serverIPs := flag.String("ip", "127.0.0.1", "Comma-separated server IPs for certificate SAN")
	serverDNS := flag.String("dns", "", "Comma-separated DNS names for certificate SAN (optional)")
	serverName := flag.String("name", "grpc-server", "Common Name for server certificate")
	caName := flag.String("ca", "My Auto CA", "Common Name of the generated CA")
	withClient := flag.Bool("client", false, "Generate client certificate")
	force := flag.Bool("force", false, "Overwrite existing files without prompt")

	flag.Parse()
	overwriteAll = *force

	log.Println("Starting certificate generation...")

	// Split IPs and DNS names
	ipList := strings.Split(*serverIPs, ",")
	dnsList := []string{}
	if *serverDNS != "" {
		dnsList = strings.Split(*serverDNS, ",")
	}

	// 1 — CA
	ca, caKey, _, err := generateCA(*caName)
	must(err)

	// 2 — Server cert (Registry)
	_, err = generateServerCert(ca, caKey, *serverName, ipList, dnsList)
	must(err)

	// 3 — Client cert (optional, not needed for our current agent auth model)
	if *withClient {
		must(generateClientCert(ca, caKey))
	}

	log.Println("All operations completed successfully.")
}
