package timer

import "time"

// Phase represents the current pomodoro phase.
type Phase int

const (
	PhaseWork Phase = iota
	PhaseShortBreak
	PhaseLongBreak
)

func (p Phase) String() string {
	switch p {
	case PhaseWork:
		return "Work"
	case PhaseShortBreak:
		return "Short Break"
	case PhaseLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

func (p Phase) IsBreak() bool {
	return p == PhaseShortBreak || p == PhaseLongBreak
}

func (p Phase) Icon() string {
	switch p {
	case PhaseWork:
		return "🍅"
	case PhaseShortBreak, PhaseLongBreak:
		return "☕"
	default:
		return ""
	}
}

// State holds the mutable runtime state of the timer.
type State struct {
	Phase         Phase
	Remaining     time.Duration
	TotalDuration time.Duration
	PomodoroCount int
	Running       bool
}

// TickMsg is sent every second by the Bubble Tea tick command.
type TickMsg time.Time

// DoneMsg is sent when the current phase finishes.
type DoneMsg struct {
	Phase Phase
}
