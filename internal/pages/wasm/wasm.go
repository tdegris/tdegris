//go:build wasm

package main

import (
	"fmt"

	"github.com/tdegris/tdegris/internal/pages/wasm/lessons"
	"github.com/tdegris/tdegris/internal/pages/wasm/ui"
	"github.com/tdegris/tdegris/internal/pages/wasm/ui/code"
	"github.com/tdegris/tdegris/internal/pages/wasm/ui/text"
	"honnef.co/go/js/dom/v2"
)

func main() {
	gui := ui.New(dom.GetWindow())
	body, err := ui.FindElementByClass[dom.HTMLElement](gui, "root_container")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}

	textElement := text.New(gui, body)
	codeElement := code.New(gui, body)

	chapters := lessons.New()
	current := chapters[0].Content[0]
	textElement.SetContent(current)
	codeElement.SetContent(current)
}
