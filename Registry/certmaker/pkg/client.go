package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"time"
)

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
