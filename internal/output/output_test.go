package output

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jamesburton/latestdirs/internal/model"
)

func TestPrintTable_Format(t *testing.T) {
	results := []model.Result{
		{
			Path:      "/projects/foo",
			Timestamp: time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC),
			Unix:      1750000200,
			Source:    "filesystem",
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintTable(results)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Modified") {
		t.Error("expected header 'Modified'")
	}
	if !strings.Contains(output, "Directory") {
		t.Error("expected header 'Directory'")
	}
	if !strings.Contains(output, "2025-06-15 14:30:00") {
		t.Errorf("expected formatted date, got: %s", output)
	}
	if !strings.Contains(output, "/projects/foo") {
		t.Error("expected path in output")
	}
}

func TestPrintJSON_ValidJSON(t *testing.T) {
	results := []model.Result{
		{
			Path:      "/projects/bar",
			Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Unix:      1735689600,
			Source:    "filesystem",
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := PrintJSON(results)
	if err != nil {
		t.Fatal(err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []model.Result
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 result, got %d", len(parsed))
	}
	if parsed[0].Path != "/projects/bar" {
		t.Errorf("expected path '/projects/bar', got %q", parsed[0].Path)
	}
}
