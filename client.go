package gorest

// New requires base url for reuse this instance.
// BaseURL includes protocol like `http://` or `https://`
// ex. https://api.github.com/repos/
func New(baseURL string) MethodSelector {
	return &client{
		baseURL:     baseURL,
		contentType: jsonContent,
	}
}

// MethodSelector sets method
// for set param call apigw.Param
// for set path call apigw.Path
type MethodSelector interface {
	Get() TerminalOperator
	Post() TerminalOperator
}

// TerminalOperator executes web api and process result
type TerminalOperator interface {
	// endpoint
	AddPath(path string) TerminalOperator
	AddURLParam(key string, value string) TerminalOperator

	// basic auth
	SetBasicAuth(username string, password string) TerminalOperator

	// header
	AddJSONHeader(key, value string) TerminalOperator

	// body
	SetJSON(json []byte) Executor
	SetJSONString(json string) Executor
	AddURLEncoded(key string, value string) URLEncoded
	AddURLEncodedList(key string, values []string) URLEncoded

	// execute
	Executor
}

type URLEncoded interface {
	AddURLEncoded(key string, value string) URLEncoded
	AddURLEncodedList(key string, values []string) URLEncoded
	Executor
}

type Executor interface {
	Unmarshal(out interface{}) (err error)
	Execute() (result interface{}, err error)
}

type client struct {
	method      RequestMethod
	contentType ContentType
	baseURL     string
	path        []string
	urlParams   []string
	username    *string
	passwd      *string
	headers     map[string]string
	params      interface{}
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
