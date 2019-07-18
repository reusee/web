package web

import "reflect"

type ElementSpec struct {
	Tag      string
	Children []Spec
}

var _ Spec = ElementSpec{}

func E(tag string, args ...any) Spec {
	elem := ElementSpec{
		Tag: tag,
	}
	for _, arg := range args {
		if s, ok := arg.(Spec); ok {
			elem.Children = append(elem.Children, s)
		} else {

			value := reflect.ValueOf(arg)
			switch value.Kind() {
			case reflect.String:
				elem.Children = append(elem.Children, TextSpec(value.String()))
			default:
				panic(me(nil, "unknown argument %#v", arg))
			}

		}
	}
	return elem
}

func (e ElementSpec) Identical(spec Spec) bool {
	e2, ok := spec.(ElementSpec)
	if !ok {
		return false
	}
	if e.Tag != e2.Tag {
		return false
	}
	if len(e.Children) != len(e2.Children) {
		return false
	}
	for i, c := range e.Children {
		if !c.Identical(e2.Children[i]) {
			return false
		}
	}
	return true
}

func (e ElementSpec) Patchable(spec Spec) bool {
	e2, ok := spec.(ElementSpec)
	if !ok {
		return false
	}
	if e.Tag != e2.Tag {
		// can't change element tag
		return false
	}
	return true
}

func (e ElementSpec) MakeElement() DOMElement {
	domElement := Document.Call("createElement", e.Tag)
	for _, child := range e.Children {
		childElement := child.MakeElement()
		domElement.Call("appendChild", childElement)
	}
	return domElement
}
