package util

import (
	"fmt"
	"time"
)

func FormatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	if minutes < 1 {
		return "Just now"
	}
	if minutes < 60 {
		return fmt.Sprintf("%dm ago", minutes)
	}
	if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}
	if days < 7 {
		return fmt.Sprintf("%dd ago", days)
	}
	return t.Format("Jan 2")
}
