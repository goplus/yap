package main

import (
	"fmt"
	"github.com/goplus/yap/ytest"
)

const _ = true

func main() {
//line ytest/demo/example/example_yunit.gox:37:1
	ytest.Gopt_App_Main(new(ytest.App), new(example))
}

type example struct {
	ytest.Case
}
//line ytest/demo/example/example_yunit.gox:1
func (this *example) Main() {
//line ytest/demo/example/example_yunit.gox:1:1
	this.Host("https://example.com", "http://localhost:8080")
//line ytest/demo/example/example_yunit.gox:2:1
	testauth := ytest.Oauth2("...")
//line ytest/demo/example/example_yunit.gox:4:1
	this.DefaultHeader.Set("User-Agent", "yaptest/0.7")
//line ytest/demo/example/example_yunit.gox:6:1
	this.Run("urlWithVar", func() {
//line ytest/demo/example/example_yunit.gox:7:1
		id := "123"
//line ytest/demo/example/example_yunit.gox:8:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_yunit.gox:9:1
		this.Send()
//line ytest/demo/example/example_yunit.gox:10:1
		fmt.Println("code:", this.Resp().Code())
//line ytest/demo/example/example_yunit.gox:11:1
		fmt.Println("body:", this.Resp().Body())
	})
//line ytest/demo/example/example_yunit.gox:14:1
	this.Run("matchWithVar", func() {
//line ytest/demo/example/example_yunit.gox:15:1
		code := ytest.Gopx_Var_Cast__0[int]()
//line ytest/demo/example/example_yunit.gox:16:1
		id := "123"
//line ytest/demo/example/example_yunit.gox:17:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_yunit.gox:18:1
		this.RetWith(code)
//line ytest/demo/example/example_yunit.gox:19:1
		fmt.Println("code:", code)
//line ytest/demo/example/example_yunit.gox:20:1
		ytest.Match__4(code, 200)
	})
//line ytest/demo/example/example_yunit.gox:23:1
	this.Run("postWithAuth", func() {
//line ytest/demo/example/example_yunit.gox:24:1
		id := "123"
//line ytest/demo/example/example_yunit.gox:25:1
		title := "title"
//line ytest/demo/example/example_yunit.gox:26:1
		author := "author"
//line ytest/demo/example/example_yunit.gox:27:1
		this.Post("https://example.com/p/" + id)
//line ytest/demo/example/example_yunit.gox:28:1
		this.Auth(testauth)
//line ytest/demo/example/example_yunit.gox:29:1
		this.Json(map[string]string{"title": title, "author": author})
//line ytest/demo/example/example_yunit.gox:33:1
		this.RetWith(200)
//line ytest/demo/example/example_yunit.gox:34:1
		fmt.Println("body:", this.Resp().Body())
	})
//line ytest/demo/example/example_yunit.gox:37:1
	this.Run("mathJsonObject", func() {
//line ytest/demo/example/example_yunit.gox:38:1
		title := ytest.Gopx_Var_Cast__0[string]()
//line ytest/demo/example/example_yunit.gox:39:1
		author := ytest.Gopx_Var_Cast__0[string]()
//line ytest/demo/example/example_yunit.gox:40:1
		id := "123"
//line ytest/demo/example/example_yunit.gox:41:1
		this.Get("https://example.com/p/" + id)
//line ytest/demo/example/example_yunit.gox:42:1
		this.RetWith(200)
//line ytest/demo/example/example_yunit.gox:43:1
		this.Json(map[string]*ytest.Var__0[string]{"title": title, "author": author})
//line ytest/demo/example/example_yunit.gox:47:1
		fmt.Println("title:", title)
//line ytest/demo/example/example_yunit.gox:48:1
		fmt.Println("author:", author)
	})
}
