package web

type TextSpec string

var _ Spec = TextSpec("")

func (t TextSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {

	elem := Document.Call("createTextNode", string(t))
	newElement = elem
	newSpec = t
	replace(elem)

	return
}
