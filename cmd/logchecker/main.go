package main

import (
	"github.com/beer-connoisseur/log-checker/english"
	"github.com/beer-connoisseur/log-checker/lowercase"
	"github.com/beer-connoisseur/log-checker/nosensitive"
	"github.com/beer-connoisseur/log-checker/nospecials"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		lowercase.Analyzer,
		english.Analyzer,
		nospecials.Analyzer,
		nosensitive.Analyzer,
	)
}
