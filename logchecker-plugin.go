package selectel

import (
	"github.com/beer-connoisseur/log-checker/english"
	"github.com/beer-connoisseur/log-checker/lowercase"
	"github.com/beer-connoisseur/log-checker/nosensitive"
	"github.com/beer-connoisseur/log-checker/nospecials"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("logchecker", New)
}

type Settings struct {
	Nosensitive struct {
		Words string `json:"words"`
	} `json:"nosensitive"`
}

type Plugin struct {
	settings Settings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: s}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.Nosensitive.Words != "" {
		if err := nosensitive.Analyzer.Flags.Set("words", p.settings.Nosensitive.Words); err != nil {
			return nil, err
		}
	}

	return []*analysis.Analyzer{
		lowercase.Analyzer,
		english.Analyzer,
		nospecials.Analyzer,
		nosensitive.Analyzer,
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
