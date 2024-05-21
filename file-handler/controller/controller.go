package controller

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// SearchDirHandler godoc
// @Summary Search directory
// @Description Open a directory and list its contents.
// @Tags directory
// @Produce json
// @Param path query string true "Directory path"
// @Success 200 {array} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /searchdir [get]
func SearchDirHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing required query parameter: path", http.StatusBadRequest)
		return
	}

	dirList, err := os.Open(path)
	if err != nil {
		http.Error(w, "Failed to open directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dirList.Close()

	files, err := dirList.Readdirnames(-1)
	if err != nil {
		http.Error(w, "Failed to read directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, files)
}

// FileHandler godoc
// @Summary Create or append to a file
// @Description Create or append to a file in the specified directory.
// @Tags file
// @Accept json
// @Produce json
// @Param directory query string true "Directory path"
// @Param filename query string true "File name"
// @Param content body map[string]string true "Content to write"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /file [post]
func FileHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("directory")
	fileName := r.URL.Query().Get("filename")
	if dirPath == "" || fileName == "" {
		http.Error(w, "Missing required query parameters: directory or filename", http.StatusBadRequest)
		return
	}

	var reqBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	content, ok := reqBody["content"]
	if !ok {
		http.Error(w, "Missing required body parameter: content", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(dirPath, fileName)
	filePtr, exists, err := openOrCreateFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer filePtr.Close()

	if exists {
		log.Println("File already exists. Appending content.")
	} else {
		log.Println("File does not exist. Creating and writing content.")
	}

	if _, err := filePtr.WriteString(content + "\n"); err != nil {
		http.Error(w, "Error writing to file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	readFileLineByLine(filePath)
	w.Write([]byte("File updated successfully"))
}

// openOrCreateFile opens a file for appending or creates it if it doesn't exist
func openOrCreateFile(filePath string) (*os.File, bool, error) {
	dirPath := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	exists := findFile(dirPath, fileName)

	var file *os.File
	var err error
	if exists {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, exists, err
		}
	} else {
		file, err = os.Create(filePath)
		if err != nil {
			return nil, exists, err
		}
	}

	return file, exists, nil
}

// findFile checks if a file exists in the specified directory
func findFile(path, file string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println("Error reading directory:", err)
		return false
	}

	for _, f := range files {
		if file == f.Name() {
			return true
		}
	}

	return false
}

// readFileLineByLine reads a file line by line and logs each line
func readFileLineByLine(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file for reading:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}
}

// writeJSON writes the provided data as JSON to the response writer
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
