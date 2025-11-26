package source

import (
	"Nosviak4/modules/gotable2"
	"encoding/json"
)

// configureTableStyles will promptly parse all the imported styles
func configureTableStyles(jsonBytes []byte) error {
	TABLEVIEWS = make(map[string]*gotable2.Style)
	return json.Unmarshal(jsonBytes, &TABLEVIEWS)
}