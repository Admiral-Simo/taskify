package main

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
func main() {




}
