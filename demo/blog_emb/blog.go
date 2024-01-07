package main

import (
	"embed"
	"io/fs"

	"github.com/goplus/yap"
)

type article struct {
	ID string
}

//go:embed yap
var yapFS embed.FS

func main() {
	fsYap, _ := fs.Sub(yapFS, "yap")
	y := yap.New(fsYap)

	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.YAP(200, "article", article{
			ID: ctx.Param("id"),
		})
	})

	y.Run(":8080")
}
