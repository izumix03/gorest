# gorest
[![Go Report Card](https://goreportcard.com/badge/github.com/izumix03/gorest)](https://goreportcard.com/report/github.com/izumix03/gorest)
[![Build Status](https://travis-ci.org/izumix03/gorest.svg?branch=master)](https://travis-ci.org/izumix03/gorest)

## install 
```$xslt
go get github.com/izumix03/gorest
```

## usage
```go
_, err := gorest.Get(`http://example.com`).
		Path(`/ticket`).
		URLParam(`name`, `foo`).
		Execute()
```


```go
_, err := gorest.Post(`http://example.com`).
		Path(`/ticket`).
		JSONStruct(Foo{name: `bar`}).
		Execute()
```

HandleBody validates http status code(default over 400 is error),  
auto close response body.

```go
var response Ticket
err := gorest.Put(`http://example.com`).
	Path(`/ticket/%s`, ticket.ID).
	Header("X-Api-Key", token.AccessToken).
	JSONStruct(Foo{name: `bazz`}).
	HandleBody(func(body []uint8) error {
	    if err := json.Unmarshal(body, &response); err != nil {
                return err
            }
            return nil
        })
```
