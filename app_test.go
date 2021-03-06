package web

import "testing"

func TestHelloWorld(t *testing.T) {
	tempElement(func(elem DOMElement) {
		app := NewApp(

			NewScope(func() (
				greetings string,
			) {
				greetings = "hello, world!"
				return
			}),

			F(func(
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
			}),

			elem,
		)

		html := getHTML(elem)
		if html != "<div><span>hello, world!</span></div>" {
			t.Fatal(html)
		}

		var greetings string
		app.Scope.Assign(&greetings)
		if greetings != "hello, world! again" {
			t.Fatal()
		}

		app.Update()
		html = getHTML(elem)
		if html != "<div><span>hello, world! again</span></div>" {
			t.Fatal()
		}

	})
}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempElement(func(elem DOMElement) {
			NewApp(
				NewScope(func() (
					greetings string,
				) {
					greetings = "hello, world!"
					return
				}),
				F(func(
					greetings string,
				) Spec {
					return E("div",
						E("span", greetings),
					)
				}),
				elem,
			)
		})
	}
}
