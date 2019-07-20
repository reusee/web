package web

import "reflect"

type Spec interface {
	Identical(spec Spec) bool
	Patchable(spec Spec) bool
	MakeElement() DOMElement
}

var specType = reflect.TypeOf((*Spec)(nil)).Elem()

type SpecConstructor any
