// Code generated by gop (Go+); DO NOT EDIT.

package main

import "github.com/goplus/yap"

const _ = true

type get struct {
	yap.Handler
	*AppV2
}
type get_p_id struct {
	yap.Handler
	*AppV2
}
type AppV2 struct {
	yap.AppV2
}
//line demo/classfile2_hello/get.yap:1
func (this *get) Main(_gop_arg0 *yap.Context) {
	this.Handler.Main(_gop_arg0)
//line demo/classfile2_hello/get.yap:1:1
	this.Html__1(`<html><body>Hello, YAP!</body></html>`)
}
func (this *get) Classfname() string {
	return "get"
}
func (this *get) Classclone() yap.HandlerProto {
	_gop_ret := *this
	return &_gop_ret
}
//line demo/classfile2_hello/get_p_#id.yap:1
func (this *get_p_id) Main(_gop_arg0 *yap.Context) {
//line demo/classfile2_hello/get.yap:1:1
	this.Handler.Main(_gop_arg0)
//line demo/classfile2_hello/get_p_#id.yap:1:1
	this.Json__1(map[string]string{"id": this.Gop_Env("id")})
}
func (this *get_p_id) Classfname() string {
	return "get_p_#id"
}
func (this *get_p_id) Classclone() yap.HandlerProto {
	_gop_ret := *this
	return &_gop_ret
}
func (this *AppV2) Main() {
	yap.Gopt_AppV2_Main(this, new(get), new(get_p_id))
}
func main() {
	new(AppV2).Main()
}
