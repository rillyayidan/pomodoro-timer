package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Send shows a desktop notification with the given title and body,
// and plays a terminal beep.
func Send(title, body string) {
	Beep()
	if err := sendDesktop(title, body); err != nil {
		fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
	}
}

// Beep writes the BEL character to stdout — works on all terminals.
func Beep() {
	fmt.Print("\a")
}

// sendDesktop dispatches a native notification based on the OS.
func sendDesktop(title, body string) error {
	switch runtime.GOOS {
	case "windows":
		// PowerShell toast notification (no extra deps needed)
		script := fmt.Sprintf(
			`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;`+
				`$template = [Windows.UI.Notifications.ToastTemplateType]::ToastText02;`+
				`$xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent($template);`+
				`$xml.GetElementsByTagName('text')[0].AppendChild($xml.CreateTextNode('%s')) | Out-Null;`+
				`$xml.GetElementsByTagName('text')[1].AppendChild($xml.CreateTextNode('%s')) | Out-Null;`+
				`$toast = [Windows.UI.Notifications.ToastNotification]::new($xml);`+
				`[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('Pomodoro').Show($toast);`,
			title, body,
		)
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
		return cmd.Run()

	case "darwin":
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, body, title)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()

	default: // Linux / WSL
		if _, err := exec.LookPath("notify-send"); err != nil {
			return err
		}
		cmd := exec.Command("notify-send", title, body)
		return cmd.Run()
	}
}
