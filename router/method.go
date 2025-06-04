package router

type Method int

const (
	GET Method = iota
	HEAD
	POST
	PUT
	DELETE
	OPTIONS
	PATCH
)

var methodString = map[Method]string{
	GET:     "GET",
	HEAD:    "HEAD",
	POST:    "POST",
	PUT:     "PUT",
	DELETE:  "DELETE",
	OPTIONS: "OPTIONS",
	PATCH:   "PATCH",
}
