package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"crypto/tls"
	"log"

	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
)

func main() {
	cert, err := tls.LoadX509KeyPair("/Users/hwanwoong/minica.pem", "/Users/hwanwoong/minica-key.pem")
	if err != nil {
		log.Fatal("Error loading certificate:", err)
	}
	tlsConfig := &tls.Config{

		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,                       // For testing purposes only
		NextProtos:         []string{"h3", "http/1.1"}, // Enable QUIC and HTTP/3
	}
	listener, err := quic.ListenAddr("localhost:8080", tlsConfig, nil)
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

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	stream, err := connection.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	go handleRequest(stream)
}

func handleRequest(stream quic.SendStream) {
	file_bin := read_file()

	_, err := stream.Write(file_bin)

	if err != nil {
		log.Fatal(err)
	}
	err = stream.Close()
}

func read_file() []byte {
	var total_file []byte
	filepath := []string{"./grpc_structure.txt", "./test_restoring.txt"}

	file := make([]*os.File, len(filepath))

	for i, path := range filepath {
		var err error
		file[i], err = os.Open(path)

		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("File does not exist. Creating empty file...")

				// Create an empty file
				file[i], err = os.Create(path)
				if err != nil {
					fmt.Println("Error creating file:", err)
					return make([]byte, 0)
				}
				defer file[i].Close()
			} else {
				fmt.Println("Error opening file:", err)
				return make([]byte, 0)
			}
		} else {
			defer file[i].Close()

			// Get file size
			fileInfo, err := file[i].Stat()
			if err != nil {
				fmt.Println("Error getting file information:", err)
				return make([]byte, 0)
			}
			fileSize := fileInfo.Size()
			file_padding := make([]byte, 512)
			file_bin := make([]byte, fileSize)

			// Read file content into the byte slice
			err = binary.Read(file[i], binary.LittleEndian, &file_bin)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return make([]byte, 0)
			}

			copy(file_padding[:fileSize], file_bin)

			total_file = append(total_file, file_padding...)

			// Process the binary data here
			fmt.Printf("Read %d bytes from file:\n", fileSize)
		}
	}

	//fmt.Println("Total file binary sum is ... ")
	//fmt.Println(total_file)

	for _, file := range file {
		file.Close()
	}

	return total_file
}
