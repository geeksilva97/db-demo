package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	filename := "output.txt"
  file_descriptor, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		panic(err)
	}

  var file_length int64
  file_length, err = file_descriptor.Seek(0, os.SEEK_END)

	// creates a slice to contain the read content
  buf := make([]byte, 15)

	// read the contet to the buf (in this case it reads 15 bytes) starting at the byte 13900
	num_bytes_read, err := file_descriptor.ReadAt(buf, 13900)

	if err != nil && err != io.EOF {
		panic(err)
	}


	fmt.Printf("Size of file %v\n", file_length)
	fmt.Printf("Num bytes read: %v\n", num_bytes_read)
	fmt.Printf("Content read:\n\n%v\n", string(buf[:num_bytes_read]))


	buf_write := []byte("writing at a specific point")
	// it overwrites the bytes that were in there
	// It increases the filez size if we write at a point greater than the file
	// num_written_bytes, err := file_descriptor.WriteAt(buf_write, 15000)

	num_written_bytes, err := file_descriptor.WriteAt(buf_write, 13900)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Num bytes written: %v\n", num_written_bytes)
	fmt.Printf("Content writen:\n\n%v\n", string(buf_write[:num_written_bytes]))

}
