package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ascii = `
______________________
| Yayy you've finnish |
| all the tasks!!!    |
|_____________________|
 \/
   /\___/\
  ( U w U )
    > ^ <
`

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

// ----------------------------------------------------------

type Appstate int

const (
	Menu Appstate = iota
	TodoList
)

type model struct {
	TaskList     []string         `json:"taskList"`
	Selected     map[int]struct{} `json:"selected"`
	cursor       int
	isAdd        bool
	newTaskInput string
}

func initialModel() model {
	var modelObj model
	//read file
	if fileContent, err := os.ReadFile("tasks.json"); err == nil {
		//unmarshal json to taskList
		trimmedData := strings.TrimSpace(string(fileContent))
		if len(trimmedData) == 0 {
			modelObj = model{
				TaskList: []string{},
				Selected: make(map[int]struct{}),
			}
		} else {
			err = json.Unmarshal(fileContent, &modelObj)
			checkErr(err)
		}
	}
	return modelObj
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.isAdd {
	case true:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyRunes:
				m.newTaskInput += string(msg.Runes)
			case tea.KeySpace:
				m.newTaskInput += " "
			case tea.KeyBackspace:
				if len(m.newTaskInput) > 0 {
					m.newTaskInput = m.newTaskInput[:len(m.newTaskInput)-1]
				}
			case tea.KeyEnter:
				if len(m.newTaskInput) > 0 {
					m.TaskList = append(m.TaskList, m.newTaskInput)
					m.isAdd = false
					m.newTaskInput = ""
				}
			case tea.KeyEsc:
				saveTasksToFile(m)
				clearTerminal()
				return m, tea.Quit
			}
		}
	case false:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl + c", "q":
				saveTasksToFile(m)
				clearTerminal()
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.TaskList)-1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.Selected[m.cursor]
				if ok {
					delete(m.Selected, m.cursor)
				} else {
					m.Selected[m.cursor] = struct{}{}
				}
			case "o":
				if !m.isAdd {
					m.isAdd = true
					m.cursor = len(m.TaskList)
				}
			case "d":
				if len(m.TaskList) > 0 {
					m.TaskList = slices.Delete(m.TaskList, m.cursor, m.cursor+1)
					m.cursor--
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var taskStyle = lipgloss.NewStyle()

	// The header
	s := fmt.Sprintf("Todo List\n\n")

	// Iterate over our choices
	for i, task := range m.TaskList {
		// style for selected task
		taskStyle = lipgloss.NewStyle()
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
			taskStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#e78a4e"))
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += taskStyle.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, task)) + "\n"
	}

	// Is adding new task function on?
	if m.isAdd {
		taskStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#e78a4e"))
		s += taskStyle.Render(fmt.Sprintf("> [ ] %s", m.newTaskInput)) + "\n"
	}

	// The footer
	s += fmt.Sprintf("\nPress q to quit.\n")

	// Send the UI for rendering
	return s
}

func saveTasksToFile(modelObj model) {
	tasksJson, err := json.MarshalIndent(modelObj, "", " ")
	checkErr(err)
	err = os.WriteFile("tasks.json", tasksJson, 0644)
	checkErr(err)
}

func main() {
	clearTerminal()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
