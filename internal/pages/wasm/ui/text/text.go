//go:build wasm

package text

import (
	"github.com/tdegris/tdegris/internal/pages/wasm/lessons"
	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"honnef.co/go/js/dom/v2"
)

type Text struct {
	gui *ui.UI
	div *dom.HTMLDivElement
}

func New(gui *ui.UI, parent dom.HTMLElement) *Text {
	return &Text{
		gui: gui,
		div: gui.CreateDIV(parent, ui.Class("text_container")),
	}
}

func (tt *Text) SetContent(les *lessons.Lesson) {
	tt.div.SetInnerHTML(les.Text)
}
