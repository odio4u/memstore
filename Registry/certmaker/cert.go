package main

import (
	"log"

	"github.com/Purple-House/memstore/registry/certmaker/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := ParseConfig()
	log.Println("Starting certificate generation...")

	// CA
	ca, caKey, _, err := pkg.GenerateCA(cfg.CAName)
	pkg.Must(err)

	// Server cert
	_, err = pkg.GenerateServerCert(ca, caKey, cfg.ServerName, cfg.ServerIPs, cfg.ServerDNS)
	pkg.Must(err)
	// Optional client cert
	if cfg.WithClient {
		pkg.Must(pkg.GenerateClientCert(ca, caKey))
	}

	log.Println("All operations completed successfully.")
}
