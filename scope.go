package web

import "reflect"

type Scope []map[reflect.Type]reflect.Value

func NewScope(updates ...any) Scope {
	s := Scope([]map[reflect.Type]reflect.Value{
		make(map[reflect.Type]reflect.Value),
	})
	s.Update(updates...)
	return s
}

func (s Scope) Set(t reflect.Type, v any) {
	s.SetValue(t, reflect.ValueOf(v))
}

func (s Scope) SetValue(t reflect.Type, v reflect.Value) {
	for i := len(s) - 1; i >= 0; i-- {
		m := s[i]
		_, ok := m[t]
		if !ok {
			continue
		}
		m[t] = reflect.ValueOf(v)
	}
	panic(me(nil, "%v not in scope", t))
}

func (s Scope) NewValue(t reflect.Type, v reflect.Value) {
	for i := len(s) - 1; i >= 0; i-- {
		m := s[i]
		_, ok := m[t]
		if !ok {
			continue
		}
		m[t] = v
	}
	s[0][t] = v
}

func (s Scope) Get(t reflect.Type) reflect.Value {
	for i := len(s) - 1; i >= 0; i-- {
		m := s[i]
		v, ok := m[t]
		if !ok {
			continue
		}
		return v
	}
	panic(me(nil, "%v not in scope", t))
}

func (s Scope) Assign(targets ...any) {
	for _, target := range targets {
		s.AssignValue(reflect.ValueOf(target))
	}
}

func (s Scope) AssignValue(target reflect.Value) {
	value := reflect.ValueOf(target).Elem()
	t := value.Type()
	value.Set(s.Get(t))
}

func (s Scope) CallValue(fn reflect.Value, targets ...any) []reflect.Value {
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

func (s Scope) Call(fn any, targets ...any) []reflect.Value {
	return s.CallValue(reflect.ValueOf(fn), targets...)
}

func (s Scope) Update(fns ...any) {
	for _, fn := range fns {
		fnValue := reflect.ValueOf(fn)
		rets := s.CallValue(fnValue)
		fnType := fnValue.Type()
		for i := 0; i < fnType.NumOut(); i++ {
			s.NewValue(fnType.Out(i), rets[i])
		}
	}
}
