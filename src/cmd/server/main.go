package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFile("/logo.png", "/static/assets/amazing_logo.png")
	r.LoadHTMLGlob("static/templates/*")
	r.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	r.Run()
}
