package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* ----------------- DATA MODELS ----------------- */

type status int

const (
	todo status = iota
	doing
	done
)

// Task structure with JSON tags for saving/loading
type Task struct {
	Status      status `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Column struct {
	Status status `json:"status"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Tasks  []Task `json:"tasks"`
}

type Board struct {
	columns  []Column
	focused  status
	cursor   int
	width    int
	height   int
	input    textinput.Model
	creating bool // State: Creating a new task
	editing  bool // State: Editing an existing task
}

/* ----------------- STYLING ----------------- */

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight)

	taskStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(subtle)

	selectedTaskStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("205")) // Pink

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

/* ----------------- PERSISTENCE ----------------- */

func (b *Board) save() {
	data, _ := json.MarshalIndent(b.columns, "", "  ")
	_ = os.WriteFile("board.json", data, 0644)
}

func (b *Board) load() {
	data, err := os.ReadFile("board.json")
	if err != nil {
		// No file found, keep default state
		return
	}
	_ = json.Unmarshal(data, &b.columns)
}

/* ----------------- INITIALIZATION ----------------- */

func NewBoard() *Board {
	ti := textinput.New()
	ti.Placeholder = "Task..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	board := &Board{
		focused:  todo,
		creating: false,
		editing:  false,
		input:    ti,
		columns: []Column{
			{Status: todo, Title: "To Do", Tasks: []Task{}},
			{Status: doing, Title: "In Progress", Tasks: []Task{}},
			{Status: done, Title: "Done", Tasks: []Task{}},
		},
	}

	// Try to load saved data
	board.load()

	return board
}

func (b Board) Init() tea.Cmd {
	return nil
}

/* ----------------- HELPER FUNCTIONS ----------------- */

func (b *Board) deleteTask() {
	col := &b.columns[b.focused]
	if len(col.Tasks) == 0 {
		return
	}
	// Remove task from slice
	col.Tasks = append(col.Tasks[:b.cursor], col.Tasks[b.cursor+1:]...)

	// Adjust cursor if it's now out of bounds
	if b.cursor >= len(col.Tasks) && b.cursor > 0 {
		b.cursor--
	}
	b.save()
}

func (b *Board) moveTask() {
	currentCol := &b.columns[b.focused]
	if len(currentCol.Tasks) == 0 {
		return
	}

	// 1. Get the task
	taskToMove := currentCol.Tasks[b.cursor]

	// 2. Delete from current column
	currentCol.Tasks = append(currentCol.Tasks[:b.cursor], currentCol.Tasks[b.cursor+1:]...)

	// 3. Determine next column
	nextStatus := b.focused + 1
	if nextStatus > done {
		nextStatus = todo
	}

	// 4. Update task status and append to next column
	taskToMove.Status = nextStatus
	b.columns[nextStatus].Tasks = append(b.columns[nextStatus].Tasks, taskToMove)

	// 5. Adjust cursor in the old column
	if b.cursor >= len(currentCol.Tasks) && b.cursor > 0 {
		b.cursor--
	}

	b.save()
}

/* ----------------- UPDATE LOOP ----------------- */

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return b, tea.Quit
		}

		// --- INPUT MODE (Creating or Editing) ---
		if b.creating || b.editing {
			switch msg.String() {
			case "enter":
				if b.creating && b.input.Value() != "" {
					// Create New
					b.columns[todo].Tasks = append(b.columns[todo].Tasks, Task{Status: todo, Title: b.input.Value()})
				} else if b.editing && b.input.Value() != "" {
					// Update Existing
					b.columns[b.focused].Tasks[b.cursor].Title = b.input.Value()
				}
				// Exit input mode
				b.creating = false
				b.editing = false
				b.input.Blur()
				b.input.Reset()
				b.save()
				return b, nil

			case "esc":
				b.creating = false
				b.editing = false
				b.input.Blur()
				b.input.Reset()
				return b, nil
			}
			// Let the textinput model handle the typing
			b.input, cmd = b.input.Update(msg)
			return b, cmd
		}

		// --- BOARD MODE (Navigation) ---
		switch msg.String() {
		case "q":
			return b, tea.Quit

		case "n":
			b.creating = true
			b.input.Focus()
			b.input.SetValue("")
			return b, textinput.Blink

		case "e":
			// Only allow edit if there is a task to edit
			if len(b.columns[b.focused].Tasks) > 0 {
				b.editing = true
				b.input.Focus()
				// Pre-fill input with current task title
				b.input.SetValue(b.columns[b.focused].Tasks[b.cursor].Title)
				return b, textinput.Blink
			}

		case "d", "backspace":
			b.deleteTask()

		case "E":
			// Export / Backup
			data, _ := json.MarshalIndent(b.columns, "", "  ")
			_ = os.WriteFile("backup_kanban.json", data, 0644)

		// Navigation
		case "h", "left":
			if b.focused > todo {
				b.focused--
				b.cursor = 0
			}
		case "l", "right":
			if b.focused < done {
				b.focused++
				b.cursor = 0
			}
		case "k", "up":
			if b.cursor > 0 {
				b.cursor--
			}
		case "j", "down":
			if b.cursor < len(b.columns[b.focused].Tasks)-1 {
				b.cursor++
			}
		case "enter", " ":
			b.moveTask()
		}
	}

	return b, nil
}

/* ----------------- VIEW (RENDER) ----------------- */

func (b Board) View() string {
	if b.width == 0 {
		return "loading..."
	}

	var cols []string
	colWidth := (b.width / 3) - 2

	// Render each column
	for i, col := range b.columns {
		var taskStrings []string
		for j, t := range col.Tasks {
			if b.focused == status(i) && b.cursor == j {
				// Highlight selected task
				taskStrings = append(taskStrings, selectedTaskStyle.Render("> "+t.Title))
			} else {
				taskStrings = append(taskStrings, taskStyle.Render(t.Title))
			}
		}

		// Highlight focused column border
		borderColor := lipgloss.Color("238") // Grey
		if b.focused == status(i) {
			borderColor = lipgloss.Color("62") // Purple
		}

		colView := columnStyle.
			Width(colWidth).
			BorderForeground(borderColor).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Render(col.Title),
					"",
					lipgloss.JoinVertical(lipgloss.Left, taskStrings...),
				),
			)
		cols = append(cols, colView)
	}

	boardView := lipgloss.JoinHorizontal(lipgloss.Left, cols...)

	// Helper footer
	helpString := helpStyle.Render("\n n: new • e: edit • d: del • E: backup • q: quit\n")

	// Render input overlay if needed
	if b.creating || b.editing {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			boardView,
			helpString,
			inputStyle.Render(b.input.View()),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, boardView, helpString)
}
