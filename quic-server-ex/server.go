package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		return
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Println("Failed to create certificate:", err)
		return
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)

	if err != nil {
		log.Fatal("Error loading certificate:", err)
	}
	tlsConfig := &tls.Config{

		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,                                   // For testing purposes only
		NextProtos:         []string{"h3", "http/1.1", "ping/1.1"}, // Enable QUIC and HTTP/3
	}

	quicConfig := &quic.Config{
		Allow0RTT:       true,
		KeepAlivePeriod: time.Minute,
	}

	listener, err := quic.ListenAddr("localhost:8080", tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("QUIC server started on localhost:8080")

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go handleSession(connection)
	}
}

func handleSession(connection quic.Connection) {
	println("New connection established:", connection.RemoteAddr().String())
	for {
		state := connection.ConnectionState()
		fmt.Printf("state: %v\n", state)

		stream, err := connection.AcceptStream(context.Background())
		println("connect stream unblocked")

		if err != nil {
			println(err.Error())
			return
		}
		go handleRequest(stream)
	}
}

func handleRequest(stream quic.ReceiveStream) {
	buffer := make([]byte, 1024)
	numBytes, err := stream.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received: %s\n", buffer[:numBytes])
}
