package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "todo.db"

type Todo struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Context   string    `json:"context"`
	Completed bool      `json:"completed"`
}

func FileExists() bool {
	info, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Createdb() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the 'todo' table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS todo (
      id INTEGER PRIMARY KEY,
      title VARCHAR(20) NOT NULL,
      context TEXT,
      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      completed INTEGER NOT NULL DEFAULT 0
    );
  `)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database created...")
}

func Getdb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetallTodos() ([]Todo, error) {
	db, err := Getdb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, title, context, completed, created_at  FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Context, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func AddTodo(title string, context string) error {
	db, err := Getdb()
	if err != nil {
		return err
	}

	_, dber := db.Exec(`
    INSERT INTO todo (title, context)
    VALUES (?, ?);
  `, title, context)
	return dber
}

func MarkComplete(todoId int) error {
	db, err := Getdb()
	if err != nil {
		return err
	}
	_, dber := db.Exec(`
		UPDATE todo SET completed=true where id = ?;
	`, todoId)

	return dber

}

func CleanTable() error {
	db, err := Getdb()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM todo")
	return err
}

func SearchTitle(title string) ([]Todo, error) {
	db, err := Getdb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, title, context, completed, created_at  FROM todo where title = ?", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Context, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
