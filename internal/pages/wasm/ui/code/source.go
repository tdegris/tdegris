package code

import (
	"fmt"
	"strings"

	"github.com/gx-org/gx/build/builder"
	"github.com/gx-org/gx/build/importers"
	"github.com/gx-org/gx/stdlib"
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
		src = append(src, ui.TextContent(child.Underlying()))
	}
	return strings.Join(src, "\n")
}

var keywordToColor = []struct {
	color string
	words []string
}{
	{
		color: "var(--language-keyword)",
		words: []string{
			"var", "const", "return", "struct", "func", "package",
		},
	},
	{
		color: "var(--type-keyword)",
		words: []string{
			"bool", "string",
			"int32", "int64",
			"bfloat64", "float32", "float64",
		},
	},
}

func format(s string) string {
	for _, color := range keywordToColor {
		fontTag := fmt.Sprintf(`<span style="color:%s;">%%s</span>`, color.color)
		for _, word := range color.words {
			s = strings.ReplaceAll(s, word, fmt.Sprintf(fontTag, word))
		}
	}
	return s
}

func (s *Source) set(src string) {
	s.lastSrc = src
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
	sel := s.code.gui.CurrentSelection(s.input)
	defer sel.SetAsCurrent()
	s.set(currentSrc)
}

func (s *Source) onRun(ev dom.Event) {
	bld := builder.New(importers.NewCacheLoader(
		stdlib.Importer(nil),
	))
	pkg := bld.NewIncrementalPackage("main")
	if err := pkg.Build(s.lastSrc); err != nil {
		s.code.out.div.SetInnerHTML(fmt.Sprintf("ERROR: %s", err.Error()))
		return
	}
	irPkg := pkg.IR()
	const fnName = "Main"
	fn := irPkg.FindFunc(fnName)
	if fn == nil {
		s.code.out.div.SetInnerHTML(fmt.Sprintf("ERROR: function %s not found", fnName))
		return
	}
	s.code.out.div.SetInnerHTML(fmt.Sprintf("Main: %v", fn))
}
