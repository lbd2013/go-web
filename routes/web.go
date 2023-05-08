package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterWebRoutes 注册 Web 相关路由
func RegisterWebRoutes(r *gin.Engine) {
	//加载静态文件
	r.Static("/web", "public\\web")

	////模板解析，解析所有templates目录下的资源
	//r.LoadHTMLFiles("public\\web\\index.html")

	////模板渲染。请求/login会返回templates/login.html文件
	//r.GET("/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", nil)
	//})
}
