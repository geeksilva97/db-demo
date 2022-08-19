package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
  META_COMMAND_SUCCESS = iota
  META_UNRECOGNIZED_COMMAND = iota
)

const (
  PREPARE_SUCCESS = iota
  PREPARE_UNRECOGNIZED_STATEMENT = iota
)

const (
  STATEMENT_INSERT = iota
  STATEMENT_SELECT = iota
)

type Statement struct {
  statement_type int
}

func do_meta_command(command *string) int {
  if *command == ".exit" {
    os.Exit(0)
  }
  return META_UNRECOGNIZED_COMMAND;
}

func prepare_statement(command string, stmt *Statement) int {
  if command[:6] == "insert" {
   stmt.statement_type = STATEMENT_INSERT 
   return PREPARE_SUCCESS
  }

  if command[:6] == "select" {
   stmt.statement_type = STATEMENT_SELECT 
   return PREPARE_SUCCESS
  }

  return PREPARE_UNRECOGNIZED_STATEMENT;
}

func execute_statement(stmt *Statement) {
  switch(stmt.statement_type) {
    case STATEMENT_SELECT:
      fmt.Printf("This is where we would select data\n")
      break

    case STATEMENT_INSERT:
      fmt.Printf("This is where we would insert data\n")
      break
  }
}

func print_prompt() {
  fmt.Printf("db_demo > ")
}

func read_input(input *string) {
  reader := bufio.NewReader(os.Stdin)
  command, _ := reader.ReadString('\n')
  *input = strings.TrimSuffix(command, "\n")
}

func main() {
  for {
      var command string

      print_prompt()
      read_input(&command)

      if command[0] == '.' {
        switch(do_meta_command(&command)) {
          case META_COMMAND_SUCCESS:
            continue

          case META_UNRECOGNIZED_COMMAND:
            fmt.Printf("Unrecognized command '%s'.\n", command)
            continue
          }
      }

      var statement Statement

      switch(prepare_statement(command, &statement)) {
        case PREPARE_SUCCESS:
          break
        case PREPARE_UNRECOGNIZED_STATEMENT:
          fmt.Printf("Unrecognized keyword at start of '%s'.\n", command)
          continue
      }

      execute_statement(&statement)
      fmt.Printf("Executed.\n")
    }
  }
