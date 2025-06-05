//go:build wasm

package main

import (
	"fmt"

	"honnef.co/go/js/dom/v2"
)

func findByClass[T dom.Element](start interface{ GetElementsByClassName(string) []dom.Element }, class string) (zero T, err error) {
	els := start.GetElementsByClassName(class)
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

func main() {
	doc := dom.GetWindow().Document()
	body, err := findByClass[dom.HTMLElement](doc, "root")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}
	body.SetInnerHTML("Hello from WASM")
}
