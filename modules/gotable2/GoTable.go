package gotable2

import (
	"strings"
)

/*
	GoTable2 is an implementation of ANSI tables which can be written to the terminal, this package
	directly implements some basic themes and functionality while keep it's self contained and
	portable from project to project.
*/

type GoTable struct {
	// vectors for the table follow a simple order and that's each vector inside this slice
	// represents a column, the value inside that slice represents the widest point in that
	// column.
	vectors []int
	
	// Header represents the data guidelines for the table
	Header 	*Row

	// Rows represents the data for the table
	Rows 	[]*Row

	// style is the style of the table for the borders etc
	style 	*Style
	
	// BorderValues adds padding onto each row
	BorderValues int

	// LongestLine implements the longest line inside the terminal
	LongestLine int
}

// NewGoTable creates a new GoTable interface
func NewGoTable(style *Style) *GoTable {
	if style == nil {
		style = DEFAULT

	}
	
	return &GoTable{BorderValues: style.BorderValues, Rows: make([]*Row, 0), Header: nil, style: style, vectors: make([]int, 0)}
}

// SetStyle sets the style of the table
func (GoTable *GoTable) SetStyle(style *Style) {
	GoTable.style = style
}

// Head will directly insert into the method and modify the vectors of the table directly.
func (gotable *GoTable) Head(row *Row) {
	for _, column := range row.Columns {
		column.Text = strings.Repeat(" ", gotable.BorderValues) + column.Text + strings.Repeat(" ", gotable.BorderValues)
		gotable.vectors = append(gotable.vectors, LenOf(column.Text))
	}

	gotable.Header = row
}

// Append will indirectly insert the row into the table
func (gotable *GoTable) Append(row *Row) {
	for i, column := range row.Columns {
		column.Text = strings.Repeat(" ", gotable.BorderValues) + column.Text + strings.Repeat(" ", gotable.BorderValues)

		index := LenOf(column.Text)
		if i >= len(gotable.vectors) {
			break
		}
		
		if gotable.vectors[i] < index {
			gotable.vectors[i] = index
		}
	}

	gotable.Rows = append(gotable.Rows, row)
}