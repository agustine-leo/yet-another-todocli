package layout

import (
	"database/sql"
	"log"
	"log/slog"
	"strings"
	"time"
	"todocli/internal/database"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type todoList []database.Todo

func UI(db *sql.DB) {
	todos := database.NewTodo(db)

	// Initialize termui
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Mock data
	todoLists := todos.All()

	// Widgets
	help := widgets.NewParagraph()
	help.Text = "[q] Quit | [Arrow Keys|Vim Keybindings] Navigate | [Space] Toggle Complete | [a] Add Task | [d] Delete Task"
	termWidth, termHeight := ui.TerminalDimensions()
	help.SetRect(0, 0, termWidth, 3)

	list := widgets.NewList()
	list.Title = "TODO List"
	list.Rows = formatTaskList(todoLists)
	list.SelectedRow = 0
	list.SetRect(0, 3, termWidth, termHeight-3)
	list.SelectedRowStyle = ui.NewStyle(ui.ColorYellow, ui.ColorBlack, ui.ModifierBold)

	// Render initial layout
	ui.Render(help, list)

	var normalMode = true

	// Input widget
	inputBox := widgets.NewParagraph()
	inputBox.Title = "Add New Task"
	inputBox.Text = ""
	inputBox.BorderStyle = ui.NewStyle(ui.ColorYellow)
	inputBox.SetRect(10, 10, 50, 15)

	// Event loop
	for e := range ui.PollEvents() {
		switch e.Type {
		case ui.KeyboardEvent:
			if normalMode {
				switch e.ID {
				case "q":
					return // Quit the application
				case "<Down>", "j":
					if list.SelectedRow < len(todoLists)-1 {
						list.SelectedRow++
					}
				case "<Up>", "k":
					if list.SelectedRow > 0 {
						list.SelectedRow--
					}
				case "<Space>":
					// Toggle completion status
					if len(todoLists) == 0 {
						continue
					}
					idx := list.SelectedRow
					todoLists[idx].Completed = !todoLists[idx].Completed
					if todoLists[idx].Completed {
						todos.Update(todoLists[idx].ID, true)
					} else {
						todos.Update(todoLists[idx].ID, false)
					}
					list.Rows = formatTaskList(todoLists)

				case "a":
					// Add a new task
					// Create a new task
					newTodo := showInputPopup("Add New Todo", &strings.Builder{})
					normalMode = false

					slog.Debug("New TODO", "todo", newTodo)
					if newTodo != "" {
						todos.Create(newTodo, false)
						todoLists = todos.All()
						list.Rows = formatTaskList(todoLists)
						ui.Render(help, list)
					}
					normalMode = true
				case "m":
					// Modify the selected task description
					if len(todoLists) == 0 {
						continue
					}
					idx := list.SelectedRow
					currentDescription := strings.Split(todoLists[idx].Description, " ")[2:]
					stringifiedCurrentDescription := strings.Join(currentDescription, " ")
					stringsBuilderCurrentDescription := strings.Builder{}
					stringsBuilderCurrentDescription.WriteString(stringifiedCurrentDescription)
					newDescription := showInputPopup("Modify Todo", &stringsBuilderCurrentDescription)
					normalMode = false
					if newDescription != "" {
						todoLists[idx].Description = newDescription
						todos.Update(todoLists[idx].ID, todoLists[idx].Completed)
						list.Rows = formatTaskList(todoLists)
						ui.Render(help, list)
					}
					normalMode = true
				case "d":
					// Delete the selected task
					if len(todoLists) == 0 {
						continue
					}
					slog.Debug("Deleting task", "index", list.SelectedRow)
					idx := list.SelectedRow
					todos.Delete(todoLists[idx].ID)
					// Refresh the list after deletion
					todoLists = todos.All()
					list.Rows = formatTaskList(todoLists)
					if list.SelectedRow > 0 {
						list.SelectedRow--
					}
					ui.Render(help, list)
				}

				// Re-render UI
				ui.Render(help, list)
			}

		case ui.ResizeEvent:
			// Adjust widget positions and sizes when the terminal is resized
			payload := e.Payload.(ui.Resize)
			width, height := payload.Width, payload.Height

			// Adjust the layout and widget sizes
			help.SetRect(0, 0, width, 3)
			list.SetRect(0, 3, width, height-3)

			// Re-render after resizing
			ui.Render(help, list)
		}

	}
}

func formatTaskList(tasks []database.Todo) []string {
	rows := make([]string, len(tasks))
	for i, t := range tasks {
		status := "[ ]"
		if t.Completed {
			status = "[x]"
		}
		rows[i] = status + " " + t.Description
	}
	return rows
}

func showInputPopup(title string, inputText *strings.Builder) string {

	// Create an input box
	inputBox := widgets.NewParagraph()
	inputBox.Title = title
	inputBox.Text = ""
	inputBox.BorderStyle = ui.NewStyle(ui.ColorYellow)
	termWidth, termHeight := ui.TerminalDimensions()
	inputBox.SetRect(10, 10, termWidth-10, termHeight-20)
	inputBox.Text = inputText.String()

	// Render the input box
	ui.Render(inputBox)

	// Event loop for capturing user input
	for e := range ui.PollEvents() {
		switch e.Type {
		case ui.KeyboardEvent:
			switch e.ID {
			case "<Enter>":
				// Finish input on Enter
				var atTheMoment = time.Now().Format("2006-01-02 15:04:05")
				return atTheMoment + " " + inputText.String()
			case "<Space>":
				// Add space character
				inputText.WriteRune(' ')
			case "<Backspace>":
				// Handle backspace
				text := inputText.String()
				if len(text) > 0 {
					inputText.Reset()
					inputText.WriteString(text[:len(text)-1])
				}
			case "<Escape>":
				// Cancel input on Escape
				return ""
			default:
				// Add character to input
				slog.Debug("Key Pressed", "key", e.ID)
				if len(e.ID) == 1 { // Filter valid characters
					inputText.WriteRune(rune(e.ID[0]))
				}
			}

			// Update input box text
			inputBox.Text = inputText.String()
			ui.Render(inputBox)
		}
	}

	return ""
}
