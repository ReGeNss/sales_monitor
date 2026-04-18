package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

const errorsLogFile = "errors.ndjson"

var errorsLogMu sync.Mutex

func getLogsDir() string {
	if d := os.Getenv("SCRAPER_LOGS_DIR"); d != "" {
		return d
	}
	return "logs"
}

type ErrorContext struct {
	Context string 
	URL     string
	Index 	int
}

type ErrorRecord struct {
	Timestamp     string `json:"timestamp"`
	Error         string `json:"error"`
	Context       string `json:"context"`
	IndexOnPage   int	 `json:"index_on_page,omitempty"`
	URL           string `json:"url,omitempty"`
	Screenshot    string `json:"screenshot"`
	ScreenshotURL string `json:"screenshot_url,omitempty"`
}

// SaveScreenshotOnError зберігає скріншот і запис про помилку для передачі в scraper_logs API
func SaveScreenshotOnError(page playwright.Page, err error, ctx ErrorContext) {
	if page == nil || (err == nil && ctx.Context == "") {
		return
	}

	logsDir := getLogsDir()
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("could not create logs directory: %v", err)
		return
	}

	safeContext := strings.ReplaceAll(strings.ReplaceAll(ctx.Context, " ", "_"), "/", "-")
	if len(safeContext) > 50 {
		safeContext = safeContext[:50]
	}

	filename := fmt.Sprintf("error_%s_%s.png",
		time.Now().Format("2006-01-02_15-04-05"),
		safeContext)
	path := filepath.Join(logsDir, filename)

	_, screenshotErr := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
	if screenshotErr != nil {
		log.Printf("could not save screenshot on error: %v", screenshotErr)
		return
	}

	record := ErrorRecord{
		Timestamp:  time.Now().Format(time.RFC3339),
		Error:      fmt.Sprintf("%v", err),
		Context:    ctx.Context,
		URL:        ctx.URL,
		IndexOnPage: ctx.Index,
		Screenshot: filename,
	}
	if err := appendErrorRecord(record, logsDir); err != nil {
		log.Printf("could not save error record: %v", err)
	}
	log.Printf("screenshot saved to %s (error: %v)", path, err)
}

func appendErrorRecord(record ErrorRecord, logsDir string) error {
	errorsLogMu.Lock()
	defer errorsLogMu.Unlock()

	f, err := os.OpenFile(filepath.Join(logsDir, errorsLogFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(record); err != nil {
		return err
	}
	return f.Sync()
}
