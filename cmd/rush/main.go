package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rush <file.rush>")
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Printf("Rush interpreter - processing file: %s\n", filename)
	
	// TODO: Implement file reading and interpretation
	fmt.Println("Phase 1: Basic parsing - Coming soon!")
}