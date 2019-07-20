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
		) (
			spec Spec,
			nextGreetings *string,
		) {
			spec = E("div",
				E("span", greetings),
			)
			greetings += " again"
			nextGreetings = &greetings
			return
		},

		elem,
	)

	html := elem.Get("innerHTML").String()
	if html != "<div><span>hello, world!</span></div>" {
		t.Fatal()
	}

	var greetings string
	app.Scope.Assign(&greetings)
	if greetings != "hello, world! again" {
		t.Fatal()
	}

	app.Update()
	html = elem.Get("innerHTML").String()
	if html != "<div><span>hello, world! again</span></div>" {
		t.Fatal()
	}

}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		elem := Document.Call("createElement", "div")
		elem.Set("id", "app")
		Document.Get("body").Call("appendChild", elem)
		NewApp(
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
		Document.Get("body").Call("removeChild", elem)
	}
}
