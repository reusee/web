package web

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"

	"github.com/reusee/e/v2"
)

func init() {
	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
}

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
