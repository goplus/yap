package main

import (
	"os"

	"github.com/goplus/yap"
)

func main() {
	y := yap.New(os.DirFS("."))

	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.YAP(200, "article", yap.H{
			"id": ctx.Param("id"),
		})
	})

	y.Run(":8080")
}
