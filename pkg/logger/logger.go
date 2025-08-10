// Package logger is a simple logger utility
package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var Color = map[string]string{
	"ERROR":   color.HiRedString("[✕]"),
	"INFO":    color.HiMagentaString("[i]"),
	"LOADING": color.HiBlueString("[*]"),
	"SUCCESS": color.HiGreenString("[✓]"),
	"WARNING": color.HiYellowString("[!]"),
}

func Error(msg string, err error) {
	fmt.Printf("%s %s", Color["ERROR"], msg)

	if os.IsNotExist(err) {
		color.Red(": [FILE/DIR DOESN'T EXIST]")
	} else if os.IsPermission(err) {
		color.HiYellow(": [PERMISSION DENIED]")
	} else if os.IsTimeout(err) {
		color.HiBlack(": [TIMEOUT]")
	} else if err == nil {
		// DD NOTHING
	} else {
		color.HiRedString(": [UNKNOWN]")
		fmt.Println(err)
	}
}

func Warning(msg string) {
	fmt.Printf("%s %s\n", Color["WARNING"], msg)
}

func Info(msg string) {
	fmt.Printf("%s %s\n", Color["INFO"], msg)
}

func Loading(msg string) {
	fmt.Printf("%s %s\n", Color["LOADING"], msg)
}

func SUCCESS(msg string) {
	fmt.Printf("%s %s\n", Color["SUCCESS"], msg)
}
