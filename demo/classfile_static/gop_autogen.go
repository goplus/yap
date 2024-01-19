package main

import "github.com/goplus/yap"

type staticfile struct {
	yap.App
}
//line demo/classfile_static/staticfile_yap.gox:1
func (this *staticfile) MainEntry() {
//line demo/classfile_static/staticfile_yap.gox:1:1
	this.Static("/foo", this.FS("public"))
//line demo/classfile_static/staticfile_yap.gox:2:1
	this.Static("/")
//line demo/classfile_static/staticfile_yap.gox:4:1
	this.Run(":8888")
}
func main() {
	yap.Gopt_App_Main(new(staticfile))
}
