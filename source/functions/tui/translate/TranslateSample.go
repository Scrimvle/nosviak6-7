package translate

import "fmt"

const TERM string = "term"

// NewText will produce the require text pool for the given information
func NewText(label string, x, y int) string {
	return fmt.Sprintf("%s.NewText(\"%s\", %d, %d)", TERM, label, x, y)
}

// NewButton will produce the require button pool for the given information
func NewButton(label string, x, y int) string {
	return fmt.Sprintf("%s.NewButton(%d, %d, \"%s\")", TERM, x, y, label)
}

// NewButtonVar will produce a NewButton and associated variable to the NewButton
func NewButtonVar(label, variable string, x, y int) string {
	return fmt.Sprintf("var %s = %s", variable, NewButton(label, x, y))
}

// NewTermFlowExec will create the event handle
func NewTermFlowExec(label string, content string) string {
	return fmt.Sprintf("func %s() -> bool { %s }", label, content)
}

// NewTermFlowExecSync will connect the event handle with the button logically
func NewTermFlowExecSync(button, funcConnect string) string {
	return fmt.Sprintf("%s.OnClick(%s)", button, funcConnect)
}

// NewInput defines certain params for the input field
func NewInput(label, id string, x, y int) string {
	return fmt.Sprintf("var %s = %s.NewInput(\"%s\", %d, %d)", id, TERM, label, x, y)
}

// Run will initialize the terminal pool
func Run() string {
	return fmt.Sprintf("%s.Run()", TERM)
}