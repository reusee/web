package web

import "testing"

func TestElementPatchSubs(t *testing.T) {
	tempElement(func(elem DOMElement) {
		app := NewApp(
			NewScope(func() int {
				return 0
			}),
			F(func(
				i int,
			) (
				spec Spec,
				nextI *int,
			) {
				if i == 0 {
					spec = E("span", "foo")
				} else if i == 1 {
					spec = E("div", E("span", "bar"))
				} else {
					spec = E("span", "ok")
				}
				i++
				nextI = &i
				return
			}),
			elem,
		)

		html := elem.Get("innerHTML").String()
		if html != "<span>foo</span>" {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != "<div><span>bar</span></div>" {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != "<span>ok</span>" {
			t.Fatal()
		}

	})
}

func TestElementPatchNotPatchable(t *testing.T) {
	tempElement(func(elem DOMElement) {
		app := NewApp(
			NewScope(func() int {
				return 0
			}),
			F(func(
				i int,
			) (
				spec Spec,
				nextI *int,
			) {
				if i == 0 {
					spec = E("span", "foo")
				} else if i == 1 {
					spec = TextSpec("bar")
				} else if i == 2 {
					spec = E("div", "yes")
				} else {
					spec = E("span", "ok")
				}
				i++
				nextI = &i
				return
			}),
			elem,
		)

		html := elem.Get("innerHTML").String()
		if html != "<span>foo</span>" {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != "bar" {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != "<div>yes</div>" {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != "<span>ok</span>" {
			t.Fatal()
		}

	})
}

func TestElementAttrs(t *testing.T) {
	tempElement(func(elem DOMElement) {
		app := NewApp(
			NewScope(func() (
				string,
				int,
			) {
				return "hello", 0
			}),
			F(func(
				s string,
				n int,
			) (
				spec Spec,
				nextN *int,
			) {
				if n == 0 {
					spec = E("div",
						A("foo", s),
					)
				} else {
					spec = E("div",
						A("bar", s),
					)
				}
				n++
				nextN = &n
				return
			}),
			elem,
		)

		html := elem.Get("innerHTML").String()
		if html != `<div foo="hello"></div>` {
			t.Fatal()
		}

		app.Update()
		html = elem.Get("innerHTML").String()
		if html != `<div bar="hello"></div>` {
			t.Fatal()
		}
	})
}

func BenchmarkDeepNestedElement(b *testing.B) {
	e := E("div")
	for i := 0; i < 512; i++ {
		e = E("div", e)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tempElement(func(elem DOMElement) {
			NewApp(
				e,
				elem,
			)
		})
	}
}

func BenchmarkManySubElements(b *testing.B) {
	var subs []any
	for i := 0; i < 512; i++ {
		subs = append(subs, E("div"))
	}
	e := E("div", subs...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tempElement(func(elem DOMElement) {
			NewApp(
				e,
				elem,
			)
		})
	}
}
