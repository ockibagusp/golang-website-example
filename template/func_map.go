package template

import (
	"html/template"
	"strconv"
)

// function map more
var FuncMapMore = func() template.FuncMap {
	list := template.FuncMap{
		"toString":      ToString,
		"hasPermission": HasPermission,
	}
	return list
}

// TODO, insyaallah
//
// Code: session_gorilla.Values["s_auth_type"]
// HTML: {{index $.session_gorilla.Values "is_auth_type" | toString}} ?
func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	// Add whatever other types you need
	default:
		return ""
	}
}

// function has parmission to User
func HasPermission(feature string) bool {
	return false
}
