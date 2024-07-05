package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
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
					return
				}
				defer file[i].Close()
			} else {
				fmt.Println("Error opening file:", err)
				return
			}
		} else {
			defer file[i].Close()

			// Get file size
			fileInfo, err := file[i].Stat()
			if err != nil {
				fmt.Println("Error getting file information:", err)
				return
			}
			fileSize := fileInfo.Size()

			// Create a byte slice to read the file
			data := make([]byte, fileSize)

			// Read file content into the byte slice
			_, err = file[i].Read(data)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			file_padding := make([]byte, 512)
			file_bin := make([]byte, fileSize)

			buf := bytes.NewBuffer(data)
			err = binary.Read(buf, binary.LittleEndian, &file_bin)
			copy(file_padding[:fileSize], file_bin)

			if err != nil {
				fmt.Println("binary read failed:", err)
			}
			fmt.Println(file_bin)
			fmt.Println(file_padding)

			total_file = append(total_file, file_padding...)

			// Process the binary data here
			fmt.Printf("Read %d bytes from file:\n%s\n", fileSize, data)
		}
	}

	fmt.Println("Total file binary sum is ... ")
	fmt.Println(total_file)

	for _, file := range file {
		file.Close()
	}
}
