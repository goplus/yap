// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/goplus/yap/ytest"
	"github.com/qiniu/x/test"
)

const _ = true

type complex struct {
	ytest.Case
	*App
}
type App struct {
	ytest.App
}
//line ytest/demo/match/complex/complex_yapt.gox:1
func (this *complex) Main() {
//line ytest/demo/match/complex/complex_yapt.gox:1:1
	d := test.Gopx_Var_Cast__0[string]()
//line ytest/demo/match/complex/complex_yapt.gox:3:1
	test.Gopt_Case_MatchAny(this, map[string]map[string]*test.Var__0[string]{"c": map[string]*test.Var__0[string]{"d": d}}, map[string]interface{}{"a": 1, "b": 3.14, "c": map[string]string{"d": "hello", "e": "world"}, "f": 1})
//line ytest/demo/match/complex/complex_yapt.gox:12:1
	fmt.Println(d)
//line ytest/demo/match/complex/complex_yapt.gox:13:1
	test.Gopt_Case_MatchAny(this, d, "hello")
}
func (this *App) Main() {
	ytest.Gopt_App_Main(this, new(complex))
}
func main() {
	new(App).Main()
}
