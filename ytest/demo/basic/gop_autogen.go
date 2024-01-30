package main

import "github.com/goplus/yap"

const _ = true

type foo struct {
	yap.App
}
//line ytest/demo/basic/foo.gox:7
func (this *foo) Main() {
//line ytest/demo/basic/foo.gox:7:1
	this.InitYap()
//line ytest/demo/basic/foo.gox:9:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line ytest/demo/basic/foo.gox:10:1
		ctx.Json__1(map[string]string{"id": ctx.Param("id")})
	})
}
func main() {
}
