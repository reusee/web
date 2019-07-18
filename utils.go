package web

import (
	"fmt"
	"syscall/js"

	"github.com/reusee/e/v2"
)

var (
	pt     = fmt.Printf
	me     = e.Default.WithStack().WithName("web")
	ce, he = e.New(me)
)

type (
	any     = interface{}
	JSValue = js.Value
)
