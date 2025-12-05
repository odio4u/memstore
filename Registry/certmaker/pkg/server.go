package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func GenerateServerCert(
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

	if len(ips) == 0 && len(dnsList) == 0 {
		return nil, errors.New("at least one IP or DNS must be provided")
	}

	server := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"MyOrg"},
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		IPAddresses: ips,
		DNSNames:    dnsList,
	}

	der, err := x509.CreateCertificate(rand.Reader, server, ca, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	Must(writePem("server.pem", "CERTIFICATE", der))
	Must(writePem("server-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(serverKey)))

	log.Println("✔ Server certificate generated.")
	log.Printf("Server SHA-256 fingerprint: %s\n", fingerprintSHA256(der))

	return der, nil
}
