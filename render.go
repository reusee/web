package web

func Render(
	specConstructor SpecConstructor,
	scope Scope,
	oldSpec Spec,
	oldElement DOMElement,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {

	var spec Spec
	scope.Call(specConstructor, &spec)

	if spec.Identical(oldSpec) {
		newElement = oldElement
		newSpec = newSpec
		return
	}

	if !spec.Patchable(oldSpec) {
		newElement = spec.MakeElement()
		newSpec = spec
		replace(newElement)
		return
	}

	//TODO patch

	return
}
