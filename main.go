package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	health "github.com/gopiesy/grpc-health-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port = flag.Int("port", 9111, "Port on which gRPC health server should listen for TCP conn.")
	root = flag.String("root", "./certs/cacert.pem", "root CA")
	cert = flag.String("cert", "./certs/server.pem", "server cert")
	key  = flag.String("key", "./certs/server.key", "server key")
)

func main() {
	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key. %s.", err)
	}

	// Load the CA certificate
	trustedCert, err := os.ReadFile(*root)
	if err != nil {
		log.Fatalf("Failed to load trusted certificate. %s.", err)
	}

	// Put the CA certificate to certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            certPool,
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}

	// Create a new TLS credentials based on the TLS configuration
	cred := credentials.NewTLS(tlsConfig)

	// Dial the gRPC server with the given credentials
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", *port), grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Unable to close gRPC channel. %s.", err)
		}
	}()

	// Create the request data
	request := &health.HealthCheckRequest{Service: "health"}

	// Create the gRPC client
	client := health.NewHealthClient(conn)
	response, err := client.Check(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to receive response. %s.", err)
	}

	// Print out response from server
	fmt.Println(response.Status)
}
