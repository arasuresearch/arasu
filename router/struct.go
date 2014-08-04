package router

import (
	"reflect"
)

type Match struct {
	Path   string
	Method string
	Action string
}

type Route struct {
	Paths   Path
	Matches []Match
}
type Path map[string]*Route
type Router struct {
	All *Route
}

type Handler struct {
	Type      reflect.Type
	Methods   []string
	FuncsArgs map[string][]string
}
type Registry map[string]Handler
