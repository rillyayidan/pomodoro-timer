package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rillyayidan/pomodoro/internal/config"
	"github.com/rillyayidan/pomodoro/internal/logger"
	"github.com/rillyayidan/pomodoro/internal/notify"
	"github.com/rillyayidan/pomodoro/internal/timer"
)

// ---- styles ----------------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			MarginBottom(1)

	phaseWorkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B"))

	phaseBreakStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4"))

	timeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFE66D")).
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1)

	statStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8E6CF")).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(1, 3).
			MarginTop(1)
)

// ---- model -----------------------------------------------------------------

// Model is the Bubble Tea model for the Pomodoro TUI.
type Model struct {
	cfg       config.Config
	state     timer.State
	progress  progress.Model
	quitting  bool
	statusMsg string
}

// New creates and initialises a new TUI model.
func New(cfg config.Config) Model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := timer.State{
		Phase:         timer.PhaseWork,
		TotalDuration: time.Duration(cfg.WorkDuration) * time.Minute,
		Remaining:     time.Duration(cfg.WorkDuration) * time.Minute,
		Running:       false,
	}
	return Model{cfg: cfg, state: s, progress: p}
}

// ---- tea.Model interface ---------------------------------------------------

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case " ": // space — start / pause
			m.state.Running = !m.state.Running
			if m.state.Running {
				m.statusMsg = ""
				return m, tick()
			}
			m.statusMsg = "Paused"
			return m, nil

		case "r": // reset current phase
			m.state.Remaining = m.state.TotalDuration
			m.state.Running = false
			m.statusMsg = "Reset"
			return m, nil

		case "s": // skip to next phase
			return m.advancePhase()
		}

	case timer.TickMsg:
		if !m.state.Running {
			return m, nil
		}
		m.state.Remaining -= time.Second
		if m.state.Remaining <= 0 {
			return m.advancePhase()
		}
		return m, tick()

	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		totalWork, count, _ := logger.TodaySummary()
		return fmt.Sprintf(
			"\n  Bye! Today: %d pomodoro(s), %d min of work focused 🍅\n\n",
			count, totalWork,
		)
	}

	// Phase label
	var phaseLabel string
	if m.state.Phase == timer.PhaseWork {
		phaseLabel = phaseWorkStyle.Render("🍅 " + m.state.Phase.String())
	} else {
		phaseLabel = phaseBreakStyle.Render("☕ " + m.state.Phase.String())
	}

	// Big countdown
	mins := int(m.state.Remaining.Minutes())
	secs := int(m.state.Remaining.Seconds()) % 60
	clock := timeStyle.Render(fmt.Sprintf("%02d:%02d", mins, secs))

	// Progress bar
	elapsed := m.state.TotalDuration - m.state.Remaining
	pct := float64(elapsed) / float64(m.state.TotalDuration)
	bar := m.progress.ViewAs(pct)

	// Status
	runStatus := "⏸  Paused  — press Space to start"
	if m.state.Running {
		runStatus = "▶  Running — press Space to pause"
	}
	if m.statusMsg != "" && !m.state.Running {
		runStatus = m.statusMsg
	}

	// Today stats
	totalWork, count, err := logger.TodaySummary()
	statsText := fmt.Sprintf("Today: %d 🍅  |  %d min focused", count, totalWork)
	if err != nil {
		statsText = "Today: stats unavailable"
	}
	stats := statStyle.Render(statsText)

	help := helpStyle.Render("[Space] start/pause   [r] reset   [s] skip   [q] quit")

	pomCount := fmt.Sprintf("Pomodoro #%d", m.state.PomodoroCount+1)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Pomodoro CLI — "+pomCount),
		phaseLabel,
		clock,
		bar,
		runStatus,
		stats,
		help,
	)

	return boxStyle.Render(inner)
}

// ---- helpers ---------------------------------------------------------------

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return timer.TickMsg(t)
	})
}

// advancePhase completes the current phase, logs it, notifies, then starts
// the next phase according to Pomodoro rules.
func (m Model) advancePhase() (tea.Model, tea.Cmd) {
	// Log completed phase
	durationMin := int(m.state.TotalDuration.Minutes())
	_ = logger.Append(m.state.Phase.String(), durationMin)

	// Notify
	var notifTitle, notifBody string
	switch m.state.Phase {
	case timer.PhaseWork:
		m.state.PomodoroCount++
		notifTitle = "Pomodoro done! 🍅"
		notifBody = "Time for a break."
	case timer.PhaseShortBreak, timer.PhaseLongBreak:
		notifTitle = "Break over! ☕"
		notifBody = "Back to work."
	}
	notify.Send(notifTitle, notifBody)

	// Determine next phase
	var nextPhase timer.Phase
	var nextDuration time.Duration
	switch m.state.Phase {
	case timer.PhaseWork:
		if m.state.PomodoroCount%m.cfg.LongBreakInterval == 0 {
			nextPhase = timer.PhaseLongBreak
			nextDuration = time.Duration(m.cfg.LongBreakDuration) * time.Minute
		} else {
			nextPhase = timer.PhaseShortBreak
			nextDuration = time.Duration(m.cfg.ShortBreakDuration) * time.Minute
		}
	default:
		nextPhase = timer.PhaseWork
		nextDuration = time.Duration(m.cfg.WorkDuration) * time.Minute
	}

	m.state.Phase = nextPhase
	m.state.TotalDuration = nextDuration
	m.state.Remaining = nextDuration
	m.state.Running = true
	m.statusMsg = ""

	return m, tick()
}
