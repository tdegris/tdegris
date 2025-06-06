package ui

import (
	"fmt"
	"strings"
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

func nodeName(el js.Value) string {
	if el.IsNull() {
		return ""
	}
	return strings.ToUpper(el.Get("nodeName").String())
}

func isDiv(el js.Value) bool {
	return nodeName(el) == "DIV"
}

func findParentDIV(el js.Value) js.Value {
	current := el
	for !isDiv(current) {
		current = current.Get("parentElement")
	}
	return current
}

func lineNumFromElement(el js.Value) int {
	line := 0
	prev := findParentDIV(el).Get("previousElementSibling")
	for !prev.IsNull() {
		prev = prev.Get("previousElementSibling")
		line++
	}
	return line
}

func textLenFromPreviousElement(el js.Value) int {
	if nodeName(el.Get("firstChild")) == "BR" {
		return 0
	}
	// Make sure that the parent is DIV
	// (moving up from text to font in a font tag)
	js.Global().Set("FIRST", el)
	for !isDiv(el.Get("parentNode")) {
		el = el.Get("parentNode")
	}
	if el.Get("parentNode").Get("childNodes").Length() <= 1 {
		return 0
	}
	// Start counting text context in the previous element
	// (ignoring the current element)
	strLen := 0
	js.Global().Set("LAST", el)
	prev := el.Get("previousSibling")
	for !prev.IsNull() {
		fmt.Println("  TEXT", prev.Get("nodeName").String(), prev.Get("textContent").String())
		strLen += len(prev.Get("textContent").String())
		js.Global().Set("LAST", prev)
		prev = prev.Get("previousSibling")
	}
	return strLen
}

func (ui *UI) CurrentSelection(el dom.HTMLElement) *Selection {
	if numRange := selection().Get("rangeCount").Int(); numRange == 0 {
		return nil
	}
	rang := selection().Call("getRangeAt", 0)
	ancestor := rang.Get("commonAncestorContainer")
	line := 0
	if len(el.InnerHTML()) > 1 { // Necessary condition to handle the edge case when there is only a single character.
		line = lineNumFromElement(ancestor)
	}
	fmt.Println("SELECT TEXT", ancestor.Get("textContent").String(), "START", rang.Get("startOffset").Int())
	return &Selection{
		ui:     ui,
		el:     el,
		rang:   rang,
		column: textLenFromPreviousElement(ancestor) + rang.Get("startOffset").Int(),
		line:   line,
	}
}

func findLeafNode(el dom.Node) dom.Node {
	current := el
	for len(current.ChildNodes()) > 0 {
		current = current.FirstChild()
	}
	return current
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
	column := sel.column
	for _, child := range lineDiv.ChildNodes() {
		textLen := len(child.TextContent())
		if column <= textLen {
			selection().Call("collapse", findLeafNode(child).Underlying(), column)
			return
		}
		column -= textLen
	}
}

func (sel *Selection) String() string {
	if sel == nil {
		return "nil"
	}
	return fmt.Sprintf("line: %d col: %d", sel.line, sel.column)
}
