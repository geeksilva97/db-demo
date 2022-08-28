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
  PREPARE_SYNTAX_ERROR = iota
  PREPARE_UNRECOGNIZED_STATEMENT = iota
)

const (
  STATEMENT_INSERT = iota
  STATEMENT_SELECT = iota
)

const (
  EXECUTE_TABLE_FULL = iota
  EXECUTE_SUCCESS = iota
)

const ROW_SIZE = 1
const PAGE_SIZE = 4096
const TABLE_MAX_PAGES = 100
const TABLE_MAX_ROWS = 100
const ROWS_PER_PAGE = PAGE_SIZE / ROW_SIZE

type Row struct {
  id uint32
  username string
  email string
}

type Table struct {
  num_rows uint32
  pages [TABLE_MAX_PAGES]interface {}
}

type Statement struct {
  statement_type int
  row_to_insert Row
}

func print_row(row *Row) {
  fmt.Printf("(%d, %s, %s)\n", row.id, row.username, row.email)
}

func new_table() *Table {
  var table Table 
  table.num_rows = 0

  for i := 0; i < TABLE_MAX_PAGES; i++ {
    table.pages[i] = nil
  }

  return &table
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
    _, err := fmt.Sscanf(command, "insert %d %s %s", &stmt.row_to_insert.id, &stmt.row_to_insert.username, &stmt.row_to_insert.email)

    if err != nil {
      fmt.Println(err)
      return PREPARE_SYNTAX_ERROR
    }

    return PREPARE_SUCCESS
  }

  if command[:6] == "select" {
   stmt.statement_type = STATEMENT_SELECT 
   return PREPARE_SUCCESS
  }

  return PREPARE_UNRECOGNIZED_STATEMENT;
}

func execute_insert(statement *Statement, table *Table) int {
  if table.num_rows >= TABLE_MAX_ROWS {
    return EXECUTE_TABLE_FULL
  }

  row := &(statement.row_to_insert)
  table.pages[table.num_rows] = row
  table.num_rows += 1

  return EXECUTE_SUCCESS
}

func execute_select(statement *Statement, table *Table) int {
  var row *Row

  for i := 0; i < int(table.num_rows); i++ {
    row = table.pages[i].(*Row) // casts interface to a *Row
    print_row(row)
  }

  return EXECUTE_SUCCESS
}

func execute_statement(stmt *Statement, table *Table) int {
  switch(stmt.statement_type) {
    case STATEMENT_SELECT:
      return execute_select(stmt, table)

    case STATEMENT_INSERT:
      return execute_insert(stmt, table)
  }

  return -1
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
  var table *Table = new_table()

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
        case PREPARE_SYNTAX_ERROR:
          fmt.Printf("Syntax error. Could not parse element.\n")
          continue
        case PREPARE_UNRECOGNIZED_STATEMENT:
          fmt.Printf("Unrecognized keyword at start of '%s'.\n", command)
          continue
      }

      result := execute_statement(&statement, table)

      switch(result) {
        case EXECUTE_SUCCESS:
          fmt.Printf("Executed.\n")
          break

        case EXECUTE_TABLE_FULL:
          fmt.Printf("Error: Table full.\n")
          break
      }
    }
  }
