package helpers

type Response struct {
	Code   int    `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	// // or,
	// Data   interface{} `form:"data" json:"data"`
	Message interface{} `form:"message" json:"message"`
}
