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
	sendDesktop(title, body)
}

// Beep writes the BEL character to stdout — works on all terminals.
func Beep() {
	fmt.Print("\a")
}

// sendDesktop dispatches a native notification based on the OS.
func sendDesktop(title, body string) {
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
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
		}

	case "darwin":

		script := fmt.Sprintf(`display notification "%s" with title "%s"`, body, title)
		cmd := exec.Command("osascript", "-e", script)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
		}

	default: // Linux / WSL
		cmd := exec.Command("notify-send", title, body)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
		}
	}
}
