package main

import (
	"os"

	"symbiosis-cloud/symbiosis-k8s-controller/internal/cmd"
)

func main() {
	code := cmd.Run()
	os.Exit(code)
}
