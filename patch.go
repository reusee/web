package web

func Patch(
	spec Spec,
	oldElement DOMElement,
	oldSpec Spec,
	replace func(DOMElement),
) (
	newElement DOMElement,
	newSpec Spec,
) {

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
