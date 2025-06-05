//go:build wasm

package code

import (
	"github.com/tdegris/tdegris/internal/pages/wasm/lessons"
	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"honnef.co/go/js/dom/v2"
)

type Code struct {
	gui *ui.UI
	src *Source
	out *Output
}

func New(gui *ui.UI, parent dom.HTMLElement) *Code {
	cd := &Code{gui: gui}
	container := gui.CreateDIV(parent, ui.Class("code_container"))
	cd.src = newSource(cd, container)
	cd.out = newOutput(cd, container)
	return cd
}

func (cd *Code) SetContent(les *lessons.Lesson) {
	cd.src.set(les.Code)
	cd.out.set("OUTPUT")
}
