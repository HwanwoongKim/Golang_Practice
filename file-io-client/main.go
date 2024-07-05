package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	println("Connection Established")

	stream, err := connection.AcceptStream(ctx)
	println("Stream Connection Accepted")
	if err != nil {
		println(err.Error())
		return
	}

	buffer := make([]byte, 64)
	file_bin := make([]byte, 0)
	bin_len := 0

	for {
		numBytes, err := stream.Read(buffer)
		println(numBytes)
		bin_len += numBytes
		log.Printf("Received: %s\n", buffer[:numBytes])
		file_bin = append(file_bin, buffer[:numBytes]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	file_n := bin_len / 512
	println(bin_len, file_n)
	files := make([][]byte, file_n)

	directory := "./"

	for i := 0; i < file_n; i++ {
		files[i] = make([]byte, 512)
		copy(files[i], file_bin[i*512:(i+1)*512])
		file_str := []string{directory, strconv.Itoa(i)}

		println(files[i], i)
		println(strings.Join(file_str, ""))
		makefile(files[i], strings.Join(file_str, ""))
	}
}

func makefile(file_bin []byte, filename string) {
	var err error
	var file *os.File

	file, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(file_bin)
}
