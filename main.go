package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

func paddingTop() {
	for range make([]int, 10) {
		fmt.Println("")
	}
}

func paddingLeft() string {
	var str string
	for range make([]int, 60) {
		str += " "
	}

	return str
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
				}
			case "d":
				if len(m.TaskList) > 0 {
					m.TaskList = slices.Delete(m.TaskList, m.cursor, m.cursor+1)
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// The header
	s := fmt.Sprintf("%sTodo List\n\n", paddingLeft())

	// Iterate over our choices
	for i, task := range m.TaskList {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s%s [%s] %s\n", paddingLeft(), cursor, checked, task)
	}

	// Is adding new task function on?
	if m.isAdd {
		s += fmt.Sprintf("  %s[ ] %s\n", paddingLeft(), m.newTaskInput)
	}

	// The footer
	s += fmt.Sprintf("\n%sPress q to quit.\n", paddingLeft())

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
	paddingTop()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
