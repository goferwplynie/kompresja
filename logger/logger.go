package logger

import (
	"github.com/fatih/color"
)

var Verbose bool = true

func SetVerbose(v bool) {
	Verbose = v
}

func Log(info any) {
	if Verbose {
		blue := color.RGB(52, 64, 235)
		blue.Printf("[INFO] %v\n", info)
	}
}

func Warn(info any) {
	if Verbose {
		yellow := color.RGB(235, 205, 52)
		yellow.Printf("[WARN] %v\n", info)
	}
}
func Error(info any) {
	if Verbose {
		red := color.RGB(235, 67, 52)
		red.Printf("[ERROR] %v\n", info)
	}
}
func Cute(info any) {
	if Verbose {
		cute := color.RGB(235, 52, 235)
		cute.Printf("[INFO] %v\n", info)
	}
}
