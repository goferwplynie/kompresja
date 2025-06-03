package logger

import (
	"github.com/fatih/color"
)

func Log(info any) {
	blue := color.RGB(52, 64, 235)
	blue.Printf("[INFO] %v\n", info)
}

func Warn(info any) {
	yellow := color.RGB(235, 205, 52)
	yellow.Printf("[WARN] %v\n", info)
}
func Error(info any) {
	red := color.RGB(235, 67, 52)
	red.Printf("[ERROR] %v\n", info)
}
func Cute(info any) {
	cute := color.RGB(235, 52, 235)
	cute.Printf("[INFO] %v\n", info)
}
