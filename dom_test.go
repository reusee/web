package web

func tempElement(fn func(DOMElement)) {
	elem := Document.Call("createElement", "div")
	elem.Set("id", "app")
	Document.Get("body").Call("appendChild", elem)
	fn(elem)
	Document.Get("body").Call("removeChild", elem)
}
