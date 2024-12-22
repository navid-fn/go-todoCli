package todo

import (
	"os"
	"strconv"
	"time"

	"github.com/mshafiee/jalali"
	"github.com/olekukonko/tablewriter"
	"navid-fn.com/command-line-tool/db"
)

func TurnTodoToTable(todos []db.Todo) tablewriter.Table {

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
	return *table
}
