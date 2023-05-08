package v1

import (
	"goweb/pkg/response"

	"github.com/gin-gonic/gin"
)

type DemoController struct {
	BaseAPIController
}

func (ctrl *DemoController) Index(c *gin.Context) {
	response.Data(c, "hello world")
}
