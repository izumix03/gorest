package gorest

import (
	"github.com/google/go-cmp/cmp"
	"net/url"
	"reflect"
	"testing"
)

func Test_client_Path(t *testing.T) {
	type fields struct {
		method      requestMethod
		contentType contentType
		baseURL     string
		paths       []string
		urlParams   []string
	}
	type args struct {
		pathFmt string
		args    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TerminalOperator
	}{
		{
			name: "simple_path",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       nil,
			},
			args: args{
				pathFmt: "/users",
				args:    nil,
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
			},
		},
		{
			name: "formatted_path",
			fields: fields{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       nil,
			},
			args: args{
				pathFmt: "/users/%s",
				args:    []interface{}{"takahiro"},
			},
			want: &client{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users/takahiro"},
			},
		},
		{
			name: "twice_call",
			fields: fields{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
			},
			args: args{
				pathFmt: "/blog/%d",
				args:    []interface{}{1},
			},
			want: &client{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users", "/blog/1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:      tt.fields.method,
				contentType: tt.fields.contentType,
				baseURL:     tt.fields.baseURL,
				paths:       tt.fields.paths,
				urlParams:   tt.fields.urlParams,
			}
			if got := cli.Path(tt.args.pathFmt, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_URLParam(t *testing.T) {
	type fields struct {
		method      requestMethod
		contentType contentType
		baseURL     string
		paths       []string
		urlParams   []string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TerminalOperator
	}{
		{
			name: "call_once",
			fields: fields{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
				urlParams:   nil,
			},
			args: args{
				key:   "key",
				value: "value",
			},
			want: &client{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
				urlParams:   []string{"key=value"},
			},
		},
		{
			name: "call_twice",
			fields: fields{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
				urlParams:   []string{"key=value"},
			},
			args: args{
				key:   "key2",
				value: "value2",
			},
			want: &client{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
				paths:       []string{"/users"},
				urlParams:   []string{"key=value", "key2=value2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:      tt.fields.method,
				contentType: tt.fields.contentType,
				baseURL:     tt.fields.baseURL,
				paths:       tt.fields.paths,
				urlParams:   tt.fields.urlParams,
			}
			if got := cli.URLParam(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URLParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func strPtr(str string) *string {
	return &str
}

func Test_client_BasicAuth(t *testing.T) {
	type fields struct {
		method      requestMethod
		contentType contentType
		baseURL     string
		headers     map[string]string
		username    *string
		password    *string
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TerminalOperator
	}{
		{
			name: "simple_call",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    nil,
				password:    nil,
			},
			args: args{
				username: "name",
				password: "pass",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("name"),
				password:    strPtr("pass"),
			},
		},
		{
			name: "headers_not_reset",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				headers:     map[string]string{"key": "value"},
				username:    nil,
				password:    nil,
			},
			args: args{
				username: "name",
				password: "pass",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				headers:     map[string]string{"key": "value"},
				username:    strPtr("name"),
				password:    strPtr("pass"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:      tt.fields.method,
				contentType: tt.fields.contentType,
				baseURL:     tt.fields.baseURL,
				headers:     tt.fields.headers,
				username:    tt.fields.username,
				password:    tt.fields.password,
			}
			got := cli.BasicAuth(tt.args.username, tt.args.password)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(client{})); diff != "" {
				t.Errorf("BasicAuth() diff= %s", diff)
			}
		})
	}
}

func Test_client_Header(t *testing.T) {
	type fields struct {
		method      requestMethod
		contentType contentType
		baseURL     string
		username    *string
		password    *string
		headers     map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TerminalOperator
	}{
		{
			name: "simple_set",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     nil,
			},
			args: args{
				key:   "key",
				value: "value",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     map[string]string{"key": "value"},
			},
		},
		{
			name: "set_two_keys",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     map[string]string{"key": "value"},
			},
			args: args{
				key:   "key2",
				value: "value2",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     map[string]string{"key": "value", "key2": "value2"},
			},
		},
		{
			name: "update_if_same_key",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     map[string]string{"key": "value"},
			},
			args: args{
				key:   "key",
				value: "value2",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				username:    strPtr("user"),
				password:    strPtr("pass"),
				headers:     map[string]string{"key": "value2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:      tt.fields.method,
				contentType: tt.fields.contentType,
				baseURL:     tt.fields.baseURL,
				username:    tt.fields.username,
				password:    tt.fields.password,
				headers:     tt.fields.headers,
			}
			if got := cli.Header(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_JSON(t *testing.T) {
	type fields struct {
		method        requestMethod
		contentType   contentType
		baseURL       string
		params        interface{}
		hasJsonStruct bool
	}
	type args struct {
		json []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JSONContent
	}{
		{
			name: "simple_set",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				json: []byte("{\"json\": true}"),
			},

			want: &client{
				method:        "POST",
				contentType:   "application/json",
				baseURL:       "https://sample.com",
				params:        []byte("{\"json\": true}"),
				hasJsonStruct: false,
			},
		},
		{
			name: "set_0length",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				json: []byte{},
			},

			want: &client{
				method:        "POST",
				contentType:   "application/json",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:        tt.fields.method,
				contentType:   tt.fields.contentType,
				baseURL:       tt.fields.baseURL,
				params:        tt.fields.params,
				hasJsonStruct: tt.fields.hasJsonStruct,
			}
			got := cli.JSON(tt.args.json)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(client{})); diff != "" {
				t.Errorf("JSON() diff = %s", diff)
			}
		})
	}
}

func Test_client_JSONString(t *testing.T) {
	type fields struct {
		method        requestMethod
		contentType   contentType
		baseURL       string
		params        interface{}
		hasJsonStruct bool
	}
	type args struct {
		json string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JSONContent
	}{
		{
			name: "simple_set",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				json: "{\"json\": true}",
			},

			want: &client{
				method:        "POST",
				contentType:   "application/json",
				baseURL:       "https://sample.com",
				params:        []byte("{\"json\": true}"),
				hasJsonStruct: false,
			},
		},
		{
			name: "set_0length",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				json: "",
			},

			want: &client{
				method:        "POST",
				contentType:   "application/json",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:        tt.fields.method,
				contentType:   tt.fields.contentType,
				baseURL:       tt.fields.baseURL,
				params:        tt.fields.params,
				hasJsonStruct: tt.fields.hasJsonStruct,
			}
			got := cli.JSONString(tt.args.json)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(client{})); diff != "" {
				t.Errorf("JSONString() = %s", diff)
			}
		})
	}
}

