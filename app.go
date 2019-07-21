package web

import (
	"reflect"
	"syscall/js"
)

type App struct {
	Container DOMElement
	Scope     Scope
	RootSpec  Spec

	Element DOMElement
	Spec    Spec
}

func NewApp(args ...any) *App {
	app := &App{}
	var fns []interface{}

	for _, arg := range args {
		value := reflect.ValueOf(arg)
		t := value.Type()

		if t.Kind() == reflect.Func && func() bool {
			for i := 0; i < t.NumOut(); i++ {
				if t.Out(i) == specType {
					return true
				}
			}
			return false
		}() {
			app.RootSpec = F(arg)

		} else if elem, ok := arg.(DOMElement); ok && elem.InstanceOf(jsElementType) {
			app.Container = elem

		} else if t.Kind() == reflect.Func {
			fns = append(fns, arg)

		} else {
			panic(me(nil, "unknown argument %#v", arg))
		}

	}
	app.Scope = NewScope(fns...)

	app.Update()

	return app
}

func (a *App) Update() {
	a.Element, a.Spec = a.RootSpec.Patch(
		a.Scope,
		a.Spec,
		&a.Element,
		func(e DOMElement) {
			if a.Element.Type() != js.TypeUndefined {
				a.Container.Call("removeChild", a.Element)
			}
			a.Container.Call("appendChild", e)
		},
	)
}
