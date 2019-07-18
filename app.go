package web

import (
	"reflect"
	"syscall/js"

	"github.com/reusee/dscope"
)

type App struct {
	Container DOMElement
	Scope     Scope
	SpecFunc  reflect.Value
	Element   DOMElement
	Spec      Spec
}

func NewApp(args ...any) *App {
	app := new(App)
	var fns []interface{}

	for _, arg := range args {
		value := reflect.ValueOf(arg)
		t := value.Type()

		if t.Kind() == reflect.Func && t.NumOut() == 1 && t.Out(0) == specType {
			app.SpecFunc = value

		} else if elem, ok := arg.(DOMElement); ok && elem.InstanceOf(jsElementType) {
			app.Container = elem

		} else if t.Kind() == reflect.Func {
			fns = append(fns, arg)

		} else {
			panic(me(nil, "unknown argument %#v", arg))
		}

	}
	app.Scope = Scope{
		Scope: dscope.New(fns...),
	}

	app.Update()

	return app
}

func (a *App) Update() {
	var spec Spec
	a.Scope.CallValue(a.SpecFunc, &spec)
	a.Element, a.Spec = Patch(spec, a.Element, a.Spec, func(e DOMElement) {
		if a.Element.Type() != js.TypeUndefined {
			a.Container.Call("removeChild", a.Element)
		}
		a.Container.Call("appendChild", e)
	})
}

func (a *App) Close() {
}
