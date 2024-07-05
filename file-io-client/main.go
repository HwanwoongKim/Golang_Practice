package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
)

func main() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // testing only
		NextProtos:         []string{"h3", "http/1.1"},
	}
	url := "localhost:8080"
	req, _ := http.NewRequest("GET", url, nil)
	var buf bytes.Buffer
	req.Write(&buf)

	ctx := context.Background()
	connection, err := quic.DialAddr(ctx, url, tlsConfig, nil)
	if err != nil {
		println(err.Error())
		return
	}

	stream, err := connection.AcceptStream(context.Background())
	if err != nil {
		println(err.Error())
		return
	}

	buffer := make([]byte, 8)

	for {
		numBytes, err := stream.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		println(numBytes)

		log.Printf("Received: %s\n", buffer[:numBytes])
	}
}
