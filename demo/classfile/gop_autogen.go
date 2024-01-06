package main

import "github.com/goplus/yap"

type article struct {
	ID string
}
type blog struct {
	yap.App
}
//line demo/classfile/blog.yapx.gox:5
func (this *blog) MainEntry() {
//line demo/classfile/blog.yapx.gox:5:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line demo/classfile/blog.yapx.gox:6:1
		ctx.Yap__1("article", article{ID: ctx.Param("id")})
	})
//line demo/classfile/blog.yapx.gox:11:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(blog))
}
