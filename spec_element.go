package web

import (
	"reflect"
	"syscall/js"
)

type ElementSpec struct {
	Tag      string
	Children []Spec
	//TODO attrs
	//TODO style
	//TODO events
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

func (e ElementSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {

	notPatchable := false
	if oldElement == nil {
		notPatchable = true
	}
	e2, ok := oldSpec.(ElementSpec)
	if !ok {
		notPatchable = true
	}
	if e2.Tag != e.Tag {
		notPatchable = true
	}

	if notPatchable {
		newElement = Document.Call("createElement", e.Tag)
		for i, child := range e.Children {
			_, newChildSpec := child.Patch(
				scope,
				nil,
				nil,
				func(elem DOMElement) {
					newElement.Call("appendChild", elem)
				},
			)
			e.Children[i] = newChildSpec
		}
		replace(newElement)
		newSpec = e
		return
	}

	// patch children
	for i, child := range e.Children {
		oldChildElement := oldElement.Get("childNodes").Index(i)
		var oldChildElementArg *DOMElement
		if oldChildElement.Type() != js.TypeUndefined {
			oldChildElementArg = &oldChildElement
		}
		_, newChildSpec := child.Patch(
			scope,
			e2.Children[i],
			oldChildElementArg,
			func(newChild DOMElement) {
				oldElement.Call("replaceChild", newChild, oldChildElement)
			},
		)
		e.Children[i] = newChildSpec
	}

	newElement = *oldElement
	newSpec = e
	return
}
