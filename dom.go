package web

import (
	"syscall/js"
)

type (
	DOMElement = JSValue
)

var (
	Global   = js.Global()
	Document = Global.Get("document")
)

var (
	jsElementType = Global.Call("eval", "Element")
)
