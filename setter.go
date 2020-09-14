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

func (cli *client) Header(key, value string) TerminalOperator {
	if len(cli.headers) == 0 {
		cli.headers = map[string]string{}
	}
	cli.headers[key] = value
	return cli
}

func (cli *client) Client(client *http.Client) TerminalOperator {
	cli.client = client
	return cli
}

func (cli *client) JSON(json []byte) JSONContent {
	if len(json) != 0 {
		cli.params = json
	}
	cli.contentType = jsonContent
	return cli
}

func (cli *client) JSONString(json string) JSONContent {
	if json != `` {
		cli.params = []byte(json)
	}
	cli.contentType = jsonContent
	return cli
}

func (cli *client) JSONStruct(data interface{}) JSONContent {
	cli.params = data
	cli.contentType = jsonContent
	cli.hasJsonStruct = true
	return cli
}

func (cli *client) URLEncoded(key string, value string) URLEncoded {
	urlValues, ok := cli.params.(url.Values)
	if !ok {
		urlValues = url.Values{}
		cli.params = urlValues
	}
	urlValues.Add(key, value)
	cli.contentType = urlEncoded
	return cli
}

func (cli *client) URLEncodedList(key string, values []string) URLEncoded {
	urlValues, ok := cli.params.(url.Values)
	if !ok {
		urlValues = url.Values{}
		cli.params = urlValues
	}
	urlValues[key] = values
	cli.contentType = urlEncoded
	return cli
}

func (cli *client) URLEncodedString(data string) URLEncoded {
	cli.params = data
	cli.contentType = urlEncoded
	cli.hasRawFormUrlEncoded = true
	return cli
}

func (cli *client) HandleResponse(f func(*http.Request, *http.Response) (*http.Response, error)) ResponseHandler {
	cli.responseHandler = f
	return cli
}
