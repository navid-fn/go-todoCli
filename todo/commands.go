package todo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/fatih/color"
	"navid-fn.com/command-line-tool/db"
)

func ListTodo() {
	todos, err := db.GetallTodos()
	if err != nil {
		log.Fatal(err)
	}

	if len(todos) == 0 {
		yellow := color.New(color.FgYellow)
		yellow.Println("📝 No todos found. Add some todos to get started!")
		return
	}
	table := TurnTodoToTable(todos)

	fmt.Println()
	table.Render()
	fmt.Println()
	MarkCompleteTodo()
	DeleteToDo()
}

func AddTodo() {
	reader := bufio.NewReader(os.Stdin)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	// Set up signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer signal.Stop(c) // Cleanup signal handler when function returns

	// Channel for input completion
	inputDone := make(chan bool)

	var title, context string
	var err error

	// Title input goroutine
	go func() {
		green.Print("📝 Enter todo title: ")
		title, err = reader.ReadString('\n')
		inputDone <- true
	}()

	// Wait for either input completion or interrupt
	select {
	case <-inputDone:
		if err != nil {
			return
		}
	case <-c:
		yellow.Println("\n⚠️  Cancelled adding todo")
		return
	}

	title = strings.TrimSpace(title)
	if title == "" {
		yellow.Println("⚠️  Title cannot be empty")
		return
	}

	// Reset for context input
	go func() {
		green.Print("📝 Enter todo description: ")
		context, err = reader.ReadString('\n')
		inputDone <- true
	}()

	// Wait for either input completion or interrupt
	select {
	case <-inputDone:
		if err != nil {
			return
		}
	case <-c:
		yellow.Println("\n⚠️  Cancelled adding todo")
		return
	}

	context = strings.TrimSpace(context)
	if context == "" {
		yellow.Println("⚠️  Description cannot be empty")
		return
	}

	err = db.AddTodo(title, context)
	if err != nil {
		red.Println("❌ Something went wrong!")
		log.Fatal(err)
	} else {
		green.Println("✅ Todo added successfully!")
	}
}

func MarkCompleteTodo() {
	var choice string
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	yellow.Println("Want to mark a ToDo?")
	green.Println("1. ✅ YES")
	red.Println("2. ❌ NO")
	fmt.Scanln(&choice)
	switch choice {
	case "1":
		var todoId int
		yellow.Println("Select ID")
		fmt.Scanln(&todoId)
		err := db.MarkComplete(todoId)
		if err != nil {
			red.Println("❌ Error occurred")
			log.Fatal(err)
		}
		green.Printf("✅ Done! Todo #%d is marked as completed\n", todoId)
	case "2":
		return
	}
}

func CleanTodoTable() {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	err := db.CleanTable()
	if err != nil {
		red.Println("❌ Error cleaning todo table:", err)
		return
	}
	green.Println("✨ Todo table cleaned successfully!")
}

func SearchTitle() {
	reader := bufio.NewReader(os.Stdin)
	green := color.New(color.FgGreen)
	green.Print("📝 Enter todo Title for search: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		panic("ERROR OCCOURED")
	}
	title = strings.TrimSpace(title)
	todos, err := db.SearchTitle(title)
	if err != nil {
		panic("ERROR OCCOURED")
	}

	if len(todos) == 0 {
		yellow := color.New(color.FgYellow)

		fmt.Println()
		yellow.Println("📝 No todos found. Add some todos to get started!")
		return
	}
	table := TurnTodoToTable(todos)

	fmt.Println()
	table.Render()
	fmt.Println()
}

func DeleteToDo() {
	var choice string
	yellow := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	yellow.Println("Want to delete a ToDo?")
	green.Println("1. ✅ YES")
	red.Println("2. ❌ NO")

	fmt.Scanln(&choice)

	switch choice {
	case "1":
		var todoId int
		yellow.Println("Select ID")
		fmt.Scanln(&todoId)
		err := db.DeleteFromTodo(todoId)
		if err != nil {
			red.Println("❌ Error occurred")
			log.Fatal(err)
		}
		green.Printf("%d Deleted Successfully", todoId)
	case "2":
		return
	}
}
