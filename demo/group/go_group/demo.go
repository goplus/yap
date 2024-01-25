package main

import (
	"github.com/goplus/yap"
)

func main() {

	y := yap.New()
	rg := y.Group("/vip")
	rg.GET("/v1", func(ctx *yap.Context) {
		ctx.JSON(200, "v1")
	})
	rg.GET("/v2", func(ctx *yap.Context) {
		ctx.JSON(200, "v2")
	})
	rg1 := y.Group("/xhy")
	rg1.GET("/a", func(ctx *yap.Context) {
		ctx.JSON(200, "a")
	})
	rg1.GET("/b", func(ctx *yap.Context) {
		ctx.JSON(200, "b")
	})
	y.Run(":8082")
}
