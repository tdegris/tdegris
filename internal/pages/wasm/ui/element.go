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

func Property(property, value string) ElementOption {
	return ElementOptionF(func(el dom.Element) {
		el.SetAttribute(property, value)
	})
}

func Listener(typ string, listener func(dom.Event)) ElementOption {
	return ElementOptionF(func(el dom.Element) {
		el.AddEventListener(typ, true, listener)
	})
}

func InnerHTML(s string) ElementOption {
	return ElementOptionF(func(el dom.Element) {
		el.SetInnerHTML(s)
	})
}