func Test_client_JSONStruct(t *testing.T) {
	type fields struct {
		method        requestMethod
		contentType   contentType
		baseURL       string
		params        interface{}
		hasJsonStruct bool
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JSONContent
	}{
		{
			name: "simple_set",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				data: struct {
					name string
				}{
					name: "name",
				},
			},

			want: &client{
				method:      "POST",
				contentType: "application/json",
				baseURL:     "https://sample.com",
				params: struct {
					name string
				}{
					name: "name",
				},
				hasJsonStruct: true,
			},
		},
		{
			name: "can_set_invalid_argument",
			fields: fields{
				method:        "POST",
				contentType:   "",
				baseURL:       "https://sample.com",
				params:        nil,
				hasJsonStruct: false,
			},
			args: args{
				data: "invalid argument because not struct",
			},

			want: &client{
				method:        "POST",
				contentType:   "application/json",
				baseURL:       "https://sample.com",
				params:        "invalid argument because not struct",
				hasJsonStruct: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:        tt.fields.method,
				contentType:   tt.fields.contentType,
				baseURL:       tt.fields.baseURL,
				params:        tt.fields.params,
				hasJsonStruct: tt.fields.hasJsonStruct,
			}
			if got := cli.JSONStruct(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_URLEncoded(t *testing.T) {
	type fields struct {
		method      requestMethod
		contentType contentType
		baseURL     string
		params      interface{}
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   URLEncoded
	}{
		{
			name: "simple_set",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				params:      nil,
			},
			args: args{
				key:   "key",
				value: "value",
			},
			want: &client{
				method:      "POST",
				contentType: "application/x-www-form-urlencoded",
				baseURL:     "https://sample.com",
				params:      url.Values{"key": []string{"value"}},
			},
		},
		{
			name: "add_value_if_same_key",
			fields: fields{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
				params:      url.Values{"key": []string{"value"}},
			},
			args: args{
				key:   "key",
				value: "value2",
			},
			want: &client{
				method:      "POST",
				contentType: "application/x-www-form-urlencoded",
				baseURL:     "https://sample.com",
				params:      url.Values{"key": []string{"value", "value2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &client{
				method:      tt.fields.method,
				contentType: tt.fields.contentType,
				baseURL:     tt.fields.baseURL,
				params:      tt.fields.params,
			}
			if got := cli.URLEncoded(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URLEncoded() = %v, want %v", got, tt.want)
			}
		})
	}
}
