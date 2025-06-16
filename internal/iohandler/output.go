package iohandler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Output dispatches to JSON, INI, TXT, or (added) Ansible-inventory formats.
func Output(data []map[string]interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return outputJSON(data)
	case "ini":
		return outputINI(data)
	case "txt":
		return outputTXT(data)
	case "ansible":
		// Not part of the original Output; route through the new helper:
		return OutputAnsibleInventory(data, "values.name", "values.variables.ip")
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

// --- New Ansible output below ---

// OutputAnsibleInventory takes parsed objects and two dot-paths,
// stripping CIDRs and emitting "hostname ansible_host=IP".
func OutputAnsibleInventory(
	objects []map[string]interface{},
	hostPath, ipPath string,
) error {
	for _, obj := range objects {
		hostname, err := lookupDotPath(obj, hostPath)
		if err != nil {
			return fmt.Errorf("host lookup %q: %w", hostPath, err)
		}
		ipCIDR, err := lookupDotPath(obj, ipPath)
		if err != nil {
			return fmt.Errorf("ip lookup %q: %w", ipPath, err)
		}
		// strip CIDR suffix
		ip := strings.SplitN(ipCIDR, "/", 2)[0]
		fmt.Printf("%s ansible_host=%s\n", hostname, ip)
	}
	return nil
}

// lookupDotPath walks a nested map by a dot-separated path,
// returning the final string value or an error.
func lookupDotPath(obj map[string]interface{}, path string) (string, error) {
	parts := strings.Split(path, ".")
	var cur interface{} = obj
	for _, p := range parts {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("not a map at %q", p)
		}
		cur, ok = m[p]
		if !ok {
			return "", fmt.Errorf("no key %q", p)
		}
	}
	s, ok := cur.(string)
	if !ok {
		return "", fmt.Errorf("value at %q is not a string", path)
	}
	return s, nil
}
