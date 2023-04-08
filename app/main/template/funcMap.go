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

// Code: session_gorilla.Values["role"]
// HTML: {{index $.session_gorilla.Values "role" | toString}} ?
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
