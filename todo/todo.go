package todo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mshafiee/jalali"
	"github.com/olekukonko/tablewriter"
	"navid-fn.com/command-line-tool/db"
)

func ListTodo() {
	todos, err := db.GetallTodos()
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Context", "Completed", "Created At"})
	
	// Table style configuration
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgHiMagentaColor},
	)

	for _, todo := range todos {
		tehran, _ := time.LoadLocation("Asia/Tehran")
		JcreatedAt := jalali.JalaliFromTime(todo.CreatedAt.In(tehran)).Format("%Y/%m/%d %H:%M")
		completed := "❌"
		if todo.Completed {
			completed = "✅"
		}
		table.Append([]string{
			strconv.Itoa(todo.Id),
			todo.Title,
			todo.Context,
			completed,
			JcreatedAt,
		})
	}

	if len(todos) == 0 {
		yellow := color.New(color.FgYellow)
		yellow.Println("📝 No todos found. Add some todos to get started!")
		return
	}

	fmt.Println()
	table.Render()
	fmt.Println()
	MarkCompleteTodo()
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
