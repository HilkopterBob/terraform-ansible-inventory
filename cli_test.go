package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runCLI(t *testing.T, stdin string, args ...string) (string, error) {
	cmdArgs := append([]string{"run", "main.go"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func TestCLIMissingInput(t *testing.T) {
	out, err := runCLI(t, "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(out, "Required flag \"input\" not set") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestCLIInvalidFormat(t *testing.T) {
	out, err := runCLI(t, "", "--input", "smoketest.json", "--format", "bad")
	if err == nil {
		t.Fatalf("expected error for bad format")
	}
	if !strings.Contains(out, "unknown inventory format") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestCLIReadStdin(t *testing.T) {
	data, err := os.ReadFile("smoketest.json")
	if err != nil {
		t.Fatalf("read smoketest: %v", err)
	}
	out, err := runCLI(t, string(data), "--input", "-", "--format", "ini")
	if err != nil {
		t.Fatalf("cli run err: %v\n%s", err, out)
	}
	if !strings.Contains(out, "[web]") || !strings.Contains(out, "ansible_host=192.168.1.10") {
		t.Fatalf("unexpected output: %s", out)
	}
}
