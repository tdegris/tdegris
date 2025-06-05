package code

import (
	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"honnef.co/go/js/dom/v2"
)

type Source struct {
	code *Code
	div  *dom.HTMLDivElement
}

func newSource(code *Code, parent dom.Element) *Source {
	return &Source{
		code: code,
		div:  code.gui.CreateDIV(parent, ui.Class("code_source_container")),
	}
}

func (s *Source) set(src string) {
	s.div.SetInnerHTML("<pre>" + src + "</pre>")
}
