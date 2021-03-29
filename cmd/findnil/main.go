package main

import (
	"github.com/gofriday/findnil"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(findnil.Analyzer) }
