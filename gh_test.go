package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseWorkflows(t *testing.T) {
	validTests := []struct {
		input string
		want  []GhWorkflow
	}{
		{
			input: ".github/workflows/test-1.yml\tactive\t12345678\n",
			want: []GhWorkflow{
				{
					Name:   ".github/workflows/test-1.yml",
					Status: "active",
					Id:     "12345678",
				},
			},
		},
		{
			input: ".github/workflows/test-1.yml\tactive\t12345678\n" +
				".github/workflows/test-2.yml\tactive\t22345678\n",
			want: []GhWorkflow{
				{
					Name:   ".github/workflows/test-1.yml",
					Status: "active",
					Id:     "12345678",
				},
				{
					Name:   ".github/workflows/test-2.yml",
					Status: "active",
					Id:     "22345678",
				},
			},
		},
	}

	invalidTests := []string{
		"invalid input",
		".github/workflows/test-2.yml\tactive\t22345678\thoge\n",
	}

	for _, test := range validTests {
		got, err := parseWorkflows([]byte(test.input))
		if err != nil {
			t.Errorf("Error parsing workflows: %s\n", err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Expected is %v but got %v\n", test.want, got)
		}
	}

	for _, test := range invalidTests {
		if _, err := parseWorkflows([]byte(test)); err == nil {
			t.Errorf("Expected error but got nil\n")
		}
	}
}

func TestParseWorkflowInputs(t *testing.T) {
	validFile, err := os.ReadFile("./testdata/valid.yml")
	if err != nil {
		t.Fatalf("Error reading file: %s\n", err)
	}

	invalidFiles := []string{
		"./testdata/invalid-input-key-name.yml",
		"./testdata/invalid-format.yml",
	}

	validTests := []struct {
		input []byte
		want  []GhWorkflowInput
	}{
		{
			input: validFile,
			want: []GhWorkflowInput{
				{
					Name:        "string-all",
					Description: "string all test",
					Default:     "string default",
					Type:        GhWorkflowInputTypeString,
					Required:    true,
				},
				{
					Name:        "string-no-type-required",
					Description: "string no type & required test",
					Default:     "",
					Type:        GhWorkflowInputTypeString,
					Required:    false,
				},
				{
					Name:        "choice-all",
					Description: "choice test",
					Default:     "optionB",
					Type:        GhWorkflowInputTypeChoice,
					Required:    true,
					Options:     []string{"optionA", "optionB", "optionC"},
				},
				{
					Name:        "boolean-default-true",
					Description: "boolean default true test",
					Default:     "true",
					Type:        GhWorkflowInputTypeBoolean,
					Required:    true,
				},
				{
					Name:        "boolean-default-false",
					Description: "boolean default false test",
					Default:     "false",
					Type:        GhWorkflowInputTypeBoolean,
					Required:    true,
				},
				{
					Name:        "boolean-no-default",
					Description: "boolean no default test",
					Default:     "",
					Type:        GhWorkflowInputTypeBoolean,
					Required:    true,
				},
				{
					Name:        "environment-all",
					Description: "environment test",
					Default:     "production",
					Type:        GhWorkflowInputTypeEnvironment,
					Required:    true,
				},
			},
		},
	}

	for _, test := range validTests {
		got, err := parseWorkflowInputs(test.input)
		if err != nil {
			t.Errorf("Error parsing workflows: %s\n", err)
		}

		if len(got) != len(test.want) {
			t.Errorf("The number of elements is different. Expected is %d but got %d\n", len(test.want), len(got))
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Expected is %v but got %v\n", test.want, got)
		}
	}

	for _, file := range invalidFiles {
		invalidFile, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Error reading file: %s\n", err)
		}
		_, err = parseWorkflowInputs(invalidFile)
		if err == nil {
			t.Errorf("Expected error but got nil\n")
		}
	}
}
