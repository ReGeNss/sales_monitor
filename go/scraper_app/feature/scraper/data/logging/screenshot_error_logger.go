package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"strings"
	"sync"
	"time"
)

const errorsLogFile = "errors.ndjson"

type errorRecord struct {
	Timestamp     string `json:"timestamp"`
	Error         string `json:"error"`
	Context       string `json:"context"`
	IndexOnPage   int    `json:"index_on_page,omitempty"`
	URL           string `json:"url,omitempty"`
	Screenshot    string `json:"screenshot"`
	ScreenshotURL string `json:"screenshot_url,omitempty"`
}

type screenshotErrorLogger struct {
	mu      sync.Mutex
	logsDir string
}

func NewScreenshotErrorLogger() gateway.ErrorLogger {
	logsDir := os.Getenv("SCRAPER_LOGS_DIR")
	if logsDir == "" {
		logsDir = "logs"
	}
	return &screenshotErrorLogger{logsDir: logsDir}
}

func (l *screenshotErrorLogger) LogError(err error, ctx gateway.ErrorContext) string {
	if err == nil && ctx.Context == "" {
		return ""
	}

	if mkErr := os.MkdirAll(l.logsDir, 0755); mkErr != nil {
		log.Printf("could not create logs directory: %v", mkErr)
		return ""
	}

	safeContext := strings.ReplaceAll(strings.ReplaceAll(ctx.Context, " ", "_"), "/", "-")
	if len(safeContext) > 50 {
		safeContext = safeContext[:50]
	}

	filename := fmt.Sprintf("error_%s_%s.png",
		time.Now().Format("2006-01-02_15-04-05"),
		safeContext)
	path := filepath.Join(l.logsDir, filename)

	record := errorRecord{
		Timestamp:   time.Now().Format(time.RFC3339),
		Error:       fmt.Sprintf("%v", err),
		Context:     ctx.Context,
		URL:         ctx.URL,
		IndexOnPage: ctx.Index,
		Screenshot:  filename,
	}
	if appendErr := l.appendRecord(record); appendErr != nil {
		log.Printf("could not save error record: %v", appendErr)
	}
	log.Printf("error logged to %s (error: %v)", path, err)
	return path
}

func (l *screenshotErrorLogger) appendRecord(record errorRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(filepath.Join(l.logsDir, errorsLogFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(record); err != nil {
		return err
	}
	return f.Sync()
}
