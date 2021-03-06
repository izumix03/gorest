package gorest

const (
	get         requestMethod = "GET"
	post                      = "POST"
	put                       = "PUT"
	jsonContent contentType   = "application/json"
	urlEncoded  contentType   = "application/x-www-form-urlencoded"
	multipart   contentType   = "multipart/form-data"
	notSet      contentType   = ""
)

// requestMethod is method name like get, post
type requestMethod string

// contentType is api content type like application/json
type contentType string
