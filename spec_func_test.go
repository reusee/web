package web

import "testing"

func TestObserverNoChange(t *testing.T) {
	tempElement(func(elem DOMElement) {
		numCalled := 0
		app := NewApp(
			elem,
			NewScope(),
			F(func() Spec {
				numCalled++
				return E("div")
			}),
		)
		if numCalled != 1 {
			t.Fatal()
		}
		app.Update()
		if numCalled != 1 {
			t.Fatal()
		}
		app.Update()
		if numCalled != 1 {
			t.Fatal()
		}
	})
}
