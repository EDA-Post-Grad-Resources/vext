package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct  {
  Id int
  Task string
}

type TodoDraft struct {
  Task string `form:"task"`
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

func AddTodo(newTodo *TodoDraft) int64 {
  db, err := sql.Open("sqlite3", "./db/dev.sqlite3")

  if err != nil {
    log.Fatal(err)
  }

  defer db.Close()

  result, err := db.Exec("INSERT INTO todos(task) VALUES(?);", newTodo.Task)
 
  id, err := result.LastInsertId()

  if err != nil {
    panic(err)
  }
  
  if err != nil {
        log.Fatal(err)
    }

  return id
}

func Delete(id int) {
 db, err := sql.Open("sqlite3", "./db/dev.sqlite3")

 if err != nil {
    log.Fatal(err)
  }

  defer db.Close()

  _, err = db.Exec("DELETE FROM todos WHERE id = ?;", id)

  if err != nil {
    log.Fatal(err)
  }

  return
}
