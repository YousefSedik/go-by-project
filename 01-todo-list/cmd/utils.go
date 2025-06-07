package cmd

import (
	"fmt"
	"time"
)

func timesince(datetime string) string {
	time_now := time.Now()
	created, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return "unknown"
	}
	duration := time_now.Sub(created)
	if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	}
	return fmt.Sprintf("%s ago", duration.String())

}
