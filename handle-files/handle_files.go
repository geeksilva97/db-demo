package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

type myStruct struct {
  myBool bool // 1 byte
  myFloat float64 // 8 bytes
  myInt int32 // 4 bytes
}

type myOptimizedStruct struct {
  myFloat float64 // 8 bytes
  myInt int32 // 4 bytes
  myBool bool // 1 byte
}

type Row struct {
  Username [32]byte
  Email [200]byte
  Id uint32
}

func main() {
  a := myStruct{}
  fmt.Println(unsafe.Sizeof(a)) // 24

  b := myOptimizedStruct{}
  fmt.Println(unsafe.Sizeof(b)) // 16

  c := Row{}
  fmt.Printf("Size of Row with fized size string %v\n\n", unsafe.Sizeof(c)) // 236 bytes

  fi, err := os.Open("input.txt")

  if err != nil {
    panic(err)
  }

  defer func() {
    if err := fi.Close(); err != nil {
      panic(err)
    }
  }()

  fo, err := os.Create("output.txt")

  if err != nil {
    panic(err)
  }

  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()

  buf := make([]byte, 1024)

  for {
    n, err := fi.Read(buf)

    if err != nil && err != io.EOF {
      panic(err)
    }

    if n == 0 {
      break
    }

    text_read := string(buf[:n])

    fmt.Printf("Read string: \n\"%s\"\nSize: %v bytes?\nLen: %v\nSizeof: %v\n\n", text_read, int(unsafe.Sizeof(text_read)) + len(text_read), int(unsafe.Sizeof(text_read)), len(text_read))

    if _, err := fo.Write(buf[:n]); err != nil {
      panic(err)
    }
  }
}
