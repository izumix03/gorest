package gorest

func concat(a string, b string) string {
	results := make([]byte, 0, 10)
	results = append(results, a...)
	results = append(results, b...)
	return string(results)
}

// join returns joined string(a + separator +  b)
func join(a string, b string, separator string) string {
	results := make([]byte, 0, 10)
	results = append(results, a...)
	results = append(results, separator...)
	results = append(results, b...)
	return string(results)
}
