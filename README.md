# 🍅 Pomodoro CLI

A terminal Pomodoro timer built with Go and Bubble Tea.

## Features
- Interactive TUI with animated progress bar
- Desktop notifications (Windows toast / macOS / Linux)
- Terminal beep at end of each phase
- Optionally disable sound and notifications with `-silent`
- Next-phase hint in the TUI
- Daily session log saved to `~/.pomodoro/logs/YYYY-MM-DD.json`
- Custom durations via flags

## Project structure

```
pomodoro/
├── cmd/
│   └── main.go                  # Entry point, flag parsing
├── internal/
│   ├── config/config.go         # Duration config struct
│   ├── timer/timer.go           # Phase types & Bubble Tea messages
│   ├── ui/ui.go                 # Bubble Tea model (Update/View)
│   ├── logger/logger.go         # JSON session logger
│   └── notify/notify.go        # Desktop notifications + beep
├── go.mod
└── README.md
```

## Installation

```powershell
git clone https://github.com/rillyayidan/pomodoro-timer.git
cd pomodoro-timer
go mod tidy
go build -o pomodoro.exe ./cmd
```

Or install directly:

```powershell
go install github.com/rillyayidan/pomodoro-timer/cmd@latest
```

## Usage

```powershell
# Default (25/5/15 min)
./pomodoro.exe

# Custom durations
./pomodoro.exe -work 50 -short 10 -long 20 -interval 3

# Disable notifications and beep
./pomodoro.exe -silent

# Show today's stats without launching TUI
./pomodoro.exe -stats

# Print the version
./pomodoro.exe -version
./pomodoro.exe -v
```

## TUI controls

| Key     | Action              |
|---------|---------------------|
| `Space` | Start / Pause       |
| `r`     | Reset current phase |
| `s`     | Skip to next phase  |
| `q`     | Quit                |

## Session logs

Logs are stored at `~/.pomodoro/logs/YYYY-MM-DD.json`:

```json
[
  {
    "date": "2025-04-08",
    "phase": "Work",
    "duration_minutes": 25,
    "completed_at": "10:30:00"
  }
]
```
