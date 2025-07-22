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
	analysistest.Run(t, testdata, analyzer.Analyzer, "ok", "fail")
}
