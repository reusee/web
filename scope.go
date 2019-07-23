package web

import "reflect"

type Scope struct {
	Values  []map[reflect.Type]reflect.Value
	Version int
}

func NewScope(updates ...any) Scope {
	s := Scope{
		Values: []map[reflect.Type]reflect.Value{
			make(map[reflect.Type]reflect.Value),
		},
		Version: 0,
	}
	s.Update(updates...)
	return s
}

func (s *Scope) Set(t reflect.Type, v any) {
	s.SetValue(t, reflect.ValueOf(v))
}

func (s *Scope) SetValue(t reflect.Type, v reflect.Value) {
	for i := len(s.Values) - 1; i >= 0; i-- {
		m := s.Values[i]
		_, ok := m[t]
		if !ok {
			continue
		}
		m[t] = reflect.ValueOf(v)
	}
	s.Values[0][t] = v
	s.Version++
}

func (s *Scope) Get(t reflect.Type) reflect.Value {
	for i := len(s.Values) - 1; i >= 0; i-- {
		m := s.Values[i]
		v, ok := m[t]
		if !ok {
			continue
		}
		return v
	}
	panic(me(nil, "%v not in scope", t))
}

func (s *Scope) Assign(targets ...any) {
	for _, target := range targets {
		s.AssignValue(reflect.ValueOf(target).Elem())
	}
}

func (s *Scope) AssignValue(target reflect.Value) {
	target.Set(s.Get(target.Type()))
}

func (s *Scope) CallValue(fn reflect.Value, targets ...any) []reflect.Value {
	var args []reflect.Value
	fnType := fn.Type()
	for i := 0; i < fnType.NumIn(); i++ {
		args = append(args, s.Get(fnType.In(i)))
	}
	rets := fn.Call(args)
	for _, target := range targets {
		v := reflect.ValueOf(target).Elem()
		for _, ret := range rets {
			if ret.Type() == v.Type() {
				v.Set(ret)
				break
			}
		}
	}
	return rets
}

func (s *Scope) Call(fn any, targets ...any) []reflect.Value {
	return s.CallValue(reflect.ValueOf(fn), targets...)
}

func (s *Scope) Update(fns ...any) {
	for _, fn := range fns {
		fnValue := reflect.ValueOf(fn)
		rets := s.CallValue(fnValue)
		fnType := fnValue.Type()
		for i := 0; i < fnType.NumOut(); i++ {
			s.SetValue(fnType.Out(i), rets[i])
		}
	}
	s.Version++
}
