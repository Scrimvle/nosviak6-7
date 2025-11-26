package translate

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/database"
	"encoding/hex"
	"fmt"
)

// Analyze will attempt to analyze the tr
func (t *Translator) Analyze() ([]string, error) {
	if err := t.Parse(); err != nil {
		return make([]string, 0), err
	}

	draw, keepPos := make([]string, 0), make(map[int]int)

	for _, tag := range t.Tags {
		switch tag.Name {

		case "text":
			draw = append(draw, NewText(tag.Content, keepPos[tag.Line], tag.Line))
			keepPos[tag.Line] = keepPos[tag.Line] + gotable2.LenOf(tag.Content)

		case "button":
			kvs := tag.Process()
			if len(kvs) == 0 {
				draw = append(draw, NewButton(tag.Content, keepPos[tag.Line], tag.Line))
				keepPos[tag.Line] = keepPos[tag.Line] + gotable2.LenOf(tag.Content)
				continue
			}


			varContext := "a" + hex.EncodeToString(*database.NewSalt(2))
			if val, ok := kvs["id"]; ok {
				varContext = fmt.Sprint(val)
			}

			draw = append(draw, NewButtonVar(tag.Content, varContext, keepPos[tag.Line], tag.Line))
			if val, ok := kvs["onclick"]; ok {
				funcForButton := varContext + "_onclick"
				draw = append(draw, NewTermFlowExec(funcForButton, fmt.Sprint(val)))
				draw = append(draw, NewTermFlowExecSync(varContext, funcForButton))
			}

			keepPos[tag.Line] = keepPos[tag.Line] + gotable2.LenOf(tag.Content)

		case "input":
			kvs := tag.Process()
			if len(kvs) == 0 {
				continue
			}

			draw = append(draw, NewInput(tag.Content, fmt.Sprint(kvs["id"]), keepPos[tag.Line], tag.Line))
			keepPos[tag.Line] = keepPos[tag.Line] + gotable2.LenOf(tag.Content)

		}
	}

	if len(draw) >= 1 {
		draw = append(draw, Run())
	}

	return append([]string{"<?swash"}, append(draw, "?>")...), nil 
}