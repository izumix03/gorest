package gorest

// Get requires base url for reuse this instance.
// BaseURL includes protocol like `http://` or `https://`
// ex. https://api.github.com/repos/
func Get(baseURL string) TerminalOperator {
	return &client{
		baseURL:     baseURL,
		contentType: jsonContent,
		method:      get,
	}
}

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

// TerminalOperator executes web api and process result
type TerminalOperator interface {
	// endpoint
	Path(path string) TerminalOperator
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

	// execute
	Executor
}

type URLEncoded interface {
	URLEncoded(key string, value string) URLEncoded
	URLEncodedList(key string, values []string) URLEncoded
	Executor
}

type Executor interface {
	Unmarshal(out interface{}) (err error)
	Execute() (result interface{}, err error)
}

type client struct {
	method        RequestMethod
	contentType   ContentType
	baseURL       string
	path          []string
	urlParams     []string
	username      *string
	passwd        *string
	headers       map[string]string
	params        interface{}
	isParamStruct bool
}

const (
	get         RequestMethod = `GET`
	post                      = `POST`
	jsonContent ContentType   = `application/json`
	urlEncoded  ContentType   = `application/x-www-form-urlencoded`
)

// RequestMethod is method name like get, post
type RequestMethod string

// ContentType is api content type like application/json
type ContentType string
