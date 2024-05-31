package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Admiral-Simo/task/db"
	"github.com/Admiral-Simo/task/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Orange = "\033[38;5;208m" // ANSI 256-color code for orange
	Green  = "\033[32m"
)

func listBetter(tasks []*models.Task) {
	// Print header
	fmt.Printf("%-5s  %-30s  %-3s\n", "ID", "Title", "Done")
	fmt.Println("-------------------------------------------")

	// Print tasks
	for _, task := range tasks {
		done := "No"
		if task.Done {
			done = "Yes"
		}

		var priorityColor string
		switch task.Priority {
		case "L":
			priorityColor = Yellow // Orange for low priority
		case "M":
			priorityColor = Orange // Yellow for medium priority
		case "H":
			priorityColor = Red // Red for high priority
		default:
			priorityColor = Reset
		}

		// Print with color
		if task.Done {
			fmt.Printf("%s%-5d  %-30s  %-3s\n", Green, task.ID, task.Title, done)
		} else {
			fmt.Printf("%s%-5d  %-30s  %-3s%s\n", priorityColor, task.ID, task.Title, done, Reset)
		}
	}
}

func mustGetTodaysTasks(client db.TaskStore) []*models.Task {
	tasks, err := client.GetAllTodaysTasks(context.TODO())
	// i need to sort tasks on the Done Field so that the "Yes" will be at the end
	var done []*models.Task
	var notDone []*models.Task
	for _, task := range tasks {
		if task.Done {
			done = append(done, task)
		} else {
			notDone = append(notDone, task)
		}
	}
	// sort the not done tasks by priority so that the high shows up first
	var h []*models.Task
	var m []*models.Task
	var l []*models.Task
	for _, task := range notDone {
		switch task.Priority {
		case "H":
			h = append(h, task)
		case "M":
			m = append(m, task)
		default:
			l = append(l, task)
		}
	}
	copy(tasks, h)
	copy(tasks[len(h):], m)
	copy(tasks[len(m)+len(h):], l)
	copy(tasks[len(m)+len(h)+len(l):], done)
	// end of the sorting
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}

func mustCreateTask(task *models.Task, client db.TaskStore) {
	if err := client.CreateTask(context.TODO(), task); err != nil {
		log.Fatal(err)
	}
}

func mustConnectDb() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	client, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// all migrations
	if err := client.AutoMigrate(&models.Task{}); err != nil {
		log.Fatal(err)
	}
	return client
}

func newTask(title string) *models.Task {
	return &models.Task{
		Title: title,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <add|list|prio|stats>\n", os.Args[0])
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	client := mustConnectDb()
	taskStore := db.NewGormTaskStore(client)

	command := os.Args[1]

	switch command {
	case "add":
		// need to combine all strings from 2 to infinity using
		title := strings.Join(os.Args[2:], " ")
		if len(title) < 4 {
			log.Fatal("the task should be 4 characters at least")
		}
		taskToAdd := newTask(title)
		mustCreateTask(taskToAdd, taskStore)
		fmt.Println("task added successfuly.")
	case "list":
		tasks := mustGetTodaysTasks(taskStore)
		listBetter(tasks)
	case "done":
		if len(os.Args) != 3 {
			fmt.Printf("usage: %s done <id>\n", os.Args[0])
			os.Exit(1)
		}
		taskID := os.Args[2]
		id, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			fmt.Printf("usage: %s done <id>\n", os.Args[0])
			os.Exit(1)
		}
		if err := taskStore.MarkDoneTask(context.TODO(), id); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("task %d successfuly done.\n", id)
	case "undone":
		if len(os.Args) != 3 {
			fmt.Printf("usage: %s undone <id>\n", os.Args[0])
			os.Exit(1)
		}
		taskID := os.Args[2]
		id, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			fmt.Printf("usage: %s undone <id>\n", os.Args[0])
			os.Exit(1)
		}
		if err := taskStore.MarkUnDoneTask(context.TODO(), id); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("task %d undone.\n", id)
	case "prio":
		if len(os.Args[2:]) != 2 {
			fmt.Printf("usage: %s <add|list|prio>\n", os.Args[0])
			os.Exit(1)
		}
		// if that works correctly we'll have 2 arguments
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("the id you have given is incorrect")
			os.Exit(1)
		}
		priority := strings.ToUpper(os.Args[3])
		err = taskStore.SetPriority(context.TODO(), id, priority)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("usage: %s <add|list|prio>\n", os.Args[0])
		os.Exit(1)

	}
}
