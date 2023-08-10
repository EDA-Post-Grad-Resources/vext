package main

import (
  "log"
  "database/sql"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
)

func main() {
  fmt.Println("Starting seed...")
  
  db, err := sql.Open("sqlite3", "../db/dev.sqlite3")

  if err != nil {
    log.Fatal(err)
  }

  defer db.Close()

 sts := `
DROP TABLE IF EXISTS todos;
CREATE TABLE todos(id INTEGER PRIMARY KEY, task TEXT);
INSERT INTO todos(task) VALUES('Buy milk');
INSERT INTO todos(task) VALUES('Buy eggs');
`
  if err != nil {
    log.Fatal(err)
  }

  _, err = db.Exec(sts)

  fmt.Println("Seed done!")
}
  
