package gorest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_client_Execute_simple_json_post(t *testing.T) {
	content := struct {
		Name string
	}{
		Name: "name",
	}
	postBody, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal %s", err)
	}

	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("invalid method %s", r.Method)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed to read body, %s", err)
			}
			if diff := cmp.Diff(
				body,
				[]byte("{\"Name\":\"name\"}"),
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err := w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	response, err := Post(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		JSON(postBody).
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_Execute_simple_urlEncoded_put(t *testing.T) {
	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Fatalf("invalid method %s", r.Method)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed to read body, %s", err)
			}
			if diff := cmp.Diff(
				body,
				[]byte("key=value"),
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err := w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	response, err := Put(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		URLEncoded("key", "value").
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_Execute_multi_urlEncoded_post(t *testing.T) {
	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("invalid method %s", r.Method)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed to read body, %s", err)
			}
			if diff := cmp.Diff(
				body,
				[]byte("key=value&key2=value2&key2=value3"),
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err := w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	response, err := Post(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		URLEncoded("key", "value").
		URLEncodedList("key2", []string{"value2", "value3"}).
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_Execute_simple_multipart_value_post(t *testing.T) {
	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("invalid method %s", r.Method)
			}

			value := r.FormValue("key")
			if diff := cmp.Diff(
				value,
				`(1,2,34,55,666)`,
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err := w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	response, err := Post(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		MultipartData("key", strings.NewReader("(1,2,34,55,666)"), false).
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_Execute_simple_multipart_file_post(t *testing.T) {
	fileName := "testdata/sample.golden"

	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("invalid method %s", r.Method)
			}

			file, header, err := r.FormFile("key")
			if err != nil {
				t.Fatalf("failed to read from file  %s", err)
			}
			defer file.Close()

			if header.Filename != fileName {
				t.Fatalf("invalid file name  got = %s", header.Filename)
			}

			buf := bytes.NewBuffer(nil)
			if _, err = io.Copy(buf, file); err != nil {
				t.Fatalf("failed to copy multipart file = %s", err)
			}

			if diff := cmp.Diff(
				buf.String(),
				"header1,header2\nvalue1,value2\n",
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err = w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	f, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("cannot open file %q: %v", "testdata/sample.golden", err)
	}
	defer f.Close()

	response, err := Post(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		MultipartData("key", f, true).
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_Execute_simple_multipart_value_as_file_post(t *testing.T) {
	fileName := "sample.csv"

	var remoteURL string
	{
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("invalid method %s", r.Method)
			}

			file, header, err := r.FormFile("key")
			if err != nil {
				t.Fatalf("failed to read from file  %s", err)
			}
			defer file.Close()

			if header.Filename != fileName {
				t.Fatalf("invalid file name  got = %s", header.Filename)
			}

			buf := bytes.NewBuffer(nil)
			if _, err = io.Copy(buf, file); err != nil {
				t.Fatalf("failed to copy multipart file = %s", err)
			}

			if diff := cmp.Diff(
				buf.String(),
				"header1,header2\nvalue1,value2\n",
			); diff != "" {
				t.Fatalf("invalid postBody, diff = %s", diff)
			}
			if _, err := w.Write([]byte("success")); err != nil {
				t.Fatalf("failed to write response %s", err)
			}
		}))
		defer server.Close()
		remoteURL = server.URL
	}

	response, err := Post(remoteURL).
		Client(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}).
		MultipartAsFormFile("key", fileName, strings.NewReader("header1,header2\nvalue1,value2\n"), false).
		Execute()
	if err != nil {
		t.Errorf("failed to post %s", err)
		return
	}
	defer CloseBody(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response %s", err)
	}
	if !reflect.DeepEqual(body, []byte("success")) {
		t.Fatalf("wrong response, got => %s", body)
	}
}

func Test_client_buildParams(t *testing.T) {
	type fields struct {
		contentType          contentType
		params               interface{}
		hasJsonStruct        bool
		hasRawFormUrlEncoded bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    io.Reader
		wantErr bool
	}{
		{
			name: "json_content",
			fields: fields{
				contentType:   "application/json",
				params:        []byte("{\"one\": \"two\"}"),
				hasJsonStruct: false,
			},
			want:    bytes.NewBuffer([]byte("{\"one\": \"two\"}")),
			wantErr: false,
		},
		{
			name: "url_encoded_raw",
			fields: fields{
				contentType:          "application/x-www-form-urlencoded",
				params:               "a=1&b=2",
				hasJsonStruct:        false,
				hasRawFormUrlEncoded: true,
			},
			want:    strings.NewReader("a=1&b=2"),
			wantErr: false,
		},
		{
			name: "url_encoded",
			fields: fields{
				contentType:          "application/x-www-form-urlencoded",
				params:               "a%3D1%26b%3D2",
				hasJsonStruct:        false,
				hasRawFormUrlEncoded: true,
			},
			want:    strings.NewReader("a%3D1%26b%3D2"),
			wantErr: false,
		},
		{
			name: "url_encoded_not_escaped",
			fields: fields{
				contentType:          "application/x-www-form-urlencoded",
				params:               "a=1&b=2",
				hasJsonStruct:        false,
				hasRawFormUrlEncoded: false,
			},
			want:    strings.NewReader("a%3D1%26b%3D2"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				contentType:          tt.fields.contentType,
				params:               tt.fields.params,
				hasJsonStruct:        tt.fields.hasJsonStruct,
				hasRawFormUrlEncoded: tt.fields.hasRawFormUrlEncoded,
			}
			got, err := cli.buildParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("buildParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
