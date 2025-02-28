package main

import (
	"fmt"      // For formatting strings
	"strings"  // For string manipulations
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// TableRow now has 5 player columns
type TableRow struct {
	RowID   string // Changed to string to handle "special" labels
	Player1 string
	Player2 string
	Player3 string
	Player4 string
	Player5 string
}

// Task struct represents a task in a stage.
type Task struct {
	Description template.HTML // A description of the task, allow safe raw HTML
}

// Stage represents a stage with a color, rows, and a list of tasks
type Stage struct {
	Color string     // Random color for the stage rows
	Rows  []TableRow // The data rows for the stage
	Tasks []Task     // The specific tasks for the stage
	Sum   string     // Placeholder for the sum (manually entered by users)
}

// generateTableData generates rows with empty player data
func generateTableData(rows int) []TableRow {
	table := make([]TableRow, rows)
	for i := 0; i < rows; i++ {
		table[i] = TableRow{
			RowID:   "Row " + strconv.Itoa(i+1), // Default RowID changed for clearer context
			Player1: "",
			Player2: "",
			Player3: "",
			Player4: "",
			Player5: "",
		}
	}
	return table
}

var color = []string{"pink", "orange", "blue", "green", "yellow", "random"}

func assignTasksToRows(stageTasks []Task, rows []TableRow) {
	for i := range rows {
		// Skip special rows
		if strings.HasPrefix(rows[i].RowID, "special") {
			continue
		}

		// Assign simplified task descriptions
		if i < len(stageTasks) {
			rows[i].RowID = string(stageTasks[i].Description)
		}
	}
}

// adjustSpecialIDs modifies specific RowIDs for special cases
func adjustSpecialIDs(stageIndex int, rows []TableRow) {
	if len(rows) >= 3 {
		// Change ID 3 in each stage
		rows[2].RowID = "special " + strconv.Itoa(stageIndex+1)
	}

	// In stage 3, change ID 4 as well
	if stageIndex == 2 && len(rows) >= 4 {
		rows[3].RowID = "special 4"
	}
}

// generateRandomFixedColor randomly chooses a color from a fixed set
func generateRandomFixedColor() string {
    colors := color // Use the color variable
	rand.Seed(time.Now().UnixNano())            // Seed the randomizer
	return colors[rand.Intn(len(colors))]       // Randomly pick one color
}

// generateTasks generates the tasks for each stage based on difficulty
func generateTasks(difficulty int) []Task {
	var tasks []Task
	switch difficulty {
	case 0: // Easy tasks, amounts 1-2
		tasks = []Task{
			{Description: template.HTML(encodeExactCount(rand.Intn(2)+1, generateRandomColor()))}, // Random amount (1-2)
			{Description: encodeNoColor(generateRandomColor())},                                  // Strikethrough task
			{Description: template.HTML(encodeTwoStrikethroughColors(generateRandomColor(), generateRandomColor()))}, // ~~colorx~~ & ~~colory~~
			{Description: template.HTML(encodeColorEquality(generateRandomColor(), generateRandomColor()))}, // New task "colorx = colory"
		}
	case 1: // Medium tasks, amounts 2-3
		tasks = []Task{
			{Description: template.HTML(encodeExactCount(rand.Intn(2)+2, generateRandomColor()))}, // Random amount (2-3)
			{Description: encodeAllDifferent(generateMultipleColors(3))},                          // All colors different
			{Description: template.HTML(encodeColorEquality(generateRandomColor(), generateRandomColor()))}, // New task "colorx = colory"
		}
	case 2: // Hard tasks, amounts fixed at 3, more complex rules
		tasks = []Task{
			{Description: template.HTML(encodeExactCount(3, generateRandomColor()))},              // Always 3
			{Description: encodeAllDifferent(generateMultipleColors(4))},                         // All colors different
			{Description: template.HTML(encodeColorEquality(generateRandomColor(), generateRandomColor()))}, // New task "colorx = colory"
			{Description: encodeNoColor(generateRandomColor())},                                  // Strikethrough task
		}
	default: // Default task description (fallback)
		tasks = []Task{
			{Description: template.HTML("Default Task")},
		}
	}

    shuffleTasks(tasks)

	return tasks
}

func encodeColorEquality(color1, color2 string) string {
	return fmt.Sprintf("%s = %s", color1, color2)
}

// encodeExactCount generates a task like "3*green" (e.g., exactly 3 green dice)
func encodeExactCount(count int, color string) template.HTML {
	return template.HTML(fmt.Sprintf("%d*%s", count, color))
}

// encodeNoColor generates a task like "<del>blue</del>" (e.g., no blue dice are allowed)
func encodeNoColor(color string) template.HTML {
	return template.HTML(fmt.Sprintf("~~%s~~", color))
}

// encodeTwoStrikethroughColors generates a task like "~~red~~ & ~~blue~~"
func encodeTwoStrikethroughColors(color1, color2 string) template.HTML {
	return template.HTML(fmt.Sprintf("~~%s~~ & ~~%s~~", color1, color2))
}


// encodeAllDifferent generates a task like "green!=yellow!=pink" (e.g., all colors must differ)
func encodeAllDifferent(colors []string) template.HTML {
	return template.HTML(strings.Join(colors, " != "))
}

// generateRandomColor selects a random color from the predefined list
func generateRandomColor() string {
	return color[rand.Intn(len(color))]
}

// generateMultipleColors generates a random selection of colors
func generateMultipleColors(count int) []string {
	predefinedColors := color
	rand.Shuffle(len(predefinedColors), func(i, j int) {
		predefinedColors[i], predefinedColors[j] = predefinedColors[j], predefinedColors[i]
	})
	selected := predefinedColors[:count]
	return selected
}

// shuffleTasks shuffles a slice of tasks randomly
func shuffleTasks(tasks []Task) {
	rand.Seed(time.Now().UnixNano()) // Properly seed randomness
	rand.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] }) // Shuffle in place
}

