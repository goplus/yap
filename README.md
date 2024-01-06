yap - Yet Another Go/Go+ HTTP Web Framework
======

[![Build Status](https://github.com/goplus/yap/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/yap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/yap)](https://goreportcard.com/report/github.com/goplus/yap)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/yap.svg?label=release)](https://github.com/goplus/yap/releases)
[![Coverage Status](https://codecov.io/gh/goplus/yap/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/yap)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/yap.svg)](https://pkg.go.dev/github.com/goplus/yap)

### Router and Parameters

demo ([hello.go](demo/hello/hello.go)):

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

### YAP Template

demo ([blog.go](demo/blog/blog.go), [article.yap](demo/blog/yap/article.yap)):

```go
import (
	"embed"
	"io/fs"

	"github.com/goplus/yap"
)

type article struct {
	ID string
}

//go:embed yap
var yapFS embed.FS

fsYap, _ := fs.Sub(yapFS, "yap")
y := yap.New(fsYap)

y.GET("/p/:id", func(ctx *yap.Context) {
	ctx.YAP(200, "article", article{
		ID: ctx.Param("id"),
	})
})

y.Run(":8080")
```
