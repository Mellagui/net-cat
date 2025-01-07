package utils

import "fmt"

// Function to save logs
func savelogs(logs string) {
	_, err := logFile.WriteString(logs)
	if err != nil {
		fmt.Println(err)
	}
}
