package test

import (
	"net/http"
	"testing"

	c "github.com/ockibagusp/golang-website-example/controllers"
	"github.com/ockibagusp/golang-website-example/router"
	"github.com/stretchr/testify/assert"
)

// setup test Handler
func setupTestHandler() http.Handler {
	return router.New(controller)
}

// Controller test
var controller *c.Controller = &c.Controller{
	DB: db,
}

func TestController(t *testing.T) {
	/*
		assert := assert.New(t)
		assert.NotNil(controller.DB)

		or,
	*/
	assert.NotNil(t, controller.DB)
}
