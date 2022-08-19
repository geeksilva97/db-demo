package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
  for {
    var parsedCommand string;
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("db_demo > ")
    command, _ := reader.ReadString('\n')
    parsedCommand = strings.TrimSuffix(command, "\n")
    fmt.Println(command)

    if parsedCommand == ".exit" {
      break
    } else {
      fmt.Printf("Unrecognized command '%s'.\n", parsedCommand)
    }
  }
}
