package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
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
	Name   string `json:"name"`
	IsDone bool   `json:"isDone"`
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func saveTasksToFile(taskList []Task) {
	tasksJson, err := json.MarshalIndent(taskList, "", "  ")
	checkErr(err)
	err = os.WriteFile("tasks.json", tasksJson, 0644)
	checkErr(err)
}

func handleChooseTasks(option *int) {
	for 1 == 1 {
		fmt.Print(
			"(1) list out all tasks\n(2) add new task\n(3) Edit a task\n(4) Delete a task\n(5) Exit\n",
		)
		fmt.Print("What do you want to do: ")
		fmt.Scan(option)

		// if they are then break
		if *option >= 1 && *option <= 5 {
			clearTerminal()
			break
		} else {
			clearTerminal()
			fmt.Println("You must enter option from 1-5!")
		}
	}
}

func printTaskList(taskList []Task) {
	if len(taskList) == 0 {
		fmt.Println("Todo list is empty!")
	} else {
		// count done tasks
		fmt.Println("Todo List")
		for i := range taskList {
			if taskList[i].IsDone {
				fmt.Printf("%d. [x] %s\n", i+1, taskList[i].Name)
			} else {
				fmt.Printf("%d. [ ] %s\n", i+1, taskList[i].Name)
			}
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
			if taskList[i].IsDone {
				fmt.Printf("%d. [x] %s\n", i+1, taskList[i].Name)
				doneTaskCount += 1
			} else {
				fmt.Printf("%d. [ ] %s\n", i+1, taskList[i].Name)
			}
		}
	}

	// choose task to complete or incomplete
	var chooseTask int
	fmt.Printf("Choose a task to perform action (1-%d) or 0 to quit: ", len(taskList))
	fmt.Scan(&chooseTask)

	//quit to menu
	if chooseTask == 0 {
		clearTerminal()
		return false
	}

	// check if task is already done or not
	taskList[chooseTask-1].IsDone = !taskList[chooseTask-1].IsDone
	if taskList[chooseTask-1].IsDone {
		clearTerminal()
		fmt.Println("Mark task as complete")
		doneTaskCount += 1
		// check if all tasks are finnished
		if doneTaskCount == len(taskList) {
			print(ascii)
		}
	} else {
		clearTerminal()
		fmt.Println("Mark task as incomplete")
	}
	fmt.Println("")

	// save to json file
	saveTasksToFile(taskList)

	return true
}

func addTask(taskList *[]Task, reader *bufio.Reader) {
	fmt.Print("Enter task name: ")
	newTask, _ := reader.ReadString('\n')
	newTask = strings.TrimSpace(newTask)
	*taskList = append(*taskList, Task{Name: newTask, IsDone: false})
	clearTerminal()
	fmt.Println("Added task successfully!")
	fmt.Println("")

	// save to json file
	saveTasksToFile(*taskList)
}

func updateTask(taskList *[]Task, reader *bufio.Reader) bool {
	var chooseTaskIndex int
	printTaskList(*taskList)

	fmt.Printf("Choose a task from (1-%d) to update or 0 to quit: ", len(*taskList))
	fmt.Scan(&chooseTaskIndex)
	if chooseTaskIndex == 0 {
		return false
	}

	fmt.Printf("Enter new name: ")
	newTaskName, _ := reader.ReadString('\n')
	newTaskName = strings.TrimSpace(newTaskName)
	(*taskList)[chooseTaskIndex-1].Name = newTaskName

	clearTerminal()
	fmt.Println("Change the name successfully!")

	//save to json file
	saveTasksToFile(*taskList)

	return true
}

func deleteTask(taskList *[]Task) bool {
	var chooseTaskIndex int
	printTaskList(*taskList)

	fmt.Printf("Choose a task from (1-%d) to delete or 0 to quit: ", len(*taskList))
	fmt.Scan(&chooseTaskIndex)
	if chooseTaskIndex == 0 {
		return false
	}

	var confirmChoice byte
	fmt.Printf("Are you sure? (y/n): ")
	fmt.Scan(&confirmChoice)
	if confirmChoice == 'n' {
		return false
	}
	*taskList = slices.Delete(*taskList, chooseTaskIndex-1, chooseTaskIndex)

	//print out success message
	clearTerminal()
	fmt.Println("Delete the task successfully!")

	//save to json file
	saveTasksToFile(*taskList)

	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var option int
	var taskList []Task

	//read file
	fileContent, err := os.ReadFile("tasks.json")
	if err == nil {
		//unmarshal json to taskList
		err = json.Unmarshal(fileContent, &taskList)
		checkErr(err)
	}

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
			if !updateTask(&taskList, reader) {
				continue
			}
		} else if option == 4 {
			if !deleteTask(&taskList) {
				continue
			}
		} else if option == 5 {
			fmt.Println("Bye byeee!")
			break
		}
	}
}
