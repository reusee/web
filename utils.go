package web

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
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

func trimHTML(src string) string {
	decoder := xml.NewDecoder(strings.NewReader(src))
	buf := new(strings.Builder)
	encoder := xml.NewEncoder(buf)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if elem, ok := token.(xml.StartElement); ok {
			for i := 0; i < len(elem.Attr); i++ {
				attr := elem.Attr[i]
				if strings.HasPrefix(attr.Name.Local, "__") {
					copy(elem.Attr[i:], elem.Attr[i+1:])
					elem.Attr = elem.Attr[:len(elem.Attr)-1]
				}
			}
			token = elem
		}
		ce(encoder.EncodeToken(token))
	}
	ce(encoder.Flush())
	return buf.String()
}

func getHTML(v js.Value) string {
	return trimHTML(v.Get("innerHTML").String())
}
