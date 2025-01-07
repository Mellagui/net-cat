package utils

import (
	"fmt"
	"os"
)


func exit() {
	var value string
	fmt.Println("Press D to stop server.")
	for value != "D" && value != "d" {
		fmt.Scanln(&value)
	}
	savelogs("Server Ended at: " + nowTime() + "\n")
	os.Exit(0)
}