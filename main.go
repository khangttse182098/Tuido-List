package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

type Task struct {
	name   string
	isDone bool
}

func handleChooseTasks(option *int) {
	for 1 == 1 {
		fmt.Print("(1) list out all tasks\n(2) add new task\n(3) Exit\n")
		fmt.Print("What do you want to do: ")
		fmt.Scan(option)

		// if they are then break
		if *option >= 1 && *option <= 3 {
			fmt.Print("\033[H\033[2J")
			break
		} else {
			fmt.Print("\033[H\033[2J")
			fmt.Println("You must enter option from 1-3!")
		}
	}
}

func listAllTasks(taskList []Task) bool {
	var doneTaskCount int
	if len(taskList) == 0 {
		fmt.Println("Todo list is empty!")
	} else {
		// count done tasks
		fmt.Println("Todo List")
		for i := range taskList {
			if taskList[i].isDone {
				fmt.Printf("%d. [x] %s\n", i+1, taskList[i].name)
				doneTaskCount += 1
			} else {
				fmt.Printf("%d. [ ] %s\n", i+1, taskList[i].name)
			}
		}
	}

	var chooseTask int
	fmt.Printf("Choose a task to perform action (1-%d) or 0 to quit: ", len(taskList))
	fmt.Scan(&chooseTask)

	//quit to menu
	if chooseTask == 0 {
		fmt.Print("\033[H\033[2J")
		return false
	}
	taskList[chooseTask-1].isDone = !taskList[chooseTask-1].isDone
	if taskList[chooseTask-1].isDone {
		fmt.Print("\033[H\033[2J")
		fmt.Println("Mark task as complete")
		doneTaskCount += 1
		// check if all tasks are finnished
		if doneTaskCount == len(taskList) {
			print(ascii)
		}
	} else {
		fmt.Print("\033[H\033[2J")
		fmt.Println("Mark task as incomplete")
	}
	fmt.Println("")
	return true
}

func addTask(taskList *[]Task, reader *bufio.Reader) {
	fmt.Print("Enter task name: ")
	newTask, _ := reader.ReadString('\n')
	newTask = strings.TrimSpace(newTask)
	*taskList = append(*taskList, Task{name: newTask, isDone: false})
	fmt.Print("\033[H\033[2J")
	fmt.Println("Added task successfully!")
	fmt.Println("")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var option int
	taskList := []Task{{"Mediate", false}, {"Wash dishes", false}, {"Drink honey lemon", false}}

	fmt.Println("Welcome to Terminal Todo List!")
	for 1 == 1 {
		// prompt the user to choose task
		handleChooseTasks(&option)

		// list out all tasks
		if option == 1 {
			if !listAllTasks(taskList) {
				continue
			}
		} else if option == 2 {
			addTask(&taskList, reader)
		} else if option == 3 {
			fmt.Println("Bye byeee!")
			break
		}
	}
}
