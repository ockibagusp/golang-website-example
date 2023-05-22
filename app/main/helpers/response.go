package helpers

import "net/http"

type Response struct {
	Code   int    `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	// // or,
	// Message   interface{} `form:"message" json:"message"`
	Data interface{} `form:"data" json:"data"`
}

var ResponseForbidden = Response{
	Code:   http.StatusForbidden,
	Status: http.StatusText(http.StatusForbidden),
	Data:   "FORBIDDEN",
}

var ResponseUnauthorized = Response{
	Code:   http.StatusUnauthorized,
	Status: http.StatusText(http.StatusUnauthorized),
	Data:   "UNAUTHORIZED",
}

var ResponseBadRequest = Response{
	Code:   http.StatusBadRequest,
	Status: http.StatusText(http.StatusBadRequest),
	Data:   "BAD REQUEST",
}
