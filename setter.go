package gorest

import (
	"fmt"
	"net/url"
)

func (cli *client) AddPath(path string) TerminalOperator {
	cli.path = append(cli.path, path)
	return cli
}

func (cli *client) AddURLParam(key string, value string) TerminalOperator {
	cli.urlParams = append(cli.urlParams, fmt.Sprintf("%s=%s", key, value))
	return cli
}

func (cli *client) SetBasicAuth(username string, password string) TerminalOperator {
	cli.username = &username
	cli.passwd = &password
	return cli
}

func (cli *client) SetJSON(json []byte) Executor {
	if len(json) != 0 {
		cli.params = json
	}
	return cli
}

func (cli *client) SetJSONString(json string) Executor {
	if json != `` {
		cli.params = []byte(json)
	}
	return cli
}

func (cli *client) AddURLEncoded(key string, value string) URLEncoded {
	mappedParams, ok := cli.params.(url.Values)
	if !ok {
		mappedParams = url.Values{}
	}
	mappedParams[key] = []string{value}
	cli.params = mappedParams
	return cli
}

func (cli *client) AddURLEncodedList(key string, values []string) URLEncoded{
	mappedParams, ok := cli.params.(url.Values)
	if !ok {
		mappedParams = url.Values{}
	}
	mappedParams[key] = values
	return cli
}

func (cli *client) AddJSONHeader(key, value string) TerminalOperator {
	if len(cli.headers) == 0 {
		cli.headers = map[string]string{}
	}
	cli.contentType = jsonContent
	cli.headers[key] = value
	return cli
}