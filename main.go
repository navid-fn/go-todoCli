package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"navid-fn.com/command-line-tool/db"
	"navid-fn.com/command-line-tool/todo"
)

func main() {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	// Set up signal handling for main menu
	cyan.Println("📝 TODO App!")
	if !db.FileExists() {
		db.Createdb()
	}

	for {
		fmt.Println("\n=========================")
		yellow.Println("Choose an option:")
		fmt.Println("=========================")
		green.Println("1. ➕ Add ToDo")
		green.Println("2. 📋 List ToDo")
		green.Println("3. 🗑️ Flush Table")
		green.Println("4. 🔍 Search Title")
		red.Println("5. Quit")
		fmt.Println("=========================")

		var choice string

		yellow.Print("Your Choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			todo.AddTodo()
		case "2":
			todo.ListTodo()
		case "3":
			todo.CleanTodoTable()
		case "4":
			todo.SearchTitle()
		case "5":
			cyan.Println("Have A nice Day! 👋")
			os.Exit(0)
		default:
			fmt.Printf("Your choice is: %s\n", choice)
			red.Println("Invalid choice! Please try again.")
		}
	}
}
