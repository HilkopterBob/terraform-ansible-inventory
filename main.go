package main

import (
	"flag"
	"fmt"
	"os"

	"terraform-ansible-inventory/internal/iohandler"
	"terraform-ansible-inventory/internal/parser"
)

func main() {
	inputPath := flag.String("input", "", "Path to input JSON file")
	format := flag.String("format", "json", "Output format: json, ini, txt")
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Input file required.")
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	results := parser.ExtractAnsibleHosts(data)

	err = iohandler.Output(results, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Output error: %v\n", err)
		os.Exit(1)
	}
}
