package gorest

import "net/http"

// Post requires base url for reuse this instance.
// BaseURL includes protocol like `http://` or `https://`
// ex. https://api.github.com/repos/
func Post(baseURL string) TerminalOperator {
	return &client{
		baseURL:     baseURL,
		contentType: jsonContent,
		method:      post,
	}
}

func Put(baseURL string) TerminalOperator {
	return &client{
		baseURL:     baseURL,
		contentType: jsonContent,
		method:      put,
	}
}

func Get(baseURL string) TerminalOperator {
	return &client{
		baseURL:     baseURL,
		contentType: jsonContent,
		method:      get,
	}
}

type client struct {
	method        requestMethod
	contentType   contentType
	baseURL       string
	paths         []string
	urlParams     []string
	username      *string
	password      *string
	headers       map[string]string
	params        interface{}
	isParamStruct bool
	handleError   func(*http.Request, *http.Response) (*http.Response, error)
}

// TerminalOperator executes web api and process result
type TerminalOperator interface {
	// endpoint
	Path(pathFmt string, args ...interface{}) TerminalOperator
	URLParam(key string, value string) TerminalOperator

	// basic auth
	BasicAuth(username string, password string) TerminalOperator

	// header
	Header(key, value string) TerminalOperator

	// body
	JSON(json []byte) Executor
	JSONString(json string) Executor
	JSONStruct(data interface{}) Executor

	URLEncoded(key string, value string) URLEncoded
	URLEncodedList(key string, values []string) URLEncoded

	// if create new response, MUST close old res.Body
	HandleResponse(func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler

	// execute
	Executor
}

// URLEncoded provides methods for sets url encoded body
type URLEncoded interface {
	URLEncoded(key string, value string) URLEncoded
	URLEncodedList(key string, values []string) URLEncoded
	Executor
}

// Executor provides methods for executing api
type Executor interface {
	Execute() (resp *http.Response, err error)
	HandleBody(f func(body []uint8) error) error
	// if create new response, MUST close old res.Body
	HandleResponse(func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler
}

// ResponseHandler provides wrapper methods with handling error(ex. http status code)
type ResponseHandler interface {
	HandleBody(f func(body []uint8) error) error
}
