package web

import "reflect"

type Spec interface {
	Patch(
		scope Scope,
		oldSpec Spec,
		oldElement *DOMElement,
		replace func(DOMElement),
	) (
		newElement DOMElement,
		newSpec Spec,
	)
}

var specType = reflect.TypeOf((*Spec)(nil)).Elem()
