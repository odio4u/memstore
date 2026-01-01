package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"time"
)

func GenerateSelfSignedAgent(commonName string, dnsSAN []string) error {
	priv, privBytes := generateEd25519Key()
	pub := priv.Public().(ed25519.PublicKey)

	serial, _ := randomSerial()

	certTmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: commonName,
		},

		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},

		DNSNames: dnsSAN,
	}

	der, err := x509.CreateCertificate(rand.Reader, certTmpl, certTmpl, pub, priv)
	if err != nil {
		return err
	}

	Must(writePem("client.pem", "CERTIFICATE", der))
	Must(writePem("client-key.pem", "PRIVATE KEY", privBytes))

	log.Printf("Client Private Key Fingerprint: %s", fingerprint(privBytes))

	return nil
}
