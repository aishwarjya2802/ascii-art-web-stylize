package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Check if server started without arguments
	if len(os.Args) == 1 {
		fmt.Println("Server started, http://localhost:8000/")
	}
	staticFiles := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", staticFiles))
	
	// Define root route handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if requested path is root
		if r.URL.Path != "/" {
			http.Error(w, "Error 400: Path not found", http.StatusBadRequest)
			return
		}

		// Serve index.html file
		http.ServeFile(w, r, "index.html")
	})

	// Define process route handler
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		// Check if request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Error 405: Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form data
		r.ParseMultipartForm(0)

		// Get form values
		inputType := r.FormValue("inputType")
		textInput := r.FormValue("message")

		// Generate ASCII art and handle errors
		result, err := generateAsciiArt(inputType, textInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write ASCII art result to response
		fmt.Fprint(w, result)
	})

	// Start HTTP server and handle errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func generateAsciiArt(inputType, textInput string) (string, error) {
	// Define the text file based on input type
	var TextFile string
	switch inputType {
	case "standard":
		TextFile = "standard.txt"
	case "shadow":
		TextFile = "shadow.txt"
	case "thinkertoy":
		TextFile = "thinkertoy.txt"
	default:
		return "", fmt.Errorf("invalid input type")
	}

	// Open the text file
	fontFile, err := os.Open(TextFile)
	if err != nil {
		return "", err
	}
	defer fontFile.Close()

	// Read the text file and store in fontList
	fontList := []string{}
	scanner := bufio.NewScanner(fontFile)
	for scanner.Scan() {
		fontList = append(fontList, scanner.Text())
	}

	// Split the input text into lines
	splitText := strings.Split(textInput, "\n")

	// Generate the ASCII art
	var result string
	for _, line := range splitText {
		if len(line) == 0 {
			result += "\n"
			continue
		}

		for i := 0; i < 8; i++ {
			for _, char := range line {
				if int(char) >= 32 && int(char) <= 126 {
					initial := fontList[i+((int(char)-32)*9)+1]
					result += initial
				} 
			}
			result += "\n"
		}
	}

	return result, nil
}
