# gorest

## usage
```$xslt
result, err := gorest.Get(`http://example.com`).
		Path(`/ticket`).
		URLParam(`name`, `foo`).
		Execute()
```


```$xslt
result, err := gorest.Post(`http://example.com`).
		Path(`/ticket`).
		JSONStruct(Foo{name: `bar`}).
		Execute()
```
