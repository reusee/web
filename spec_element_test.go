package web

import "testing"

func TestElementPatchSubs(t *testing.T) {
	tempElement(func(elem DOMElement) {
		app := NewApp(
			func() int {
				return 0
			},
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
			func() int {
				return 0
			},
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
