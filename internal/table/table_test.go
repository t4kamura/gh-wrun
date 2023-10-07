package table

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	input := [][]string{
		{"Targets", "Git branch", "main"},
		{"Targets", "Workflow", "test.yml"},
		{"Inputs", "env", "test"},
		{"Inputs", "message", "test"},
		{"Inputs", "server", "app"},
	}

	want := "+---------+------------+----------+\n" +
		"| Targets | Git branch | main     |\n" +
		"+         +------------+----------+\n" +
		"|         | Workflow   | test.yml |\n" +
		"+---------+------------+----------+\n" +
		"| Inputs  | env        | test     |\n" +
		"+         +------------+          +\n" +
		"|         | message    |          |\n" +
		"+         +------------+----------+\n" +
		"|         | server     | app      |\n" +
		"+---------+------------+----------+\n"

	orgStdout := os.Stdout
	pr, pw, _ := os.Pipe()
	os.Stdout = pw

	Render(input)

	pw.Close()
	os.Stdout = orgStdout

	buf := bytes.Buffer{}
	io.Copy(&buf, pr)
	output := buf.String()

	if output != want {
		t.Errorf("Expected is %s but got %s\n", want, output)
	}
}
