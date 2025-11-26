package packages

import (
	"Nosviak4/source"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// GLAMOUR is the import path for glamour
const GLAMOUR = "glamour"

// Introducing Glamour, the dynamic text gradient package! With support for background and foreground gradients, it brings your text to life in stunning
// colors. Originally written in Golang and ported to Swash, Glamour offers flexibility and efficiency. Easily create captivating designs with precise
// control over colors, angles, and more. Experience high-performance execution, even with large amounts of text. Elevate your projects with Glamour,
// the ultimate tool for visually striking and elegant text gradients.

// GLAMOURFUNCS represents all the objects ported into the glamour package.
var GLAMOURFUNCS map[string]any = map[string]any{
	"escape":      EscapeCode,
	"foreground":  Foreground,
	"background":  Background,
	"newGradient": NewDerivative,
	"newThemeGradient": NewDerivativeFromTheme,
}

const (
	Foreground int    = 38
	Background int    = 48
	EscapeCode string = "<escape>"
)

type Gradient struct {
	rgb [][3]int
}

// NewDerivative will initialize the gradient object
func NewDerivative() *Gradient {
	return &Gradient{
		rgb: make([][3]int, 0),
	}
}

// NewDerivativeFromTheme will apply a theme to the gradient object
func NewDerivativeFromTheme(themeName string) *Gradient {
	theme, err := source.GetTheme(themeName)
	if err != nil {
		theme, err = source.GetTheme(source.OPTIONS.String("default_theme"))
		if err != nil || theme == nil {
			return NewDerivative()
		}
	}

	gradientMarshaller := NewDerivative()
	for _, colour := range theme.Glamour.Colours {
		gradientMarshaller.AppendRgbToGradient(colour[0], colour[1], colour[2])
	}

	return gradientMarshaller
}

// AppendRgbToGradient will append the colours into the curve
func (G *Gradient) AppendRgbToGradient(red int, green int, blue int) {
	G.rgb = append(G.rgb, [3]int{red, green, blue})
}

// Marshal will create the curve
func (G *Gradient) Marshal(mode int, content string) string {
	return G.Apply(content, make([][3]int, 0), mode, 0)
}

// ResetRGB will remove all the gradients on the colour
func (G *Gradient) ResetRGB() {
	G.rgb = make([][3]int, 0)
}

// curve will perform the linear interpolation on the given colours
func (G *Gradient) Curve(steps int) [][3]int {
	if len(G.rgb) <= 0 {
		return make([][3]int, steps)
	}

	dest := make([][3]int, 0)

	slope, incrementer := (float64(len(G.rgb))-1)/float64(steps-1), 0.0
	for step := 0; step < steps; step++ {
		curve, force := divmod(incrementer, 1)
		if curve >= float64(len(G.rgb)-1) {
			curve = float64(len(G.rgb) - 2)
			force = 1.0
		}

		dest = append(dest, [3]int{int(math.Round(float64(G.rgb[int(curve)][0])*(1-force) + float64(G.rgb[int(curve)+1][0])*force)), int(math.Round(float64(G.rgb[int(curve)][1])*(1-force) + float64(G.rgb[int(curve)+1][1])*force)), int(math.Round(float64(G.rgb[int(curve)][2])*(1-force) + float64(G.rgb[int(curve)+1][2])*force))})
		incrementer += slope
	}

	return dest
}

// apply will attempt to apply the steps to content
func (G *Gradient) Apply(content string, steps [][3]int, mode, x int) string {
	tokens, toggles := Split(content), make([]int, 0)

	lenOf := len(tokens)
	for p := 0; p < len(tokens) && strings.Count(content, EscapeCode) >= 1; p++ {
		if len(tokens[p:]) < len(EscapeCode) || !strings.HasPrefix(strings.Join(tokens[p:][:len(EscapeCode)], ""), EscapeCode) {
			continue
		}

		// implements any escape codes being found within the start
		if escape := strings.Join(tokens[p:][:len(EscapeCode)], "")[len(EscapeCode):]; len(escape) >= 1 {
			tokens[p + len(EscapeCode)] = strings.Join(tokens[p:][:len(EscapeCode)], "")[len(EscapeCode):] + tokens[p + len(EscapeCode)]
		}

		// removes the trace and appends into the toggles array
		toggles = append(toggles, p)
		tokens, lenOf = append(tokens[:p], tokens[p + len(EscapeCode):]...), lenOf - len(EscapeCode)
	}

	if len(steps) == 0 {
		steps = G.Curve(lenOf)
	}

	// buf stores our output and depth decides if its escaped
	buf, depth := bytes.NewBuffer(make([]byte, 0)), 0

	for pos, colour := range steps {
		if pos >= len(tokens) || x > 0 && pos >= x {
			break
		}

		if slices.Contains(toggles, pos) {
			depth++

			if depth % 2 != 0 {
				buf.WriteString(tokens[pos])
				continue
			}
		}

		if depth % 2 != 0 {
			buf.WriteString(tokens[pos])
			continue
		}

		buf.WriteString(fmt.Sprintf("\x1b[" + strconv.Itoa(mode) + ";2;%d;%d;%dm%s", colour[0], colour[1], colour[2], tokens[pos]))
	}

	return buf.String()
}

// divmod is embedded into the marshal function
func divmod(x, y float64) (float64, float64) {
	return math.Floor(x / y), x - y*math.Floor(x/y)
}

// Split will take the string and split the text char by char
func Split(input string) []string {
	slices := strings.Split(input, "")
	
	writeOver := make([]string, 0)

	buf := bytes.NewBuffer(make([]byte, 0))
	if strings.HasPrefix(input, "\x1b") || strings.HasPrefix(input, "\033") {
		buf.WriteString(strings.SplitAfter(input, "m")[0])
	}

	for i := buf.Len(); i < len(slices); i++ {
		writeOver = append(writeOver, buf.String() + slices[i])
		buf.Reset()

		if len(slices) <= i + 1 {
			break
		}

		if slices[i + 1] == "\x1b" || slices[i + 1] == "\033" {
			escape := strings.SplitAfter(strings.Join(slices[i + 1:], ""), "m")[0]
			if i + len(escape) >= len(slices) {
				writeOver[len(writeOver) - 1] += escape
				break
			}

			writeOver[len(writeOver) - 1] += escape
			i += len(escape)
		}

	}

	return writeOver
}
