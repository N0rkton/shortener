// Package main - Multichecker with several analyzers.
// usage:
// go install path/to/analyzer
// go vet -vettool=$(which analyzername) path/to/files
package main

import (
	"encoding/json"
	"github.com/alexkohler/nakedret"
	critic "github.com/go-critic/go-critic/checkers/analyzer"
	sqlclosecheck "github.com/ryanrolds/sqlclosecheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"os"
)

// Config — имя файла конфигурации.
const Config = `cmd/staticlint/config.json`

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string
}

func main() {
	data, err := os.ReadFile(Config)
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		ErrMainExit,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, sqlclosecheck.NewAnalyzer())
	mychecks = append(mychecks, nakedret.NakedReturnAnalyzer(5))
	mychecks = append(mychecks, critic.Analyzer)
	multichecker.Main(
		mychecks...,
	)
}
