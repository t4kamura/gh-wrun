package main

import (
	"reflect"
	"testing"
)

func TestParseBranchesResult(t *testing.T) {
	tests := []struct {
		name string
		out  []byte
		want []string
	}{
		{
			name: "empty",
			out:  []byte{},
			want: []string{},
		},
		{
			name: "one branch",
			out:  []byte("  origin/main\n"),
			want: []string{"main"},
		},
		{
			name: "two branches",
			out:  []byte("  origin/main\n  origin/feature\n"),
			want: []string{"main", "feature"},
		},
		{
			name: "two branches with HEAD",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n"),
			want: []string{"main"},
		},
		{
			name: "two branches with HEAD and spaces",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n  origin/feature\n"),
			want: []string{"main", "feature"},
		},
		{
			name: "two branches with HEAD and spaces and tabs",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n  origin/feature\n\torigin/feature2\n"),
			want: []string{"main", "feature", "feature2"},
		},
		{
			name: "two branches with HEAD and spaces and tabs and newline",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n  origin/feature\n\torigin/feature2\n\n"),
			want: []string{"main", "feature", "feature2"},
		},
		{
			name: "two branches with HEAD and spaces and tabs and newline and spaces",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n  origin/feature\n\torigin/feature2\n\n  origin/feature3\n"),
			want: []string{"main", "feature", "feature2", "feature3"},
		},
		{
			name: "two branches with HEAD and spaces and tabs and newline and spaces and spaces",
			out:  []byte("  origin/HEAD -> origin/main\n  origin/main\n  origin/feature\n\torigin/feature2\n\n  origin/feature3\n  origin/feature4\n"),
			want: []string{"main", "feature", "feature2", "feature3", "feature4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBranchesResult(&tt.out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBranchesResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
