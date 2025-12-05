package main

import (
	"flag"
	"strings"

	"github.com/Purple-House/memstore/registry/certmaker/pkg"
)

type Config struct {
	ServerIPs  []string
	ServerDNS  []string
	ServerName string
	CAName     string
	WithClient bool
	Force      bool
}

func ParseConfig() *Config {
	serverIPs := flag.String("ip", "127.0.0.1", "Comma-separated server IPs for certificate SAN")
	serverDNS := flag.String("dns", "", "Comma-separated DNS names for certificate SAN (optional)")
	serverName := flag.String("name", "grpc-server", "Common Name for server certificate")
	caName := flag.String("ca", "My Auto CA", "Common Name of the generated CA")
	withClient := flag.Bool("client", false, "Generate client certificate")
	force := flag.Bool("force", false, "Overwrite existing files without prompt")
	flag.Parse()

	cfg := &Config{
		ServerIPs:  strings.Split(*serverIPs, ","),
		ServerDNS:  []string{},
		ServerName: *serverName,
		CAName:     *caName,
		WithClient: *withClient,
		Force:      *force,
	}

	if *serverDNS != "" {
		cfg.ServerDNS = strings.Split(*serverDNS, ",")
	}

	pkg.OverwriteAll = cfg.Force

	return cfg
}
