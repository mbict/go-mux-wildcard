package mux

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

var ErrDuplicateRoute = errors.New("duplicate route")

type route struct {
	path    string
	childs  routes
	handler http.Handler
}

type routes []*route

func (t *routes) AddRoute(path string, h http.Handler) error {
	var addToSegment func(t *routes, segments []string) error
	addToSegment = func(t *routes, segments []string) error {
		isLastSegment := len(segments) == 1
		segment := segments[0]

		var foundPathSegment *route
		for _, v := range *t {
			if v.path == segment {
				if v.handler != nil && isLastSegment {
					return ErrDuplicateRoute
				}

				foundPathSegment = v
				break
			}
		}

		//no path segment found
		if foundPathSegment == nil {
			foundPathSegment = &route{
				path:    segment,
				childs:  routes{},
				handler: nil,
			}
			*t = append(*t, foundPathSegment)

			sort.Slice(*t, func(i, j int) bool {
				return len((*t)[i].path) >= len((*t)[j].path)
			})
		}

		if isLastSegment {
			foundPathSegment.handler = h
			return nil
		}

		return addToSegment(&foundPathSegment.childs, segments[1:])
	}

	segments := splitPathWildcards(path)

	return addToSegment(t, segments)
}

func (t *routes) Match(path string) http.Handler {
	for _, r := range *t {
		n := len(r.path)
		if len(path) >= n && path[0:n] == r.path {
			//check child routes
			if len(r.childs) > 0 {
				//fmt.Println("===" ,path)
				path := path[n:]

				//skip
				if len(path) >= 1 {
					//skip if this is not the exact path
					if path[0] != '/' {
						continue
					}

					//remove param
					path = removeParam(path[1:])
					if h := r.childs.Match(path); h != nil {
						return h
					}
				}
			}
			return r.handler
		}
	}
	return nil
}

func splitPathWildcards(path string) []string {
	res := []string{}
	for i := strings.Index(path, "/*/"); i >= 0; i = strings.Index(path, "/*/") {
		res = append(res, path[:i])
		path = path[i+2:]
	}

	res = append(res, path)
	return res
}

func removeParam(path string) string {
	if i := strings.IndexByte(path, '/'); i >= 0 {
		return path[i:]
	}
	return ""
}
