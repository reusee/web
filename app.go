package web

type App struct {
	Container DOMElement
	Scope     Scope
	RootSpec  Spec

	Element *DOMElement
	Spec    Spec
}

func NewApp(args ...any) *App {
	app := &App{}

	for _, arg := range args {

		if spec, ok := arg.(Spec); ok {
			app.RootSpec = spec

		} else if elem, ok := arg.(DOMElement); ok && elem.InstanceOf(jsElementType) {
			app.Container = elem

		} else if scope, ok := arg.(Scope); ok {
			app.Scope = scope

		} else {
			panic(me(nil, "unknown argument %#v", arg))
		}

	}

	app.Update()

	return app
}

func (a *App) Update() {
	a.Element, a.Spec = a.RootSpec.Patch(
		a.Scope,
		a.Spec,
		a.Element,
		func(e DOMElement) {
			if a.Element != nil {
				a.Container.Call("removeChild", a.Element)
			}
			a.Container.Call("appendChild", e)
		},
	)
}
