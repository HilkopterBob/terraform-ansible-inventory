package iohandler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Output(data []map[string]interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return outputJSON(data)
	case "ini":
		return outputINI(data)
	case "txt":
		return outputTXT(data)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputJSON(data []map[string]interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func outputINI(data []map[string]interface{}) error {
	for i, obj := range data {
		fmt.Printf("[host%d]\n", i)
		for k, v := range obj {
			fmt.Printf("%s = %v\n", k, v)
		}
		fmt.Println()
	}
	return nil
}

func outputTXT(data []map[string]interface{}) error {
	for i, obj := range data {
		fmt.Printf("Host %d:\n", i)
		for k, v := range obj {
			fmt.Printf("  %s: %v\n", k, v)
		}
		fmt.Println()
	}
	return nil
}
