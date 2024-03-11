package main

import (
	"github.com/goplus/yap"
)

func main() {
	y := yap.New()
	y.GET("/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`)
	})
	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.JSON(200, yap.H{
			"id": ctx.Param("id"),
		})
	})
	y.Run("localhost:8080")
}
