package main

import "github.com/goplus/yap"

const _ = true

type blog struct {
	yap.App
}
//line demo/classfile_delimiter/blog_yap.gox:1
func (this *blog) MainEntry() {
//line demo/classfile_delimiter/blog_yap.gox:1:1
	this.SetDelims("${", "}$")
//line demo/classfile_delimiter/blog_yap.gox:2:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line demo/classfile_delimiter/blog_yap.gox:3:1
		ctx.Yap__1("article", map[string]string{"id": ctx.Param("id")})
	})
//line demo/classfile_delimiter/blog_yap.gox:8:1
	this.Run(":8080")
}
func main() {
//line demo/classfile_delimiter/blog_yap.gox:8:1
	yap.Gopt_App_Main(new(blog))
}
