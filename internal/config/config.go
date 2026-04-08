package config

// Config holds all timer durations (in minutes).
type Config struct {
	WorkDuration       int
	ShortBreakDuration int
	LongBreakDuration  int
	LongBreakInterval  int // after how many pomodoros a long break triggers
}

// Default returns the classic Pomodoro configuration.
func Default() Config {
	return Config{
		WorkDuration:       25,
		ShortBreakDuration: 5,
		LongBreakDuration:  15,
		LongBreakInterval:  4,
	}
}
