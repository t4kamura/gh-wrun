package subproc

import (
	"os"
	"path"
	"reflect"
	"testing"
)

const testDataDir = "../../testdata"

func TestParseWorkflows(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		want      []GhWorkflow
		expectErr bool
	}{
		{
			name:  "single workflow",
			input: ".github/workflows/test-1.yml\tactive\t12345678\n",
			want: []GhWorkflow{
				{
					Name:   ".github/workflows/test-1.yml",
					Status: "active",
					Id:     "12345678",
				},
			},
			expectErr: false,
		},
		{
			name: "multiple workflows",
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
			expectErr: false,
		},
		{
			name:      "just a string",
			input:     "invalid input",
			want:      nil,
			expectErr: true,
		},
		{
			name:      "Many tab separated values",
			input:     ".github/workflows/test-2.yml\tactive\t22345678\thoge\n",
			want:      nil,
			expectErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseWorkflows([]byte(test.input))
			if test.expectErr && err == nil {
				t.Errorf("Expected error but got nil\n")
			} else if !test.expectErr && err != nil {
				t.Errorf("Error parsing workflows: %s\n", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Expected is %v but got %v\n", test.want, got)
			}
		})
	}
}

func TestParseWorkflowInputs(t *testing.T) {
	// valid tests
	testCases := []struct {
		name          string
		inputFileName string
		want          []GhWorkflowInput
		expectErr     bool
	}{
		{
			name:          "parse valid types",
			inputFileName: "valid-all-types.yml",
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
			expectErr: false,
		},
		// want blank inputs
		{
			name:          "no inputs property",
			inputFileName: "valid-no-inputs-property.yml",
			want:          []GhWorkflowInput{},
			expectErr:     false,
		},
		{
			name:          "no inputs childs",
			inputFileName: "valid-no-inputs-childs.yml",
			want:          []GhWorkflowInput{},
			expectErr:     false,
		},
		{
			name:          "invalid format",
			inputFileName: "invalid-format.yml",
			want:          nil,
			expectErr:     true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.ReadFile(path.Join(testDataDir, test.inputFileName))
			if err != nil {
				t.Fatalf("Error reading file: %s\n", err)
			}

			got, err := parseWorkflowInputs(file)

			if test.expectErr && err == nil {
				t.Errorf("Expected error but got nil\n")
			} else if !test.expectErr && err != nil {
				t.Errorf("Error parsing workflows: %s\n", err)
			}

			if !test.expectErr && len(got) != len(test.want) {
				t.Errorf("The number of elements is different. Expected is %d but got %d\n", len(test.want), len(got))
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Expected is %v but got %v\n", test.want, got)
			}
		})
	}
}
