package web

import "reflect"

type NamedSpec struct {
	Name string
	Spec Spec
}

var _ Spec = NamedSpec{}

func Named(name string, spec Spec) Spec {
	return NamedSpec{name, spec}
}

type _Name string

var _nameType = reflect.TypeOf((*_Name)(nil)).Elem()

func (n NamedSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement *DOMElement,
	newSpec Spec,
) {
	//TODO use sub scope
	scope.Set(_nameType, _Name(n.Name))
	return n.Spec.Patch(scope, oldSpec, oldElement, replace)
}
