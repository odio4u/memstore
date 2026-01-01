package main

import (
	"log"

	"github.com/Purple-House/memstore/certengine"
	"github.com/Purple-House/memstore/certengine/pkg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := certengine.ParseConfig()
	log.Println("Starting certificate generation...")

	// CA
	// ca, caKey, _, err := pkg.GenerateCA(cfg.CAName)
	// pkg.Must(err)
	// Arch update - skipping CA generation for now

	// client cert
	if cfg.WithClient {
		pkg.Must(pkg.GenerateSelfSignedClient(cfg.ServerName, cfg.ServerDNS))
	} else {
		_, err := pkg.GenerateSelfSignedServer(cfg.ServerName, cfg.ServerIPs, cfg.ServerDNS)
		pkg.Must(err)
	}

	log.Println("All operations completed successfully.")
}
