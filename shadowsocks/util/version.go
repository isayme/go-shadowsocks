package util

import "fmt"

// version info
var (
	Name    = "@isayme/go-shadowsocks"
	Version = "unknown"
)

// PrintVersion print version
func PrintVersion() {
	fmt.Printf("name: %s\n", Name)
	fmt.Printf("version: %s\n", Version)
}
