[![Build Status](https://travis-ci.org/mbict/go-mux-wildcard.png?branch=master)](https://travis-ci.org/mbict/go-mux-wildcard)
[![GoDoc](https://godoc.org/github.com/mbict/go-mux-wildcard?status.png)](http://godoc.org/github.com/mbict/go-mux-wildcard)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-mux-wildcard)](http://gocover.io/github.com/mbict/go-mux-wildcard)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-mux-wildcard)](http://goreportcard.com/report/mbict/go-mux-wildcard)

# Mux Router with wildcard support

A simplified / stripdown version of the golang http.ServeMux with wildcard support

## Why i created this

I created this simplified version to be able to match on partial paths.
The muxer is mainly used for preselection of handlers based on the beginning of the path.

This is something i do when creating go-kit services.

## How it works

The matcher will start at the longest possible path an tries to (partial) match it with the request path.
If it found a match the http.Handler will be invoked and matching stops. 
When the path pattern does not match it will try the next path pattern.
If none of the path patterns match, the default not found handler will be invoked. 

## Wildcards
An path segment can be wildcarded with the * and must be followed with a next path name.

```go
    mux.Handle("/api/base", handler)
    mux.Handle("/api/base/*/subresource", handler)
    mux.Handle("/api/base/*/otherpart", handler)
```

## Parts i removed

- Redirect explicit slash
- Exact match on routes not ending with a slash
- Hostname matching
- Mutex locks (this muxer is initialized once at startup)