package mux

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func serve(code int) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(code)
	}
}

func TestMux(t *testing.T) {
	testing.Short()
	mux := New()

	routes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"/test", serve(200)},
		{"/testa", serve(201)},
		{"/test/long", serve(202)},
		{"/test/long/", serve(203)},
		{"/test/a/wildcard", serve(204)},
		{"/test/*/wildcard", serve(205)},
		{"/api/test", serve(206)},
		{"/api/test/*/foobar", serve(207)},
	}

	for _, r := range routes {
		mux.Handle(r.pattern, r.handler)
	}

	tc := map[string]struct {
		path string
		code int
	}{
		"longest path":                       {"/test/long/is/good", 203},
		"longest path with trailing slash":   {"/test/long/", 203},
		"longest path without traling slash": {"/test/long", 202},
		"shortest path":                      {"/test", 200},
		"variant a shortest path":            {"/testa", 201},
		"path ignoring the wildcard":         {"/test/a/wildcard", 204},
		"path with wildcard":                 {"/test/foo/wildcard", 205},

		"base path":                           {"/api/test", 206},
		"base path with params":               {"/api/test/123", 206},
		"base path with partial match on alt": {"/api/test/123/foo", 206},
		"alt path for base path":              {"/api/test/1234/foobar", 207},

		"no partial match": {"/tes", 404},
	}

	for name, test := range tc {
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: test.path,
			},
		}
		h := mux.Handler(req)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != test.code {
			t.Errorf("%s : %s want code %d got %d", name, test.path, rr.Code, test.code)
		}
	}
}

func TestMuxPanicDuplicatePath(t *testing.T) {
	expected := "mux: duplicate path for `/test`"
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
		} else if r != expected {
			t.Errorf("Expected panic with message `%s` but got `%s`", expected, r)
		}
	}()

	mux := New()
	mux.HandleFunc("/test", serve(200))
	mux.HandleFunc("/test", serve(200))
}

func TestMuxPanicNilHandler(t *testing.T) {
	expected := "mux: nil handler"
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
		} else if r != expected {
			t.Errorf("Expected panic with message `%s` but got `%s`", expected, r)
		}
	}()

	mux := New()
	mux.HandleFunc("/test", nil)
}
