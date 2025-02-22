package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
	"strconv"
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

// Stage represents a stage with a color, rows, and a sum field
type Stage struct {
	Color string     // Random color for the stage rows
	Rows  []TableRow // The data rows for the stage
	Sum   string     // Placeholder for the sum (manually entered by users)
}

// generateTableData generates rows with empty player data
func generateTableData(rows int) []TableRow {
	table := make([]TableRow, rows)
	for i := 0; i < rows; i++ {
		table[i] = TableRow{
			RowID:   string(i + 1), // Default IDs as numbers (converted to string)
			Player1: "",            // Empty field
			Player2: "",
			Player3: "",
			Player4: "",
			Player5: "",
		}
	}
	return table
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
	colors := []string{"pink", "yellow", "blue"} // List of available colors
	rand.Seed(time.Now().UnixNano())            // Seed the randomizer
	return colors[rand.Intn(len(colors))]       // Randomly pick one color
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
	adjustSpecialIDs(0, rows1) // Adjust IDs for special rows in stage 1
	stages = append(stages, Stage{Color: generateRandomFixedColor(), Rows: rows1})

	// Stage 2: 3 rows
	rows2 := generateTableData(3)
	adjustSpecialIDs(1, rows2) // Adjust IDs for special rows in stage 2
	stages = append(stages, Stage{Color: generateRandomFixedColor(), Rows: rows2})

	// Stage 3: 4 rows
	rows3 := generateTableData(4)
	adjustSpecialIDs(2, rows3) // Adjust IDs for special rows in stage 3
	stages = append(stages, Stage{Color: generateRandomFixedColor(), Rows: rows3})

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