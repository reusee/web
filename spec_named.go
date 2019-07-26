package web

type _NamedSpec struct {
	Spec
}

var _ Spec = _NamedSpec{}

type HasName interface {
	GetName() string
	SetName(string)
}

func Named(name string, spec Spec) Spec {
	if obj, ok := spec.(HasName); ok {
		obj.SetName(name)
	} else {
		panic(me(nil, "%#v has no name", spec))
	}
	return _NamedSpec{spec}
}
