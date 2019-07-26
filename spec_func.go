package web

import "reflect"

type FuncSpec struct {
	Func any
	Name string
}

var (
	_ Spec    = FuncSpec{}
	_ HasName = &FuncSpec{}
)

func F(fn any) *FuncSpec {
	return &FuncSpec{
		Func: fn,
	}
}

func (f FuncSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement *DOMElement,
	newSpec Spec,
) {
	var spec Spec

	// optimize against observer
	if f.Name == "" {
		f.Name = reflect.TypeOf(f.Func).Name()
	}
	if f.Name != "" {
		observer, ok := oldSpec.(ObserverSpec)
		if ok {
			if f.Name == observer.Name {
				if observer.NoChange(scope) {
					newElement = oldElement
					newSpec = observer
					return
				}
			}
		}
	}

	rets := scope.Call(f.Func, &spec)
	fnType := reflect.TypeOf(f.Func)
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
		t := fnType.Out(i).Elem()
		scope.SetValue(t, ret.Elem())
	}
	newElement, newSpec = spec.Patch(scope, oldSpec, oldElement, replace)

	// wrap to observer if fn is named
	if f.Name != "" {
		newSpec = ObserverSpec{
			Name:         f.Name,
			ScopeVersion: scope.Version,
			Spec:         newSpec,
		}
	}

	return
}

func (f *FuncSpec) GetName() string {
	return f.Name
}

func (f *FuncSpec) SetName(name string) {
	f.Name = name
}

type ObserverSpec struct {
	Name         string
	ScopeVersion int
	Spec         Spec
}

var _ Spec = ObserverSpec{}

func (o ObserverSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement *DOMElement,
	newSpec Spec,
) {
	return o.Spec.Patch(scope, oldSpec, oldElement, replace)
}

func (o ObserverSpec) NoChange(newScope Scope) bool {
	if o.ScopeVersion == newScope.Version {
		return true
	}
	//TODO
	return false
}
