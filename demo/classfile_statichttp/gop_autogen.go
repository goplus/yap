package main

import (
	"github.com/goplus/yap"
	"github.com/qiniu/x/http/fs"
)

const _ = true

type statichttp struct {
	yap.App
}
//line demo/classfile_statichttp/statichttp_yap.gox:3
func (this *statichttp) MainEntry() {
//line demo/classfile_statichttp/statichttp_yap.gox:3:1
	this.Static__2("/", fs.Http("https://goplus.org"), false)
//line demo/classfile_statichttp/statichttp_yap.gox:4:1
	this.Run(":8888")
}
func main() {
//line demo/classfile_statichttp/statichttp_yap.gox:4:1
	yap.Gopt_App_Main(new(statichttp))
}
