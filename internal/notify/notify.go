package notify

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

// Send shows a desktop notification with the given title and body,
// and plays a terminal beep.
func Send(title, body string) error {
	Beep()
	return sendDesktop(title, body)
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
		titleLiteral := strconv.Quote(title)
		bodyLiteral := strconv.Quote(body)
		script := fmt.Sprintf(
			`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;`+
				`$template = [Windows.UI.Notifications.ToastTemplateType]::ToastText02;`+
				`$xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent($template);`+
				`$xml.GetElementsByTagName("text")[0].AppendChild($xml.CreateTextNode(%s)) | Out-Null;`+
				`$xml.GetElementsByTagName("text")[1].AppendChild($xml.CreateTextNode(%s)) | Out-Null;`+
				`$toast = [Windows.UI.Notifications.ToastNotification]::new($xml);`+
				`[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Pomodoro").Show($toast);`,
			titleLiteral, bodyLiteral,
		)
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
		return cmd.Run()

	case "darwin":
		script := fmt.Sprintf("display notification %s with title %s", strconv.Quote(body), strconv.Quote(title))
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
