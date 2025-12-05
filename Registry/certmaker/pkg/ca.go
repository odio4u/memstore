package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"time"
)

func GenerateCA(caName string) (*x509.Certificate, *rsa.PrivateKey, []byte, error) {
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

		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, err
	}

	Must(writePem("ca.pem", "CERTIFICATE", der))
	Must(writePem("ca-key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caKey)))

	log.Println("✔ CA certificate generated.")
	log.Printf("CA SHA-256 fingerprint: %s\n", fingerprintSHA256(der))

	return ca, caKey, der, nil
}
