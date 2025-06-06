package ui

import (
	"fmt"
	"strconv"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

type UI struct {
	win dom.Window
}

func New(win dom.Window) *UI {
	return &UI{win}
}

func (ui *UI) CreateDIV(parent dom.Element, opts ...ElementOption) *dom.HTMLDivElement {
	el := ui.win.Document().CreateElement("div")
	parent.AppendChild(el)
	for _, opt := range opts {
		opt.Apply(el)
	}
	return el.(*dom.HTMLDivElement)
}

func (ui *UI) CreateBR(parent dom.Element, opts ...ElementOption) *dom.HTMLBRElement {
	el := ui.win.Document().CreateElement("br")
	parent.AppendChild(el)
	for _, opt := range opts {
		opt.Apply(el)
	}
	return el.(*dom.HTMLBRElement)
}

type EventFunc func(ev dom.Event)

func (ui *UI) CreateButton(parent dom.Element, text string, f EventFunc, opts ...ElementOption) *dom.HTMLButtonElement {
	el := ui.win.Document().CreateElement("button")
	el.SetTextContent(text)
	parent.AppendChild(el)
	for _, opt := range opts {
		opt.Apply(el)
	}
	el.AddEventListener("click", true, func(ev dom.Event) {
		go f(ev)
	})
	return el.(*dom.HTMLButtonElement)
}

func FindElementByClass[T dom.Element](ui *UI, class string) (zero T, err error) {
	els := ui.win.Document().GetElementsByClassName(class)
	if len(els) == 0 {
		return zero, fmt.Errorf("not element of class %s found", class)
	}
	if len(els) > 1 {
		return zero, fmt.Errorf("too many elements of class %s found", class)
	}
	el := els[0]
	elT, ok := el.(T)
	if !ok {
		return zero, fmt.Errorf("node %s:%T cannot be converted %T", el, el, zero)
	}
	return elT, nil
}

type Selection struct {
	ui     *UI
	el     dom.HTMLElement
	line   int
	column int
	rang   js.Value
}

func selection() js.Value {
	return js.Global().Call("getSelection")
}

func lineFromElement(el js.Value) (int, bool) {
	line := el.Get("dataset").Get("line")
	if line.IsUndefined() {
		return -1, false
	}
	lineI, err := strconv.Atoi(line.String())
	if err != nil {
		fmt.Printf("ERROR: Invalid line number: %q\n", line)
		return -1, false
	}
	return lineI, true
}

func (ui *UI) CurrentSelection(el dom.HTMLElement) *Selection {
	if numRange := selection().Get("rangeCount").Int(); numRange == 0 {
		return nil
	}
	rang := selection().Call("getRangeAt", 0)
	line, lineOk := lineFromElement(rang.Get("commonAncestorContainer").Get("parentElement"))
	if !lineOk {
		line = len(el.ChildNodes()) - 1
	}
	return &Selection{
		ui:     ui,
		el:     el,
		rang:   rang,
		column: rang.Get("startOffset").Int(),
		line:   line,
	}
}

func (sel *Selection) SetAsCurrent() {
	if sel == nil {
		return
	}
	children := sel.el.ChildNodes()
	if sel.line >= len(children) {
		return
	}
	lineDiv := children[sel.line]
	lineNum, lineOk := lineFromElement(lineDiv.Underlying())
	if !lineOk || lineNum != sel.line {
		return
	}
	textLine := lineDiv.FirstChild()
	selection().Call("collapse", textLine.Underlying(), sel.column)
}

func (sel *Selection) String() string {
	if sel == nil {
		return "nil"
	}
	return fmt.Sprintf("line: %d col: %d", sel.line, sel.column)
}
