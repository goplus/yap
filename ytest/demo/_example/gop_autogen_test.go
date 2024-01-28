package main

import (
	"fmt"
	"github.com/goplus/yap/ytest"
	"testing"
)

type case_example struct {
	ytest.Case
}
//line ytest/demo/example/example_ytest.gox:1
func (this *case_example) Main() {
//line ytest/demo/example/example_ytest.gox:1:1
	this.Host("https://example.com", "http://localhost:8080")
//line ytest/demo/example/example_ytest.gox:2:1
	testauth := ytest.Oauth2("...")
//line ytest/demo/example/example_ytest.gox:4:1
	this.DefaultHeader.Set("User-Agent", "yaptest/0.7")
//line ytest/demo/example/example_ytest.gox:6:1
	this.Run("urlWithVar", func() {
//line ytest/demo/example/example_ytest.gox:7:1
		id := "123"
//line ytest/demo/example/example_ytest.gox:8:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_ytest.gox:9:1
		this.Send()
//line ytest/demo/example/example_ytest.gox:10:1
		fmt.Println("code:", this.Resp().Code())
//line ytest/demo/example/example_ytest.gox:11:1
		fmt.Println("body:", this.Resp().Body())
	})
//line ytest/demo/example/example_ytest.gox:14:1
	this.Run("matchWithVar", func() {
//line ytest/demo/example/example_ytest.gox:15:1
		code := ytest.Gopx_Var_Cast__0[int]()
//line ytest/demo/example/example_ytest.gox:16:1
		id := "123"
//line ytest/demo/example/example_ytest.gox:17:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_ytest.gox:18:1
		this.RetWith(code)
//line ytest/demo/example/example_ytest.gox:19:1
		fmt.Println("code:", code)
//line ytest/demo/example/example_ytest.gox:20:1
		ytest.Match__4(code, 200)
	})
//line ytest/demo/example/example_ytest.gox:23:1
	this.Run("postWithAuth", func() {
//line ytest/demo/example/example_ytest.gox:24:1
		id := "123"
//line ytest/demo/example/example_ytest.gox:25:1
		title := "title"
//line ytest/demo/example/example_ytest.gox:26:1
		author := "author"
//line ytest/demo/example/example_ytest.gox:27:1
		this.Post("https://example.com/p/" + id)
//line ytest/demo/example/example_ytest.gox:28:1
		this.Auth(testauth)
//line ytest/demo/example/example_ytest.gox:29:1
		this.Json(map[string]string{"title": title, "author": author})
//line ytest/demo/example/example_ytest.gox:33:1
		this.RetWith(200)
//line ytest/demo/example/example_ytest.gox:34:1
		fmt.Println("body:", this.Resp().Body())
	})
//line ytest/demo/example/example_ytest.gox:37:1
	this.Run("matchJsonObject", func() {
//line ytest/demo/example/example_ytest.gox:38:1
		title := ytest.Gopx_Var_Cast__0[string]()
//line ytest/demo/example/example_ytest.gox:39:1
		author := ytest.Gopx_Var_Cast__0[string]()
//line ytest/demo/example/example_ytest.gox:40:1
		id := "123"
//line ytest/demo/example/example_ytest.gox:41:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_ytest.gox:42:1
		this.RetWith(200)
//line ytest/demo/example/example_ytest.gox:43:1
		this.Json(map[string]*ytest.Var__0[string]{"title": title, "author": author})
//line ytest/demo/example/example_ytest.gox:47:1
		fmt.Println("title:", title)
//line ytest/demo/example/example_ytest.gox:48:1
		fmt.Println("author:", author)
	})
}
func Test_example(t *testing.T) {
	ytest.Gopt_Case_TestMain(new(case_example), t)
}
