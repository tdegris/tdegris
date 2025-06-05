package ui

import "honnef.co/go/js/dom/v2"

type (
	ElementOption interface {
		Apply(dom.Element)
	}

	ElementOptionF func(dom.Element)
)

func (f ElementOptionF) Apply(el dom.Element) {
	f(el)
}

func Class(class string) ElementOption {
	return ElementOptionF(func(el dom.Element) {
		el.SetAttribute("class", class)
	})
}
