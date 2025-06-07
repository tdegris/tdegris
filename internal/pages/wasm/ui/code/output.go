package code

import (
	"html"
	"strings"

	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"honnef.co/go/js/dom/v2"
)

type Output struct {
	code *Code
	div  *dom.HTMLDivElement
}

func newOutput(code *Code, parent dom.Element) *Output {
	return &Output{
		code: code,
		div:  code.gui.CreateDIV(parent, ui.Class("code_output_container")),
	}
}

func (s *Output) set(src string) {
	src = html.EscapeString(src)
	src = strings.ReplaceAll(src, "\n", "<br>")
	s.div.SetInnerHTML(src)
}
