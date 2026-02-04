package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)



type status int

const (
	todo status = iota
	doing
	done
)


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
	creating bool 
	editing  bool 
}



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
				Foreground(lipgloss.Color("205")) 

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)



func (b *Board) save() {
	data, _ := json.MarshalIndent(b.columns, "", "  ")
	_ = os.WriteFile("board.json", data, 0644)
}

func (b *Board) load() {
	data, err := os.ReadFile("board.json")
	if err != nil {
		
		return
	}
	_ = json.Unmarshal(data, &b.columns)
}



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

	
	board.load()

	return board
}

func (b Board) Init() tea.Cmd {
	return nil
}



func (b *Board) deleteTask() {
	col := &b.columns[b.focused]
	if len(col.Tasks) == 0 {
		return
	}
	
	col.Tasks = append(col.Tasks[:b.cursor], col.Tasks[b.cursor+1:]...)

	
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

	
	taskToMove := currentCol.Tasks[b.cursor]

	
	currentCol.Tasks = append(currentCol.Tasks[:b.cursor], currentCol.Tasks[b.cursor+1:]...)

	
	nextStatus := b.focused + 1
	if nextStatus > done {
		nextStatus = todo
	}

	
	taskToMove.Status = nextStatus
	b.columns[nextStatus].Tasks = append(b.columns[nextStatus].Tasks, taskToMove)

	
	if b.cursor >= len(currentCol.Tasks) && b.cursor > 0 {
		b.cursor--
	}

	b.save()
}



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

		
		if b.creating || b.editing {
			switch msg.String() {
			case "enter":
				if b.creating && b.input.Value() != "" {
					
					b.columns[todo].Tasks = append(b.columns[todo].Tasks, Task{Status: todo, Title: b.input.Value()})
				} else if b.editing && b.input.Value() != "" {
					
					b.columns[b.focused].Tasks[b.cursor].Title = b.input.Value()
				}
				
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
			
			b.input, cmd = b.input.Update(msg)
			return b, cmd
		}

		
		switch msg.String() {
		case "q":
			return b, tea.Quit

		case "n":
			b.creating = true
			b.input.Focus()
			b.input.SetValue("")
			return b, textinput.Blink

		case "e":
			
			if len(b.columns[b.focused].Tasks) > 0 {
				b.editing = true
				b.input.Focus()
				
				b.input.SetValue(b.columns[b.focused].Tasks[b.cursor].Title)
				return b, textinput.Blink
			}

		case "d", "backspace":
			b.deleteTask()

		case "E":
			
			data, _ := json.MarshalIndent(b.columns, "", "  ")
			_ = os.WriteFile("backup_kanban.json", data, 0644)

		
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



func (b Board) View() string {
	if b.width == 0 {
		return "loading..."
	}

	var cols []string
	colWidth := (b.width / 3) - 2

	
	for i, col := range b.columns {
		var taskStrings []string
		for j, t := range col.Tasks {
			if b.focused == status(i) && b.cursor == j {
				
				taskStrings = append(taskStrings, selectedTaskStyle.Render("> "+t.Title))
			} else {
				taskStrings = append(taskStrings, taskStyle.Render(t.Title))
			}
		}

		
		borderColor := lipgloss.Color("238") 
		if b.focused == status(i) {
			borderColor = lipgloss.Color("62") 
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

	
	helpString := helpStyle.Render("\n n: new • e: edit • d: del • E: backup • q: quit\n")

	
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
