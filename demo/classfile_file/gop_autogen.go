package main

import "github.com/goplus/yap"

type hello struct {
	yap.App
}
//line classfile_file/hello_yap.gox:1
func (this *hello) MainEntry() {
//line classfile_file/hello_yap.gox:1:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(hello))
}
