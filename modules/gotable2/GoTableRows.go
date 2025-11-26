package gotable2

import "strings"

/*
	With gotable2 we work y through x which isn't conventional as normally you follow x through y to
	find the vectors of the table
*/

type Row struct {
	Columns []*Column
}

// Columns are embedded inside the rows of the tables
type Column struct {
	Text 		string
	Align
}

// headerString will build the string representation of the row
func (tablet *GoTable) headerString() []string {
	destination := make([]string, 0)
	if tablet.style.hasAboveHeader() {
		destination = append(destination, *tablet.style.AboveHeaderLeft)
		for index, vector := range tablet.vectors {
			destination[len(destination)-1] += strings.Repeat(*tablet.style.AboveHeaderHorizontal, vector)
			if index + 1 < len(tablet.vectors) {
				destination[len(destination)-1] += *tablet.style.AboveHeaderIntersection
			}
		}

		destination[len(destination) - 1] += *tablet.style.AboveHeaderRight
	}
	
	destination = append(destination, tablet.style.HeaderLeft)
	for index, vector := range tablet.Header.Columns {
		destination[len(destination)-1] += vector.Pad(tablet.vectors[index], vector.Text)
		if index + 1 < len(tablet.vectors) {
			destination[len(destination)-1] += tablet.style.HeaderIntersection
		}
	}
	
	destination[len(destination)-1] += tablet.style.HeaderRight
	if tablet.style.hasBelowHeader() {
		destination = append(destination, *tablet.style.BelowHeaderLeft)
		for index, vector := range tablet.vectors {
			destination[len(destination)-1] += strings.Repeat(*tablet.style.BelowHeaderHorizontal, vector)
			if index + 1 < len(tablet.vectors) {
				destination[len(destination)-1] += *tablet.style.BelowHeaderIntersection
			}
		}

		destination[len(destination) - 1] += *tablet.style.BelowHeaderRight
	}
	
	return destination
}

// values will implement all the individual rows of values presented in the table
func (gotable *GoTable) valueString(destination []string) []string {
	for _, value := range gotable.Rows {
		destination = append(destination, gotable.style.ValueLeft)
		for index, vector := range gotable.vectors {
			destination[len(destination)-1] += value.Columns[index].Pad(vector, value.Columns[index].Text)
			if index + 1 < len(gotable.vectors) {
				destination[len(destination)-1] += gotable.style.ValueIntersection
			}
		}

		destination[len(destination)-1] += gotable.style.ValueRight
	}

	if gotable.style.hasBelowBody() {
		destination = append(destination, *gotable.style.BelowBodyLeft)
		for index, vector := range gotable.vectors {
			destination[len(destination)-1] += strings.Repeat(*gotable.style.BelowBodyHorizontal, vector)
			if index + 1 < len(gotable.vectors) {
				destination[len(destination)-1] += *gotable.style.BelowBodyIntersection
			}
		}

		destination[len(destination) - 1] += *gotable.style.BelowBodyRight
	}

	return destination
}