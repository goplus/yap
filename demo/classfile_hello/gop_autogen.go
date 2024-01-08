package main

import "github.com/goplus/yap"

type hello struct {
	yap.App
}
//line demo/classfile_hello/hello_yap.gox:1
func (this *hello) MainEntry() {
//line demo/classfile_hello/hello_yap.gox:1:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line demo/classfile_hello/hello_yap.gox:2:1
		ctx.Json__1(map[string]string{"id": ctx.Param("id")})
	})
//line demo/classfile_hello/hello_yap.gox:6:1
	this.Handle("/", func(ctx *yap.Context) {
//line demo/classfile_hello/hello_yap.gox:7:1
		ctx.Html__1(`<html><body>Hello, <a href="/p/123">Yap</a>!</body></html>`)
	})
//line demo/classfile_hello/hello_yap.gox:10:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(hello))
}
