package web

import "testing"

func TestHelloWorld(t *testing.T) {
	elem := Document.Call("createElement", "div")
	elem.Set("id", "app")
	Document.Get("body").Call("appendChild", elem)

	app := NewApp(

		func() (
			greetings string,
		) {
			greetings = "hello, world!"
			return
		},

		func(
			greetings string,
		) Spec {
			return E("div",
				E("span", greetings),
			)
		},

		elem,
	)
	defer app.Close()

	html := elem.Get("innerHTML").String()
	if html != "<div><span>hello, world!</span></div>" {
		t.Fatal()
	}

}
