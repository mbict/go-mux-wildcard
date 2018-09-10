package mux

import (
	"net/http/httptest"
	"testing"
)

func TestTreeAddDuplicateRoute(t *testing.T) {
	h := serve(100)

	tt := &routes{}

	if err := tt.AddRoute("/foo", h); err != nil {
		t.Errorf("should add route but got error '%s'", err)
	}
	if err := tt.AddRoute("/foo", h); err != ErrDuplicateRoute {
		t.Errorf("should return duplicate route error but got '%s'", err)
	}
}

func TestTreeAddDuplicateRouteWithWildcard(t *testing.T) {
	h := serve(100)

	tt := &routes{}

	if err := tt.AddRoute("/foo/*/baz", h); err != nil {
		t.Errorf("should add route but got error '%s'", err)
	}
	if err := tt.AddRoute("/foo/*/baz", h); err != ErrDuplicateRoute {
		t.Errorf("should return duplicate route error but got '%s'", err)
	}
}

func TestTreeMatchHandler(t *testing.T) {
	testroutes := []struct {
		pattern string
		code    int
	}{
		{"/foo", 100},
		{"/foo/bar", 101},
		{"/foo/*/b", 102},
		{"/foo/*/baz", 103},
		{"/foo/*/*/baz", 104},
	}

	testpaths := []struct {
		path         string
		expectedCode int
	}{
		{"/foo", 100},
		{"/foo/lalalalal", 100},
		{"/foo/bar", 101},
		{"/foo/bar/b", 101},
		{"/foo/bar1/b", 101},
		{"/foo/baz/b", 102},
		{"/foo/a/baz", 103},
		{"/foo/a/b", 102},
		{"/foo/a/be", 102},
		{"/foo/a/ba", 102},
		{"/foo/a/baz", 103},
		{"/foo/a/bazert", 103},
		{"/foo/a/b/baz", 102},
		{"/foo/a/c/baz", 104},
	}

	tt := &routes{}

	for _, v := range testroutes {
		if err := tt.AddRoute(v.pattern, serve(v.code)); err != nil {
			t.Errorf("should add route `%s` without error but got error '%s'", v.pattern, err)
		}
	}

	for _, test := range testpaths {
		h := tt.Match(test.path)

		if h == nil {
			t.Errorf("expected a handler for route `%s`", test.path)
			continue
		}

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, nil)

		if rr.Code != test.expectedCode {
			t.Errorf("%s should return code %d but got %d", test.path, test.expectedCode, rr.Code)
		}
	}
}

func TestTreeMatchHandlerNoMatch(t *testing.T) {
	testroutes := []string{
		"/foo",
		"/bar/*/foo",
	}

	testpaths := []string{
		"/fo",
		"/bar",
		"/bar/blah",
		"/bar/blah/booboo",
	}

	tt := &routes{}

	for _, v := range testroutes {
		if err := tt.AddRoute(v, serve(100)); err != nil {
			t.Errorf("should add route `%s` without error but got error '%s'", v, err)
		}
	}

	for _, testpath := range testpaths {
		h := tt.Match(testpath)

		if h != nil {
			t.Errorf("expected a nil handler for route `%s`, but got a route back", testpath)
		}
	}
}
