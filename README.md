# gorest
[![Go Report Card](https://goreportcard.com/badge/github.com/izumix03/gorest)](https://goreportcard.com/report/github.com/izumix03/gorest)
[![Build Status](https://travis-ci.org/izumix03/gorest.svg?branch=master)](https://travis-ci.org/izumix03/gorest)

## install 
```$xslt
go get github.com/izumix03/gorest
```

## usage
```$xslt
res, err := gorest.Get(`http://example.com`).
		Path(`/ticket`).
		URLParam(`name`, `foo`).
		Execute()
```


```$xslt
res, err := gorest.Post(`http://example.com`).
		Path(`/ticket`).
		JSONStruct(Foo{name: `bar`}).
		Execute()
```
