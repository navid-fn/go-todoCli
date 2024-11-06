package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mshafiee/jalali"
	"navid-fn.com/command-line-tool/db"
)

func listTodo() {
	todos, err := db.GetallTodos()
	if err != nil {
		log.Fatal(err)
	}

	for _, todo := range todos {
		tehran, _ := time.LoadLocation("Asia/Tehran")
		JcreatedAt := jalali.JalaliFromTime(todo.CreatedAt.In(tehran)).Format("%Y/%m/%d %H:%M")
		fmt.Printf("ID %d, Title: %s, Context: %s, Completed: %t, CreatedAt: %s \n",todo.Id, todo.Title, todo.Context, todo.Completed,
			JcreatedAt)
	}
}

func addTodo() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter todo title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Enter todo description: ")
	context, _ := reader.ReadString('\n')

	err := db.AddTodo(title, context)
	if err != nil {
		fmt.Println("Something is Wrong!")
		log.Fatal(err)
	} else {
		fmt.Println("Added Successfully")
	}
}

func markCompleteTodo() {
	var choice string
	fmt.Println("Want to mark ToDo?")
	fmt.Println("1. YES")
	fmt.Println("2. NO")
	fmt.Scanln(&choice)
	switch choice {
	case "1":
		var todoId int
		fmt.Println("Select ID")
		fmt.Scanln(&todoId)
		err := db.MarkComplete(todoId)
		if err != nil {
			fmt.Println("Error Accouired")
			log.Fatal(err)
		}
		fmt.Printf("Done! id %d is set to True \n", todoId)
	case "2":
		return
	}

}

func main() {
	fmt.Println("TODO App!")
	if !db.FileExists() {
		db.Createdb()
	}
	for {
		fmt.Println("1. Add ToDo")
		fmt.Println("2. List ToDo")
		fmt.Println("3. Quit")

		var choice string
		fmt.Print("Your Choice:")
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			addTodo()
		case "2":
			listTodo()
			markCompleteTodo()
		case "3":
			fmt.Println("Have A nice Day!")
			return
		default:
			fmt.Println("Invalid choice! try again!")
		}
	}
}
