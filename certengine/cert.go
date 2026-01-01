package certengine

import (
	"log"

	"github.com/Purple-House/memstore/certengine/pkg"
)

func buildCerts() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting certificate generation...")

	cfg := ParseConfig()

	if cfg.WithClient {
		pkg.Must(pkg.GenerateSelfSignedAgent(cfg.ServerName, cfg.ServerDNS))
	} else {
		_, err := pkg.GenerateSelfSignedPublicFacing(cfg.ServerName, cfg.ServerIPs, cfg.ServerDNS)
		pkg.Must(err)
	}

	log.Println("All operations completed successfully.")
}
