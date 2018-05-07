# gorest

## usage
```$xslt
gorest.Get(`http://example.com`).
		AddPath(`/ticket`).
		AddURLParam(`name`, `foo`).
		Execute()
```
