package web

import "reflect"

type FuncSpec struct {
	Func any
}

var _ Spec = FuncSpec{}

func F(fn any) FuncSpec {
	return FuncSpec{fn}
}

func (f FuncSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {
	var spec Spec
	rets := scope.Call(f.Func, &spec)
	specConstructorType := reflect.TypeOf(f.Func)
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
	return spec.Patch(scope, oldSpec, oldElement, replace)
}
