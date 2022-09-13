# Notes about the implementation

[Sep 12, 2022]

- Today we're going to fix the strings to calc the row size

Here's an example:
 
  ```go
  package main

  import "fmt"

  func main() {
    var fixedSizeString [5]byte

    // copy(dest, source) - dest has to be a slice of bytes
    copy(fixedSizeString[:], "Hello world")

    fmt.Println(fixedSizeString)
    fmt.Println(string(fixedSizeString[:])) // string() also expects a slice of bytes
  }
  ```

Links
  - https://dlintw.github.io/gobyexample/public/memory-and-sizeof.html
  - https://www.includehelp.com/golang/read-a-structure-from-the-file.aspx

Sep 10 - In the tutorial the pages is a `void*` but such a type does not exist in Go. I tried to use `interface` but
this type only define methods. So I decided to create a new struct called `Page` to hold the rows in an array.

Found out that `binary` doesn work to convert complex structures. It's better to use `gob`. also exported fields starts
with Uppercase [https://stackoverflow.com/questions/65842245/what-does-the-error-binary-write-invalid-type-mean]

TODO: Have to find a way to fix string in Go because the ROW_SIZE has to be fixed.
  - https://forum.golangbridge.org/t/solved-string-size-of-20-character/15783/5
  - https://stackoverflow.com/questions/8039245/convert-string-to-fixed-size-byte-array-in-go
