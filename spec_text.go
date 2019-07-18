package web

type TextSpec string

var _ Spec = TextSpec("")

func (t TextSpec) Identical(spec Spec) bool {
	t2, ok := spec.(TextSpec)
	if !ok {
		return false
	}
	return t == t2
}

func (t TextSpec) MakeElement() DOMElement {
	return Document.Call("createTextNode", string(t))
}

func (t TextSpec) Patchable(_ Spec) bool {
	return false
}
