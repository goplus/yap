package main

import (
	"os"

	"github.com/goplus/yap"
)

func main() {
	y := yap.New(os.DirFS("."))

	y.GET("/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`)
	})
	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.YAP(200, "article", yap.H{
			"id": ctx.Param("id"),
		})
	})

	y.Run(":8888")
}
