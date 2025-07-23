package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jakobmoellerdev/deferrlint/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(wd), "testdata")

	entries, err := os.ReadDir(filepath.Join(testdata, "src"))
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %s", err)
	}

	folders := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, filepath.Base(entry.Name()))
			t.Logf("Adding %q for analysis test", entry.Name())
		} else {
			t.Logf("Skipping non-folder %s", entry.Name())
		}
	}

	analysistest.Run(t, testdata, analyzer.Analyzer, folders...)
}
