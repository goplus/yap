// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/qiniu/x/test"
	"github.com/goplus/yap/ytest"
)

const _ = true

type simple struct {
	ytest.Case
	*App
}
type App struct {
	ytest.App
}
//line ytest/demo/match/simple/simple_yapt.gox:1
func (this *simple) Main() {
//line ytest/demo/match/simple/simple_yapt.gox:1:1
	id := test.Gopx_Var_Cast__0[int]()
//line ytest/demo/match/simple/simple_yapt.gox:2:1
	test.Gopt_Case_MatchAny(this, id, 1+2)
//line ytest/demo/match/simple/simple_yapt.gox:3:1
	test.Gopt_Case_MatchAny(this, id, 3)
//line ytest/demo/match/simple/simple_yapt.gox:4:1
	fmt.Println(id)
//line ytest/demo/match/simple/simple_yapt.gox:6:1
	test.Gopt_Case_MatchAny(this, id, 5)
}
func (this *App) Main() {
	ytest.Gopt_App_Main(this, new(simple))
}
func main() {
	new(App).Main()
}
