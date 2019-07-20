package web

import "reflect"

func Render(
	specConstructor SpecConstructor,
	scope Scope,
	oldSpec Spec,
	oldElement DOMElement,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {

	var spec Spec
	rets := scope.Call(specConstructor, &spec)
	specConstructorType := reflect.TypeOf(specConstructor)
	for i, ret := range rets {
		if ret.Type() == specType {
			continue
		}
		if ret.Kind() != reflect.Ptr {
			panic(me(nil, "bad return type, must be pointer: %v", ret.Type()))
		}
		if ret.IsNil() {
			continue
		}
		t := specConstructorType.Out(i).Elem()
		scope.SetValue(t, ret.Elem())
	}

	if spec.Identical(oldSpec) {
		newElement = oldElement
		newSpec = newSpec
		return
	}

	if !spec.Patchable(oldSpec) {
		newElement = spec.MakeElement()
		newSpec = spec
		replace(newElement)
		return
	}

	//TODO patch

	return
}
