// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"github.com/goplus/yap/ytest"
	"github.com/qiniu/x/stringutil"
	"github.com/qiniu/x/test"
	"testing"
)

type case_bar struct {
	ytest.CaseApp
}
type case_foo struct {
	ytest.CaseApp
}
//line ytest/demo/foo/bar_ytest.gox:1
func (this *case_bar) Main() {
//line ytest/demo/foo/bar_ytest.gox:1:1
	this.TestServer("foo.com", new(AppV2))
//line ytest/demo/foo/bar_ytest.gox:3:1
	id := "123"
//line ytest/demo/foo/bar_ytest.gox:4:1
	this.Get(stringutil.Concat("http://foo.com/p/", id))
//line ytest/demo/foo/bar_ytest.gox:5:1
	this.Send()
//line ytest/demo/foo/bar_ytest.gox:6:1
	test.Gopt_Case_MatchTBase(this, this.Resp().Code(), 200)
//line ytest/demo/foo/bar_ytest.gox:7:1
	this.Json(map[string]string{"id": id})
}
//line ytest/demo/foo/foo_ytest.gox:1
func (this *case_foo) Main() {
//line ytest/demo/foo/foo_ytest.gox:1:1
	this.Mock("foo.com", new(AppV2))
//line ytest/demo/foo/foo_ytest.gox:3:1
	id := "123"
//line ytest/demo/foo/foo_ytest.gox:4:1
	this.Get(stringutil.Concat("http://foo.com/p/", id))
//line ytest/demo/foo/foo_ytest.gox:5:1
	this.RetWith(200)
//line ytest/demo/foo/foo_ytest.gox:6:1
	this.Json(map[string]string{"id": id})
}
func Test_bar(t *testing.T) {
	ytest.Gopt_CaseApp_TestMain(new(case_bar), t)
}
func Test_foo(t *testing.T) {
	ytest.Gopt_CaseApp_TestMain(new(case_foo), t)
}
