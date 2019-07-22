package web

import (
	"fmt"
	"strconv"
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

func toString(arg any) string {
	if str, ok := arg.(string); ok {
		return str
	}
	if i, ok := arg.(int); ok {
		return strconv.Itoa(i)
	}
	if b, ok := arg.(bool); ok {
		if b {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%v", arg)
}
