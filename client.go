package gorest

import (
	"io"
	"net/http"
)

// Get requires base url for reuse this instance.
// BaseURL includes protocol like `http://` or `https://`
// ex. https://api.github.com/repos/
func Get(baseURL string) TerminalOperator {
	return &client{
		baseURL: baseURL,
		method:  get,
	}
}

// Post requires base url for reuse this instance.
// BaseURL includes protocol like `http://` or `https://`
// ex. https://api.github.com/repos/
func Post(baseURL string) TerminalOperator {
	return &client{
		baseURL: baseURL,
		method:  post,
	}
}

func Put(baseURL string) TerminalOperator {
	return &client{
		baseURL: baseURL,
		method:  put,
	}
}

type client struct {
	method            requestMethod
	contentType       contentType
	baseURL           string
	paths             []string
	urlParams         []string
	username          *string
	password          *string
	headers           map[string]string
	params            interface{}
	hasJsonStruct     bool
	multipartSettings []multipartSetting
	handleError       func(*http.Request, *http.Response) (*http.Response, error)
	client            *http.Client
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

	// client
	Client(client *http.Client) TerminalOperator

	// body

	JSON(json []byte) JSONContent
	JSONString(json string) JSONContent
	// JSONStruct requires struct, not pointer.
	// if receive invalid, error occurs when executing.
	JSONStruct(data interface{}) JSONContent
	// URLEncoded add value if same key
	URLEncoded(key string, value string) URLEncoded
	// URLEncodedList replace values if same key
	URLEncodedList(key string, values []string) URLEncoded
	MultipartData(key string, value io.Reader) Multipart
	MultipartAsFormFile(key string, fileName string, reader io.Reader) Multipart
	// HandleResponse require response handler,
	// if create a new response, MUST close old res.Body
	HandleResponse(func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler

	// execute

	Executor
}

type JSONContent interface {
	JSON(json []byte) JSONContent
	JSONString(json string) JSONContent
	// JSONStruct requires struct, not pointer.
	// if receive invalid, error occurs when executing.
	JSONStruct(data interface{}) JSONContent
	Executor
}

// URLEncoded provides methods for sets url encoded body
type URLEncoded interface {
	// URLEncoded add value if same key
	URLEncoded(key string, value string) URLEncoded
	// URLEncodedList replace values if same key
	URLEncodedList(key string, values []string) URLEncoded
	Executor
}

// Multipart provides methods for set multipart data(including file)
type Multipart interface {
	MultipartData(key string, value io.Reader) Multipart
	MultipartAsFormFile(key string, fileName string, reader io.Reader) Multipart
	Executor
}

// Executor provides methods for executing api
type Executor interface {
	Execute() (resp *http.Response, err error)
	HandleBody(f func(body []uint8) error) error
	// HandleResponse require response handler,
	// if create a new response, MUST close old res.Body
	HandleResponse(func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler
}

// ResponseHandler provides wrapper methods with handling error(ex. http status code)
type ResponseHandler interface {
	HandleBody(f func(body []uint8) error) error
}
