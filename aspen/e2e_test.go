package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
)

type End2EndTestCase struct {
	fileName string
	stdout   string
}

func (tc *End2EndTestCase) Run(t *testing.T) {
	if tc == nil {
		return
	}

	cmd := exec.Command("./aspen", tc.fileName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("%s: could not run aspen file: %v", tc.fileName, err)
		return
	}

	stdout := out.String()

	if stdout != tc.stdout {
		t.Errorf("%s: expected stdout be be:\n%s\ngot:\n%s", tc.fileName, tc.stdout, stdout)
	}
}

func NewEnd2EndTestCase(file string, t *testing.T) *End2EndTestCase {
	source, err := OpenFile(file)

	if err != nil {
		t.Errorf("%s: %v", file, err)
	}

	tokens, err := ScanSource(source)

	if err != nil {
		t.Errorf("%s: %v", file, err)
	}

	stdout := tokens[0].value.(string)

	return &End2EndTestCase{fileName: file, stdout: stdout}
}

func TestEnd2End(t *testing.T) {
	// Build aspen binary

	cmd := exec.Command("go", "build")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not build aspen binary: %v", err)
	}

	matches, err := filepath.Glob("test_cases/e2e/*.aspen")

	if err != nil {
		t.Error("could not glob files")
		return
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := NewEnd2EndTestCase(match, t)
		tc.Run(t)
	}
}
