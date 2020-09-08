package gorest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Execute executes api and return result, error
func (cli *client) Execute() (*http.Response, error) {
	req, err := cli.buildRequest()
	if err != nil {
		return nil, err
	}

	return doRequest(req)
}

func (cli *client) HandleBody(f func(body []uint8) error) error {
	req, err := cli.buildRequest()
	if err != nil {
		return err
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}
	if cli.handleError != nil {
		res, err = func() (*http.Response, error) {
			return cli.handleError(req, res)
		}()
		if err != nil {
			return err
		}
	}
	defer CloseBody(res.Body)

	if err := cli.handleByStatusCode(res); err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return f(body)
}

func CloseBody(body io.ReadCloser) {
	_, _ = io.Copy(ioutil.Discard, body)
	_ = body.Close()
}

func (cli *client) handleByStatusCode(res *http.Response) error {
	if res.StatusCode >= 400 {
		var responseBody []byte
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			responseBody = body
		}
		return &InvalidStatusCodeError{
			StatusCode:   res.StatusCode,
			ResponseBody: responseBody,
		}
	}
	return nil
}

type InvalidStatusCodeError struct {
	StatusCode   int
	ResponseBody []byte
}

func (i *InvalidStatusCodeError) Error() string {
	return fmt.Sprintf("StatusCode: %d, responseBody: %v", i.StatusCode, string(i.ResponseBody))
}

func (cli *client) buildRequest() (*http.Request, error) {
	endpoint := concat(cli.baseURL, strings.Join(cli.paths, ``))
	urlParamString := strings.Join(cli.urlParams, `&`)
	if urlParamString != `` {
		endpoint = join(endpoint, urlParamString, `?`)
	}

	body, err := cli.buildParams()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(string(cli.method), endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(`Content-Type`, string(cli.contentType))
	for key, val := range cli.headers {
		req.Header.Set(key, val)
	}

	if cli.username != nil && cli.password != nil {
		req.SetBasicAuth(*cli.username, *cli.password)
	}

	return req, nil
}

func (cli *client) buildParams() (io.Reader, error) {
	if cli.hasJsonStruct {
		var err error
		cli.params, err = json.Marshal(cli.params)
		if err != nil {
			return nil, err
		}
	}

	switch cli.contentType {
	case jsonContent:
		if cli.params == nil {
			return nil, nil
		}
		jsonBytes, ok := cli.params.([]byte)
		if !ok {
			// maybe JSONStruct receive invalid data
			return nil, errors.New("invalid body")
		}
		return bytes.NewBuffer(jsonBytes), nil
	case urlEncoded:
		if cli.params == nil {
			return nil, nil
		}
		values, ok := cli.params.(url.Values)
		if !ok {
			return nil, errors.New(`invalid request body parameters`)
		}
		return strings.NewReader(values.Encode()), nil
	case notSet:
		if len(cli.multipartSettings) == 0 {
			return nil, nil
		}
		body, err := cli.setupMultipartRequest()
		if err != nil {
			return nil, fmt.Errorf(`invalid request body parameters, %s`, err)
		}
		return body, nil
	default:
		return nil, errors.New(`unsupported content type for request body`)
	}
}

// Do sends an HTTP request and returns an HTTP response
func doRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

//
//// decodeBody decode response body and stores it in the value pointed to by out
//func decodeBody(resp *http.Response, out interface{}, equals io.WriteCloser) error {
//	defer resp.Body.Close()
//	// Symmetric API Testing
//	if equals != nil {
//		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, equals))
//		defer equals.Close()
//	}
//	decoder := json.NewDecoder(resp.Body)
//	return decoder.Decode(out)
//}
