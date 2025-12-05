package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"log"
	"net"
	"strings"
	"time"
)

func GenerateSelfSignedServer(commonName string, ipList, dnsList []string) ([]byte, error) {
	var ips []net.IP
	for _, raw := range ipList {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		ip := net.ParseIP(raw)
		if ip == nil {
			return nil, errors.New("invalid IP: " + raw)
		}
		ips = append(ips, ip)
	}

	if len(ips) == 0 && len(dnsList) == 0 {
		return nil, errors.New("server certificate requires at least one IP or DNS SAN")
	}

	priv, privBytes := generateEd25519Key()

	serial, _ := randomSerial()

	certTmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: commonName,
		},

		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		IPAddresses: ips,
		DNSNames:    dnsList,
	}
	pub := priv.Public().(ed25519.PublicKey)

	// Self-signed: certificate template = parent
	der, err := x509.CreateCertificate(rand.Reader, certTmpl, certTmpl, pub, priv)
	if err != nil {
		return nil, err
	}

	Must(writePem("server.pem", "CERTIFICATE", der))
	Must(writePem("server-key.pem", "PRIVATE KEY", privBytes))

	log.Printf("Server Private Key Fingerprint: %s", fingerprint(privBytes))

	return der, nil
}
