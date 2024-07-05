package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/klauspost/reedsolomon"
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
		bin_len += numBytes
		//log.Printf("Received: %s\n", buffer[:numBytes])
		file_bin = append(file_bin, buffer[:numBytes]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	file_n := bin_len / 512

	enc, err := reedsolomon.New(2, 1)
	if err != nil {
		log.Fatal(err)
	}

	rs_bin := make([][]byte, file_n)

	for i := 0; i < file_n; i++ {
		rs_bin[i] = append(rs_bin[i], file_bin[i*512:(i+1)*512]...)
	}

	rs_bin[0] = nil

	fmt.Println(rs_bin)
	err = enc.Reconstruct(rs_bin)

	if err != nil {
		log.Fatal(err)
	}

	var file_rs_decode []byte

	for i := 0; i < file_n; i++ {
		file_rs_decode = append(file_rs_decode, rs_bin[i]...)
	}

	fmt.Println(file_rs_decode)

	//println(bin_len, file_n)
	files := make([][]byte, file_n)

	directory := "./"

	for i := 0; i < file_n; i++ {
		files[i] = make([]byte, 512)
		copy(files[i], file_bin[i*512:(i+1)*512])
		file_str := []string{directory, strconv.Itoa(i)}

		//println(files[i], i)
		//println(strings.Join(file_str, ""))
		makefile(files[i], strings.Join(file_str, ""))
	}
}

func makefile(file_bin []byte, filename string) {
	var err error
	var file *os.File
	var real_file_len int

	for i := 0; i < len(file_bin); i++ {
		if file_bin[i] == byte(0) {
			real_file_len = i
			break
		}
	}

	file_real := make([]byte, real_file_len)
	file_real = file_bin[:real_file_len]

	file, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, &file_real)
}
