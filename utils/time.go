package utils

import "time"

func nowTime() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time in "YYYY-MM-DD HH:mm:ss"
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}
