package pkg

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
)

var OverwriteAll bool

func generateEd25519Key() (ed25519.PrivateKey, []byte) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal("failed to generate Ed25519 key:", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatal("failed to marshal Ed25519 private key:", err)
	}

	return priv, der
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func randomSerial() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func fingerprintSHA256(der []byte) string {
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:])
}

func fingerprint(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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

func writePem(filename string, blockType string, bytes []byte) error {
	if _, err := os.Stat(filename); err == nil {
		if !OverwriteAll {
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
