package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single completed pomodoro session.
type Entry struct {
	Date        string `json:"date"`
	Phase       string `json:"phase"`
	Duration    int    `json:"duration_minutes"`
	CompletedAt string `json:"completed_at"`
}

// logFile returns the path to today's log file inside ~/.pomodoro/logs/.
func logFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".pomodoro", "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	today := time.Now().Format("2006-01-02")
	return filepath.Join(dir, today+".json"), nil
}

// Append adds an entry to today's log file.
func Append(phase string, durationMinutes int) error {
	path, err := logFile()
	if err != nil {
		return err
	}

	var entries []Entry
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &entries); err != nil {
			return err
		}
	}

	entries = append(entries, Entry{
		Date:        time.Now().Format("2006-01-02"),
		Phase:       phase,
		Duration:    durationMinutes,
		CompletedAt: time.Now().Format("15:04:05"),
	})

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// TodaySummary returns total work minutes and count of pomodoros logged today.
func TodaySummary() (totalWork int, pomodoroCount int, err error) {
	path, err := logFile()
	if err != nil {
		return 0, 0, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return 0, 0, err
	}
	for _, e := range entries {
		if e.Phase == "Work" {
			totalWork += e.Duration
			pomodoroCount++
		}
	}
	return totalWork, pomodoroCount, nil
}
