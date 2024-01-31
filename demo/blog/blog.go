package main

import (
	"os"

	"github.com/goplus/yap"
)

func main() {
	y := yap.New(os.DirFS("."))
	y.SetDelims("${", "}")
	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.YAP(200, "article", yap.H{
			"id":   1,
			"Name": "aaaa",
			"Base": "bbbbb",
			"URL":  "http://fefawd.baidu.com",
			"Ann":  true,
			"Items": []string{
				"awdawd",
			},
		})
	})

	y.Run("127.0.0.1:8080")
}
