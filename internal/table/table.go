package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// RenderTable renders a table to stdout
func Render(d [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(d)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render()
}
