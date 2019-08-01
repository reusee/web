package web

import (
	"reflect"
	"sync"
	"sync/atomic"
	"syscall/js"
)

type ElementSpec struct {
	Tag      string
	Attrs    []ElementAttr
	Children []Spec
	//TODO attrs
	//TODO style
	//TODO events
}

var _ Spec = ElementSpec{}

func E(tag string, args ...any) Spec {
	elem := ElementSpec{
		Tag: tag,
	}
	for _, arg := range args {
		if s, ok := arg.(Spec); ok {
			elem.Children = append(elem.Children, s)
		} else if attr, ok := arg.(ElementAttr); ok {
			elem.Attrs = append(elem.Attrs, attr)
		} else if attrs, ok := arg.([]ElementAttr); ok {
			elem.Attrs = append(elem.Attrs, attrs...)
		} else {

			value := reflect.ValueOf(arg)
			switch value.Kind() {
			case reflect.String:
				elem.Children = append(elem.Children, TextSpec(value.String()))
			default:
				panic(me(nil, "unknown argument %#v", arg))
			}

		}
	}
	return elem
}

type ElementAttr [2]string

func A(args ...any) (ret []ElementAttr) {
	for i := 0; i < len(args); i += 2 {
		ret = append(ret, ElementAttr{
			toString(args[i]),
			toString(args[i+1]),
		})
	}
	return
}

var nextElementID int32

func (e ElementSpec) Patch(
	scope Scope,
	oldSpec Spec,
	oldElement *DOMElement,
	replace func(DOMElement),
) (
	newElement *DOMElement,
	newSpec Spec,
) {

	defer func() {
		// recycle element
		if oldElement == nil || newElement == oldElement {
			return
		}
		elementRecycleChan <- oldElement
	}()

	notPatchable := false
	if oldElement == nil {
		notPatchable = true
	}
	e2, ok := oldSpec.(ElementSpec)
	if !ok {
		notPatchable = true
	}
	if e2.Tag != e.Tag {
		notPatchable = true
	}

	if notPatchable {
		elem := Document.Call("createElement", e.Tag)
		id := atomic.AddInt32(&nextElementID, 1)
		elem.Set("__id__", id)
		for _, kv := range e.Attrs {
			elem.Call("setAttribute", kv[0], kv[1])
		}
		for i, child := range e.Children {
			_, newChildSpec := child.Patch(
				scope,
				nil,
				nil,
				func(newChild DOMElement) {
					elem.Call("appendChild", newChild)
				},
			)
			e.Children[i] = newChildSpec
		}
		replace(elem)
		newElement = &elem
		newSpec = e
		return
	}

	// patch attrs
	attrNames := make(map[string]bool)
	for _, kv := range e.Attrs {
		oldElement.Call("setAttribute", kv[0], kv[1])
		attrNames[kv[0]] = true
	}
	for _, kv := range e2.Attrs {
		if !attrNames[kv[0]] {
			oldElement.Call("removeAttribute", kv[0])
		}
	}

	// patch children
	for i, child := range e.Children {
		oldChildElement := oldElement.Get("childNodes").Index(i)
		var oldChildElementArg *DOMElement
		if oldChildElement.Type() != js.TypeUndefined {
			oldChildElementArg = &oldChildElement
		}
		_, newChildSpec := child.Patch(
			scope,
			e2.Children[i],
			oldChildElementArg,
			func(newChild DOMElement) {
				oldElement.Call("replaceChild", newChild, oldChildElement)
			},
		)
		e.Children[i] = newChildSpec
	}

	newElement = oldElement
	newSpec = e
	return
}

// element recycler

var elementRecycleChan = make(chan *js.Value, 1024)

var elemFinalizers = make(map[int32][]func())

var elemFinalizersL sync.RWMutex

func setFinalizer(id int32, fn func()) {
	elemFinalizersL.Lock()
	elemFinalizers[id] = append(elemFinalizers[id], fn)
	elemFinalizersL.Unlock()
}

func init() {
	go func() {
		for {
			select {

			case elem := <-elementRecycleChan:
				idValue := elem.Get("__id__")
				if idValue.Type() != js.TypeUndefined {
					id := int32(idValue.Int())
					elemFinalizersL.RLock()
					for _, fn := range elemFinalizers[id] {
						fn()
					}
					elemFinalizersL.RUnlock()
				}

			}
		}
	}()
}
