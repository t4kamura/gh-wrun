package main

import (
	"testing"

	actionlint "github.com/rhysd/actionlint"
)

func TestConvertActionLintErrors(t *testing.T) {
	alErrors := []*actionlint.Error{
		{
			Message:  "message1",
			Filepath: "filepath_1.yml",
			Line:     1,
			Column:   1,
			Kind:     "kind1",
		},
		{
			Message:  "message2",
			Filepath: "filepath_2.yml",
			Line:     2,
			Column:   2,
			Kind:     "kind2",
		},
	}
	want := []*LinterError{
		{
			Message: "message1",
			Line:    1,
			Column:  1,
		},
		{
			Message: "message2",
			Line:    2,
			Column:  2,
		},
	}

	lErrors := convertActionlintErrors(alErrors)

	if len(lErrors) != len(want) {
		t.Fatalf("len(lErrors) = %d, want %d", len(lErrors), len(want))
	}

	for i, e := range lErrors {
		if e.Message != want[i].Message {
			t.Errorf("lErrors[%d].Message = %s, want %s", i, e.Message, want[i].Message)
		}
		if e.Line != want[i].Line {
			t.Errorf("lErrors[%d].Line = %d, want %d", i, e.Line, want[i].Line)
		}
		if e.Column != want[i].Column {
			t.Errorf("lErrors[%d].Column = %d, want %d", i, e.Column, want[i].Column)
		}
	}

	if convertActionlintErrors(nil) != nil {
		t.Errorf("convertActionlintErrors(nil) != nil")
	}
}
