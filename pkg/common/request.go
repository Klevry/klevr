package common

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
}

func (r *Request) Param(name string) string {
	return strings.Join(r.URL.Query()[name], "")
}
