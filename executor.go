package gorest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (cli *client) Execute() (result interface{}, err error) {
	return result, cli.Unmarshal(&result)
}

func (cli *client) Unmarshal(out interface{}) (err error) {
	req, err := cli.buildRequest()
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	resp, err := doRequest(req)
	if err != nil {
		return
	}

	err = decodeBody(resp, &out, nil)
	if err != nil {
		log.Printf(`status is %s`, resp.Status)
	}
	return
}

func (cli *client) buildRequest() (req *http.Request, err error) {
	endpoint := concat(cli.baseURL, strings.Join(cli.path, ``))
	urlParamString := strings.Join(cli.urlParams, `&`)
	if urlParamString != `` {
		endpoint = join(endpoint, urlParamString, `?`)
	}

	body, err := cli.buildParams()
	if err != nil {
		return
	}
	req, err = http.NewRequest(string(cli.method), endpoint, body)
	req.Header.Set(`Content-Type`, string(cli.contentType))

	for key, val := range cli.headers {
		req.Header.Set(key, val)
	}

	if cli.username != nil && cli.passwd != nil {
		req.SetBasicAuth(*cli.username, *cli.passwd)
	}
	fmt.Printf("%+v\n", req)
	return
}

func (cli *client) buildParams() (io.Reader, error) {
	if cli.isParamStruct {
		var err error
		cli.params, err = json.Marshal(cli.params)
		if err != nil {
			return nil, err
		}
	}

	switch cli.contentType {
	case jsonContent:
		jsonBytes, ok := cli.params.([]byte)
		if !ok {
			return nil, nil
		}
		return bytes.NewBuffer(jsonBytes), nil
	case urlEncoded:
		vals, ok := cli.params.(url.Values)
		if !ok {
			return nil, errors.New(`invalid request body parameters`)
		}
		return strings.NewReader(vals.Encode()), nil
	default:
		return nil, errors.New(`unsupported content type for request body`)
	}
}

// Do sends an HTTP request and returns an HTTP response
func doRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

// decodeBody decode response body and stores it in the value pointed to by out
func decodeBody(resp *http.Response, out interface{}, f io.WriteCloser) error {
	defer resp.Body.Close()
	// Symmetric API Testing
	if f != nil {
		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
		defer f.Close()
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
