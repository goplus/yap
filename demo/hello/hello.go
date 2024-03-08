package main

import (
	"github.com/goplus/yap"
)

func main() {
	y := yap.New()
	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.JSON(200, yap.H{
			"id": ctx.Param("id"),
		})
	})
	y.Handle("/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`)
	})
	y.Run(":8080")
}
