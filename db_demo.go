package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"
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

const ROW_SIZE = int(unsafe.Sizeof(Row{}))
const PAGE_SIZE = 4096
const TABLE_MAX_PAGES = 100
const ROWS_PER_PAGE = PAGE_SIZE / ROW_SIZE
const TABLE_MAX_ROWS = 100

type Pager struct {
  file_descriptor *os.File
  file_length uint32
  pages [TABLE_MAX_PAGES]interface {}
}

type Row struct {
  id uint32 // 4 bytes
  username string // 18 bytes
  email string // 18 bytes
}

type Table struct {
  num_rows int
  pager *Pager
}

type Statement struct {
  statement_type int
  row_to_insert Row
}

func print_row(row *Row) {
  fmt.Printf("(%d, %s, %s)\n", row.id, row.username, row.email)
}

func db_open(filename string) *Table {
  pager := pager_open(filename)
  num_rows := int(pager.file_length ) / ROW_SIZE

  var table Table 
  table.pager = pager
  table.num_rows = num_rows

  return &table
}

func pager_open(filename string) *Pager {
  file_descriptor, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0755)

  if err != nil {
    fmt.Println("Unable to open file")
    os.Exit(-1)
  }

  var file_length int64
  file_length, err = file_descriptor.Seek(0, os.SEEK_END)

  var pager Pager
  pager.file_descriptor = file_descriptor
  pager.file_length = uint32(file_length)


  for i := 0; i < TABLE_MAX_PAGES; i++ {
    pager.pages[i] = nil
  }

  return &pager 
}

func get_page(pager *Pager, page_num uint32) interface{} {
  if page_num > TABLE_MAX_PAGES {
    fmt.Printf("Tried to fetch page number out of bounds. %d > %d\n", page_num, TABLE_MAX_PAGES)
    os.Exit(-1)
  } 

  if pager.pages[page_num] == nil {
    num_pages := pager.file_length / PAGE_SIZE

    if pager.file_length % PAGE_SIZE > 0 {
      num_pages += 1
    }

    if page_num <= num_pages {
      pager.file_descriptor.Seek(int64(page_num) * PAGE_SIZE, os.SEEK_SET) 
      bytes_read, err := pager.file_descriptor.Read(make([]byte, PAGE_SIZE))

      if err != nil {
        fmt.Printf("Error reading file: %d\n", err)
        os.Exit(-1)
      }

      pager.pages[page_num] = bytes_read
    }
  }

  return pager.pages[page_num]
}

func pager_flush(pager *Pager, page_num int, size int) {
  fmt.Printf("Flushing page %d\n", page_num)
  if pager.pages[page_num] == nil {
    fmt.Printf("Tried to flush a null page\n")
    os.Exit(-1)
  }

  page := pager.pages[page_num].(*Row)

  fmt.Printf("Id: %v\nNome: %v\nEmail: %v\n", page.id, page.username, page.email)

  writer := new(bytes.Buffer)

  binary.Write(writer, binary.LittleEndian, pager.pages[page_num])

  fmt.Printf("Bytes in buffer: %b\n", writer.Bytes())

  written_bytes, err := pager.file_descriptor.WriteAt(writer.Bytes(), int64(page_num) * PAGE_SIZE)

  if err != nil {
    fmt.Printf("Error writing\n")
    panic(err)
  }

  fmt.Printf("Written %v bytes \n", written_bytes)
}

func db_close(table *Table) {
  fmt.Println("Closing database...\n")
  pager := table.pager
  num_full_pages := table.num_rows / ROWS_PER_PAGE

  fmt.Printf("num_full_pages: %v\n", num_full_pages)

  for i := 0; i < num_full_pages; i++ {
    if pager.pages[i] != nil {
      pager_flush(pager, i, PAGE_SIZE)
      pager.pages[i] = nil
    }
  }

  num_additional_rows := table.num_rows % ROWS_PER_PAGE
  fmt.Printf("num_additional_rows: %v\n", num_additional_rows)
  if num_additional_rows > 0 {
    page_num := num_full_pages

    if pager.pages[page_num] != nil {
      pager_flush(pager, page_num, num_additional_rows * PAGE_SIZE)
      pager.pages[page_num] = nil
    }
  }

  err := pager.file_descriptor.Close()

  if err != nil {
    fmt.Printf("Error closing db file\n")
    panic(err)
  }

  for i := 0; i < TABLE_MAX_PAGES; i++ {
    pager.pages[i] = nil
  }
}

func do_meta_command(command *string, table *Table) int {
  if *command == ".exit" {
    db_close(table)
    os.Exit(0)
  }
  return META_UNRECOGNIZED_COMMAND;
}

func prepare_statement(command string, stmt *Statement) int {
  if command[:6] == "insert" {
    stmt.statement_type = STATEMENT_INSERT 
    _, err := fmt.Sscanf(command, "insert %d %s %s", &stmt.row_to_insert.id, &stmt.row_to_insert.username, &stmt.row_to_insert.email)

    if err != nil {
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
  table.pager.pages[table.num_rows] = row
  table.num_rows += 1

  return EXECUTE_SUCCESS
}

func execute_select(statement *Statement, table *Table) int {
  var row *Row

  for i := 0; i < int(table.num_rows); i++ {
    row = table.pager.pages[i].(*Row) // casts interface to a *Row
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
  if len(os.Args) < 2 {
    fmt.Println("Must supply a database filename")
    os.Exit(-1)
  }

  fmt.Printf("------- Database info ------ \nRow size: %v\nPage size: %v\nRows per page: %v\nMax rows: %v\nMax pages: %v\n\n", ROW_SIZE, PAGE_SIZE, ROWS_PER_PAGE, TABLE_MAX_ROWS, TABLE_MAX_PAGES)

  filename := os.Args[1]

  var table *Table = db_open(filename)

  for {
      var command string

      print_prompt()
      read_input(&command)

      if command[0] == '.' {
        switch(do_meta_command(&command, table)) {
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
