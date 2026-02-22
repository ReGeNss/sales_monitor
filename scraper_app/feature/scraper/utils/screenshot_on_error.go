package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const logsDir = "logs"

func SaveScreenshotOnError(page playwright.Page, err error, context string) {
	if page == nil || (err == nil && context == "") {
		return
	}

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("could not create logs directory: %v", err)
		return
	}

	safeContext := strings.ReplaceAll(strings.ReplaceAll(context, " ", "_"), "/", "-")
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
	log.Printf("screenshot saved to %s (error: %v)", path, err)
}
