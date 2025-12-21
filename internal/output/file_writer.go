package output

import (
	"fmt"
	"os"
	"time"
)

type FileWriter struct {
	Dir      string
	Now      func() time.Time
	Location *time.Location
}

func NewDefaultFilteWriter() *FileWriter {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return &FileWriter{
		Dir:      "./releasePrMarkdown",
		Now:      time.Now,
		Location: loc,
	}
}

func (w *FileWriter) WriteReleaseMarkdown(markdown string) (string, error) {

	if err := os.MkdirAll(w.Dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}

	now := w.Now().In(w.Location)
	filePath := fmt.Sprintf(w.Dir+`/release_%s.md`, now.Format("20060102_1504"))

	if err := os.WriteFile(filePath, []byte(markdown), 0o644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", filePath, err)
	}
	return filePath, nil
}
