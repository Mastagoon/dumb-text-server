package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const dir = "./data"
const MAX_STR = 1024 * 4
const SWEEP_INTERVAL = time.Hour * 2

func init() {
	os.Mkdir(dir, 0755)
}

func handleGetFile(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	path := filepath.Join(dir, key)

	data, err := os.ReadFile(path)

	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	lines := strings.Split(string(data), "\n")

	reverseStrings(lines)

	w.Write([]byte(strings.Join(lines, "\n")))
}

func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func generateKey(w http.ResponseWriter, r *http.Request) {
	key := uuid.NewString()
	path := filepath.Join(dir, key)

	file, err := os.Create(path)
	if err != nil {
		http.Error(w, "failed to create file", http.StatusInternalServerError)
		return
	}

	file.Close()
	w.Write([]byte(key))
}

type postData struct {
	Data string `json:"data"`
}

func handleWriteToFile(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(string(data)) > MAX_STR {
		http.Error(w, "Input too long", http.StatusBadRequest)
		return
	}

	key := r.URL.Path[1:]
	path := filepath.Join(dir, key)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusNotFound)
		return
	}
	defer file.Close()

	// seperator
	now := time.Now().Format("15:04:05")
	_, err = file.WriteString(fmt.Sprintf("\n\n\n\n%s\n%s", string(data), fmt.Sprintf("@%s - :", now)))

	if err != nil {
		http.Error(w, "Failed to write", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func main() {
	// sweep old files
	sweepFiles()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /generate", generateKey)
	mux.HandleFunc("GET /{key}", handleGetFile)
	mux.HandleFunc("POST /{key}", handleWriteToFile)

	fmt.Println("listening on :6969")
	http.ListenAndServe(":6969", mux)
}

func sweepFiles() {
	fmt.Println("Sweeping files...")
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read dir", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("failed to stat file", err)
			continue
		}

		if (info.ModTime().Add(time.Hour * 24)).Before(time.Now()) {
			fmt.Printf("Deleting file: %s\n", path)
			os.Remove(path)
		}
	}
	time.AfterFunc(SWEEP_INTERVAL, sweepFiles)
}
