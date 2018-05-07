package gorest

func concat(a string, b string) string {
	vals := make([]byte, 0, 10)
	vals = append(vals, a...)
	vals = append(vals, b...)
	return string(vals)
}

// join returns joined string(a + separator +  b)
func join(a string, b string, separator string) string {
	vals := make([]byte, 0, 10)
	vals = append(vals, a...)
	vals = append(vals, separator...)
	vals = append(vals, b...)
	return string(vals)
}
