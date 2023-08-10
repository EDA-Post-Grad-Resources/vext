package db

import (
  "log"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type Todo struct  {
  Id int
  Task string
}

func GetTodos() []Todo {
  db, err := sql.Open("sqlite3", "./db/dev.sqlite3")

  if err != nil {
    log.Fatal(err)
  }

  defer db.Close()

  rows, err := db.Query("SELECT id, task FROM todos")

  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  var todos []Todo

  for rows.Next() {
    var todo Todo
    err := rows.Scan(&todo.Id, &todo.Task)

    if err != nil {
      log.Fatal(err)
    }

    todos = append(todos, todo)
  }

  err = rows.Err()

  if err != nil {
    log.Fatal(err)
  }

  return todos
}