func renderTable(w http.ResponseWriter, r *http.Request) {
	// Parse the external HTML template file
	tmpl, err := template.ParseFiles("templates/structure.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	// Create 3 stages with varying rows and random colors
	stages := []Stage{}

	// Stage 1: 3 rows
	rows1 := generateTableData(3)
	stage1Tasks := generateTasks(0)          // Easy tasks
	assignTasksToRows(stage1Tasks, rows1)   // Link tasks to rows
	adjustSpecialIDs(0, rows1)              // Adjust IDs for special rows in stage 1
	stages = append(stages, Stage{
		Color: generateRandomFixedColor(),
		Rows:  rows1,
		Tasks: stage1Tasks, // Easy tasks for Stage 1
	})

	// Stage 2: 3 rows
	rows2 := generateTableData(3)
	stage2Tasks := generateTasks(1)          // Intermediate tasks
	assignTasksToRows(stage2Tasks, rows2)    // Link tasks to rows
	adjustSpecialIDs(1, rows2)               // Adjust IDs for special rows in stage 2
	stages = append(stages, Stage{
		Color: generateRandomFixedColor(),
		Rows:  rows2,
		Tasks: stage2Tasks, // Intermediate tasks for Stage 2
	})

	// Stage 3: 4 rows
	rows3 := generateTableData(4)
	stage3Tasks := generateTasks(2)          // Advanced tasks
	assignTasksToRows(stage3Tasks, rows3)    // Link tasks to rows
	adjustSpecialIDs(2, rows3)               // Adjust IDs for special rows in stage 3
	stages = append(stages, Stage{
		Color: generateRandomFixedColor(),
		Rows:  rows3,
		Tasks: stage3Tasks, // Advanced tasks for Stage 3
	})

	// Render the template with the stages data
	if err := tmpl.Execute(w, stages); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func main() {
	// Set up the server to listen on `/`.
	http.HandleFunc("/", renderTable)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
