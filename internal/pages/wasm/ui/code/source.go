package code

import (
	"fmt"
	"strings"

	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"honnef.co/go/js/dom/v2"
)

type Source struct {
	code      *Code
	container *dom.HTMLDivElement
	input     *dom.HTMLDivElement
	control   *dom.HTMLDivElement

	lastSrc string
}

func newSource(code *Code, parent dom.Element) *Source {
	s := &Source{
		code:      code,
		container: code.gui.CreateDIV(parent, ui.Class("code_source_container")),
	}
	s.input = code.gui.CreateDIV(parent,
		ui.Class("code_source_textinput_container"),
		ui.Property("contenteditable", "true"),
		ui.Listener("input", s.onSourceChange),
	)
	s.input.AddEventListener("input", true, func(ev dom.Event) {
	})
	s.control = code.gui.CreateDIV(parent,
		ui.Class("code_source_controls_container"),
	)
	code.gui.CreateButton(s.control, "Run", s.onRun)
	return s
}

func (s *Source) extractSource() string {
	var src []string
	for _, child := range s.input.ChildNodes() {
		src = append(src, child.TextContent())
	}
	return strings.Join(src, "\n")
}

func format(s string) string {
	s = strings.ReplaceAll(s, "le", `<font color="#9900FF">le</font>`)
	return s
}

func (s *Source) set(src string) {
	parent := s.input
	for _, child := range parent.ChildNodes() {
		parent.RemoveChild(child)
	}
	for _, line := range strings.Split(src, "\n") {
		if line == "" {
			line = "<br>"
		} else {
			line = format(line)
		}
		s.code.gui.CreateDIV(parent,
			ui.InnerHTML(line),
		)
	}
}

func (s *Source) onSourceChange(dom.Event) {
	currentSrc := s.extractSource()
	if currentSrc == s.lastSrc {
		return
	}
	fmt.Println("InnerHTML", s.input.InnerHTML())
	sel := s.code.gui.CurrentSelection(s.input)
	defer sel.SetAsCurrent()
	fmt.Println("Selection", sel)
	s.set(currentSrc)
}

func (s *Source) onRun(ev dom.Event) {
	fmt.Println("onRun")
}
