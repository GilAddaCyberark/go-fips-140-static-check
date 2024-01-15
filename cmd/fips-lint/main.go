package main

import (
	fipslint "go-fips-140-static-check/pkg/fips-lint"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(fipslint.Analyzer)
}
