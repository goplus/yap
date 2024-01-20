yap - Yet Another Go/Go+ HTTP Web Framework
======

[![Build Status](https://github.com/goplus/yap/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/yap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/yap)](https://goreportcard.com/report/github.com/goplus/yap)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/yap.svg?label=release)](https://github.com/goplus/yap/releases)
[![Coverage Status](https://codecov.io/gh/goplus/yap/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/yap)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/yap.svg)](https://pkg.go.dev/github.com/goplus/yap)

### Router and Parameters

demo in Go ([hello.go](demo/hello/hello.go)):

```go
import "github.com/goplus/yap"

y := yap.New()
y.GET("/p/:id", func(ctx *yap.Context) {
	ctx.JSON(200, yap.H{
		"id": ctx.Param("id"),
	})
})
y.Handle("/", func(ctx *yap.Context) {
	ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">Yap</a>!</body></html>`)
})
y.Run(":8080")
```

demo in Go+ classfile ([hello_yap.gox](demo/classfile_hello/hello_yap.gox)):

```go
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}
handle "/", ctx => {
	ctx.html `<html><body>Hello, <a href="/p/123">Yap</a>!</body></html>`
}

run ":8080"
```

### Static files

Static files server demo in Go:

```go
y := yap.New(os.DirFS("."))
y.Static("/foo", y.FS("public"))
y.Static("/") // means: y.Static("/", y.FS("static"))
y.Run(":8888")
```

Static files server demo in Go+ classfile ([staticfile_yap.gox](demo/classfile_static/staticfile_yap.gox)):

```go
static "/foo", FS("public")
static "/"
run ":8888"
```

Static files server also can use a `http.FileSystem` instead of `fs.FS` object (See [yapserve](https://github.com/xushiwei/yapserve) for more information):

```go
import "github.com/qiniu/x/http/fs"

static "/", fs.http("https://goplus.org"), false // false means not allow to redirect
run ":8888"
```

### YAP Template

demo in Go ([blog.go](demo/blog/blog.go), [article_yap.html](demo/blog/yap/article_yap.html)):

```go
import (
	"os"

	"github.com/goplus/yap"
)

y := yap.New(os.DirFS("."))

y.GET("/p/:id", func(ctx *yap.Context) {
	ctx.YAP(200, "article", yap.H{
		"id": ctx.Param("id"),
	})
})

y.Run(":8080")
```

demo in Go+ classfile ([blog_yap.gox](demo/classfile_blog/blog_yap.gox), [article_yap.html](demo/classfile_blog/yap/article_yap.html)):

```go
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```
