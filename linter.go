package main

import (
	"io"

	actionlint "github.com/rhysd/actionlint"
)

type LinterError struct {
	Message string
	Line    int
	Column  int
}

// Lint is a wrapper function for actionlint
func lint(out []byte) ([]*LinterError, error) {
	var opts actionlint.LinterOptions
	linter, err := actionlint.NewLinter(io.Discard, &opts)
	if err != nil {
		return nil, err
	}

	alErrors, err := linter.Lint("<stdin>", out, nil)
	if err != nil {
		return nil, err
	}

	return convertActionlintErrors(alErrors), nil
}

// convertActionlintError converts actionlint.Error to LinterError
func convertActionlintErrors(alErrors []*actionlint.Error) []*LinterError {
	var lErrors []*LinterError
	for _, e := range alErrors {
		lErrors = append(lErrors, &LinterError{
			Message: e.Message,
			Line:    e.Line,
			Column:  e.Column,
		})
	}
	return lErrors
}
