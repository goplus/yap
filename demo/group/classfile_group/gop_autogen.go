package main

import "github.com/goplus/yap"

type demo_yap struct {
	yap.App
}
//line demo/group/classfile_group/demo_yap.gox:1
func (this *demo_yap) MainEntry() {
//line demo/group/classfile_group/demo_yap.gox:1:1
	this.Group("/vip")
//line demo/group/classfile_group/demo_yap.gox:2:1
	this.Get("/v1", func(ctx *yap.Context) {
//line demo/group/classfile_group/demo_yap.gox:3:1
		ctx.Json__1(map[string]string{"v1": "v1"})
	})
//line demo/group/classfile_group/demo_yap.gox:7:1
	this.Get("/v2", func(ctx *yap.Context) {
//line demo/group/classfile_group/demo_yap.gox:8:1
		ctx.Json__1(map[string]string{"v2": "v2"})
	})
//line demo/group/classfile_group/demo_yap.gox:12:1
	this.Group("/xhy")
//line demo/group/classfile_group/demo_yap.gox:13:1
	this.Get("/a", func(ctx *yap.Context) {
//line demo/group/classfile_group/demo_yap.gox:14:1
		ctx.Json__1(map[string]string{"xhy": "a"})
	})
//line demo/group/classfile_group/demo_yap.gox:18:1
	this.Get("/b", func(ctx *yap.Context) {
//line demo/group/classfile_group/demo_yap.gox:19:1
		ctx.Json__1(map[string]string{"xhy": "b"})
	})
//line demo/group/classfile_group/demo_yap.gox:25:1
	this.Run__1(":8083")
}
func main() {
	yap.Gopt_App_Main(new(demo_yap))
}
