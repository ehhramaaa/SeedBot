package helper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrettyLog(level, message string) {
	cleanMessage := strings.TrimSuffix(message, "\n")

	level = strings.ToUpper(level)

	var levelColor *color.Color
	switch level {
	case "INFO":
		levelColor = color.New(color.FgWhite) // Blue for INFO
	case "ERROR":
		levelColor = color.New(color.FgRed) // Red for ERROR
	case "WARNING":
		levelColor = color.New(color.FgYellow) // Yellow for WARNING
	case "INPUT":
		levelColor = color.New(color.FgCyan) // Cyan for INPUT
	case "SUCCESS":
		levelColor = color.New(color.FgGreen) // Cyan for INPUT
	default:
		levelColor = color.New(color.FgWhite) // White for default
	}

	// Print the log message with color
	if level == "INPUT" {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s", message)
	} else {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s\n", cleanMessage)
	}
}
