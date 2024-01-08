package main

import "github.com/goplus/yap"

type blog struct {
	yap.App
}
//line demo/classfile_blog/blog_yap.gox:1
func (this *blog) MainEntry() {
//line demo/classfile_blog/blog_yap.gox:1:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line demo/classfile_blog/blog_yap.gox:2:1
		ctx.Yap__1("article", map[string]string{"id": ctx.Param("id")})
	})
//line demo/classfile_blog/blog_yap.gox:7:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(blog))
}
