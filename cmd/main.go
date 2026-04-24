package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rillyayidan/pomodoro-timer/internal/config"
	"github.com/rillyayidan/pomodoro-timer/internal/logger"
	"github.com/rillyayidan/pomodoro-timer/internal/ui"
)

func main() {
	// ---- flags ---------------------------------------------------------------
	work := flag.Int("work", 25, "Work duration in minutes")
	short := flag.Int("short", 5, "Short break duration in minutes")
	long := flag.Int("long", 15, "Long break duration in minutes")
	silent := flag.Bool("silent", false, "Disable desktop notifications and beep")
	interval := flag.Int("interval", 4, "Pomodoros before a long break")
	statsOnly := flag.Bool("stats", false, "Show today's stats and exit")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
Pomodoro CLI — a terminal Pomodoro timer

Usage:
  pomodoro [flags]

Flags:`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
Controls (inside the TUI):
  Space   start / pause
  r       reset current phase
  s       skip to next phase
  q       quit`)
	}
	flag.Parse()

	// ---- stats-only mode -----------------------------------------------------
	if *statsOnly {
		totalWork, count, err := logger.TodaySummary()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading logs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Today's summary\n")
		fmt.Printf("  Pomodoros completed : %d 🍅\n", count)
		fmt.Printf("  Total focus time    : %d minutes\n", totalWork)
		return
	}

	// ---- validate flags ------------------------------------------------------
	if *work <= 0 || *short <= 0 || *long <= 0 || *interval <= 0 {
		fmt.Fprintln(os.Stderr, "All duration/interval values must be > 0")
		os.Exit(1)
	}

	// ---- build config & run --------------------------------------------------
	cfg := config.Default()
	cfg.WorkDuration = *work
	cfg.ShortBreakDuration = *short
	cfg.LongBreakDuration = *long
	cfg.LongBreakInterval = *interval
	cfg.Silent = *silent

	model := ui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
