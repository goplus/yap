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

### YAP Template

demo in Go ([blog.go](demo/blog/blog.go), [article.yap](demo/blog/yap/article.yap)):

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

demo in Go+ classfile ([blog_yap.gox](demo/classfile_blog/blog_yap.gox), [article.yap](demo/classfile_blog/yap/article.yap)):

```go
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```
