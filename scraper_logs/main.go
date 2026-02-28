package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	logsDir         = getEnv("SCRAPER_LOGS_DIR", "logs")
	baseURL         = getEnv("SCRAPER_LOGS_PUBLIC_URL", "")
	errorsLogFile   = "errors.ndjson"
	defaultPort     = "9092"
	apiErrorsPath   = "/api/errors"
	screenshotsPath = "/screenshots/"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type ErrorRecord struct {
	Timestamp     string `json:"timestamp"`
	Error         string `json:"error"`
	Context       string `json:"context"`
	IndexOnPage   int    `json:"index_on_page,omitempty"`
	URL           string `json:"url,omitempty"`
	Screenshot    string `json:"screenshot"`
	ScreenshotURL string `json:"screenshot_url"`
}

func main() {
	port := os.Getenv("SCRAPER_LOGS_PORT")
	if port == "" {
		port = defaultPort
	}

	http.HandleFunc(apiErrorsPath, handleErrors)
	http.HandleFunc(screenshotsPath, handleScreenshot)

	log.Printf("scraper-logs server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleErrors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := filepath.Join(logsDir, errorsLogFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.Write([]byte("[]"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var records []ErrorRecord
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rec ErrorRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		baseURL := os.Getenv("SCRAPER_LOGS_PUBLIC_URL")
		if baseURL == "" {
			baseURL = "http://localhost:9092"
		}
		rec.ScreenshotURL = strings.TrimSuffix(baseURL, "/") + "/screenshots/" + rec.Screenshot
		records = append(records, rec)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleScreenshot(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/screenshots/")
	if filename == "" || strings.Contains(filename, "..") {
		http.NotFound(w, r)
		return
	}

	path := filepath.Join(logsDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(data)
}
