package python3

// TypeName returns the name of the type of the given object.
func TypeName(o *Object) string {
	return o.Type().GetAttr("__name__").String()
}
