package helpers

type ResponseSuccess struct {
	Code   int         `form:"code" json:"code"`
	Status string      `form:"status" json:"status"`
	Data   interface{} `form:"data" json:"data"`
}

type ResponseError struct {
	Code    int         `form:"code" json:"code"`
	Status  string      `form:"status" json:"status"`
	Message interface{} `form:"message" json:"message"`
}
