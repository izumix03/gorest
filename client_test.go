package gorest

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want TerminalOperator
	}{
		{
			name: "method_is_get",
			args: args{
				baseURL: "https://sample.com",
			},
			want: &client{
				method:      "GET",
				contentType: "",
				baseURL:     "https://sample.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.baseURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want TerminalOperator
	}{
		{
			name: "method_is_post",
			args: args{
				baseURL: "https://sample.com",
			},
			want: &client{
				method:      "POST",
				contentType: "",
				baseURL:     "https://sample.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Post(tt.args.baseURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPut(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want TerminalOperator
	}{
		{
			name: "method_is_put",
			args: args{
				baseURL: "https://sample.com",
			},
			want: &client{
				method:      "PUT",
				contentType: "",
				baseURL:     "https://sample.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Put(tt.args.baseURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() = %v, want %v", got, tt.want)
			}
		})
	}
}
