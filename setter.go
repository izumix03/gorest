package gorest

import (
	"fmt"
	"net/http"
	"net/url"
)

func (cli *client) Path(pathFmt string, args ...interface{}) TerminalOperator {
	cli.paths = append(cli.paths, fmt.Sprintf(pathFmt, args...))
	return cli
}

func (cli *client) URLParam(key string, value string) TerminalOperator {
	cli.urlParams = append(cli.urlParams, fmt.Sprintf("%s=%s", key, value))
	return cli
}

func (cli *client) BasicAuth(username string, password string) TerminalOperator {
	cli.username = &username
	cli.password = &password
	return cli
}

func (cli *client) JSON(json []byte) Executor {
	if len(json) != 0 {
		cli.params = json
	}
	cli.contentType = jsonContent
	return cli
}

func (cli *client) JSONString(json string) Executor {
	if json != `` {
		cli.params = []byte(json)
	}
	cli.contentType = jsonContent
	return cli
}

func (cli *client) JSONStruct(data interface{}) Executor {
	cli.params = data
	cli.contentType = jsonContent
	cli.isParamStruct = true
	return cli
}

func (cli *client) URLEncoded(key string, value string) URLEncoded {
	mappedParams, ok := cli.params.(url.Values)
	if !ok {
		mappedParams = url.Values{}
	}
	mappedParams[key] = []string{value}
	cli.params = mappedParams
	cli.contentType = urlEncoded
	return cli
}

func (cli *client) URLEncodedList(key string, values []string) URLEncoded {
	mappedParams, ok := cli.params.(url.Values)
	if !ok {
		mappedParams = url.Values{}
	}
	mappedParams[key] = values
	cli.contentType = urlEncoded
	return cli
}

func (cli *client) Header(key, value string) TerminalOperator {
	if len(cli.headers) == 0 {
		cli.headers = map[string]string{}
	}
	cli.headers[key] = value
	return cli
}

func (cli *client) HandleResponse(f func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler {
	cli.handleError = f
	return cli
}
