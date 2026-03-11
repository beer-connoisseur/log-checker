package selectel_test

import (
	"testing"

	"github.com/beer-connoisseur/log-checker/english"
	"github.com/beer-connoisseur/log-checker/lowercase"
	"github.com/beer-connoisseur/log-checker/nosensitive"
	"github.com/beer-connoisseur/log-checker/nospecials"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()

	analysistest.Run(t, testdata, lowercase.Analyzer, "lowercase")
	analysistest.Run(t, testdata, english.Analyzer, "english")
	analysistest.Run(t, testdata, nospecials.Analyzer, "nospecials")
	analysistest.Run(t, testdata, nosensitive.Analyzer, "nosensitive")
}
