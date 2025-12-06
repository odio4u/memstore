package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"runtime/debug"

	"github.com/Purple-House/memstore/registry/pkg/maps"
	memstore "github.com/Purple-House/memstore/registry/pkg/memstore"
	mapper "github.com/Purple-House/memstore/registry/proto"
	wal "github.com/Purple-House/memstore/registry/wal"
	walpb "github.com/Purple-House/memstore/registry/wal/proto"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func gracefulShutdown(server *grpc.Server) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Attempt graceful shutdown
	server.Stop()

}

func main() {
	fmt.Println("Registry Service for Ingress Tunnel")

	cert, err := tls.LoadX509KeyPair("certmaker/server.pem", "certmaker/server-key.pem")
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}

	servertLs := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,

		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
		},

		SessionTicketsDisabled:   true,
		PreferServerCipherSuites: true,

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		Renegotiation: tls.RenegotiateNever,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	port = fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			stack := string(debug.Stack())
			log.Printf("[PANIC RECOVERED] %v\nSTACK TRACE:\n%s", p, stack)
			return fmt.Errorf("internal server error")
		}),
	}

	waler, err := wal.OpenWAL()
	if err != nil {
		log.Fatalf("failed to open WAL: %v", err)
	}
	defer waler.Close()

	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(servertLs)),
		grpc.UnaryInterceptor(grpc_recovery.UnaryServerInterceptor(recoveryOpts...)),
		grpc.StreamInterceptor(grpc_recovery.StreamServerInterceptor(recoveryOpts...)),
	)

	store := memstore.NewMemStore()
	mapper.RegisterMapsServer(s, &maps.RPCMap{
		MemStore: store,
		WALer:    waler,
	})
	reflection.Register(s)

	_ = waler.Replay(func(wr *walpb.WalRecord) error {
		return wal.ApplyRecord(store, wr)

	})

	// Start the server
	fmt.Println("Server is running on port 50051")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	gracefulShutdown(s)

}
