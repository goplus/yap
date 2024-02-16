yap - Yet Another Go/Go+ HTTP Web Framework
======

[![Build Status](https://github.com/goplus/yap/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/yap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/yap)](https://goreportcard.com/report/github.com/goplus/yap)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/yap.svg?label=release)](https://github.com/goplus/yap/releases)
[![Coverage Status](https://codecov.io/gh/goplus/yap/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/yap)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/yap.svg)](https://pkg.go.dev/github.com/goplus/yap)

This repo contains three [Go+ classfiles](https://github.com/goplus/gop/blob/main/doc/classfile.md): `yap` (a HTTP Web Framework), `yaptest` (a HTTP Test Framework) and `ydb` (a Go+ Database Framework).

The classfile `yap` has the file suffix `_yap.gox`. The classfile `yaptest` has the file suffix `_ytest.gox`. And the classfile `ydb` has the file suffix `_ydb.gox`.

Before using `yap`, `yaptest` or `ydb`, you need to add `github.com/goplus/yap` to `go.mod`:

```sh
gop get github.com/goplus/yap@latest
```

For more details, see [YAP Web Framework Manual](doc/manual.md).

### How to use in Go+

First let us initialize a hello project:

```sh
gop mod init hello
```

Then we have it reference a classfile called `yap` as the HTTP Web Framework:

```sh
gop get github.com/goplus/yap@latest
```

We can use it to implement a static file server:

```coffee
static "/foo", FS("public")
static "/"    # Equivalent to static "/", FS("static")

run ":8080"
```

We can also add the ability to handle dynamic GET/POST requests:

```coffee
static "/foo", FS("public")
static "/"    # Equivalent to static "/", FS("static")

get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

Save this code to `hello_yap.gox` file and execute:

```sh
mkdir -p yap/static yap/public    # Static resources can be placed in these directories
gop mod tidy
gop run .
```

A simplest web program is running now. At this time, if you visit http://localhost:8080/p/123, you will get:

```
{"id":"123"}
```

### yap: HTTP Web Framework

This classfile has the file suffix `_yap.gox`.

#### Router and Parameters

demo in Go+ classfile ([hello_yap.gox](demo/classfile_hello/hello_yap.gox)):

```coffee
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

#### Static files

Static files server demo in Go+ classfile ([staticfile_yap.gox](demo/classfile_static/staticfile_yap.gox)):

```coffee
static "/foo", FS("public")
static "/"

run ":8080"
```

#### YAP Template

demo in Go+ classfile ([blog_yap.gox](demo/classfile_blog/blog_yap.gox), [article_yap.html](demo/classfile_blog/yap/article_yap.html)):

```coffee
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

### yaptest: HTTP Test Framework

This classfile has the file suffix `_ytest.gox`.

Suppose we have a web server named `foo` ([demo/foo/foo_yap.gox](ytest/demo/foo/foo_yap.gox)):

```go
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

Then we create a yaptest file ([demo/foo/foo_ytest.gox](ytest/demo/foo/foo_ytest.gox)):

```go
mock "foo.com", new(foo)

run "test get /p/$id", => {
	id := "123"
	get "http://foo.com/p/${id}"
	ret 200
	json {
		"id": id,
	}
}
```

The directive `mock` creates the `foo` server by [mockhttp](https://pkg.go.dev/github.com/qiniu/x/mockhttp). Then we call the directive `run` to run a subtest.

You can change the directive `mock` to `testServer` (see [demo/foo/bar_ytest.gox](ytest/demo/foo/bar_ytest.gox)), and keep everything else unchanged:

```go
testServer "foo.com", new(foo)

run "test get /p/$id", => {
	id := "123"
	get "http://foo.com/p/${id}"
	ret 200
	json {
		"id": id,
	}
}
```

The directive `testServer` creates the `foo` server by [net/http/httptest](https://pkg.go.dev/net/http/httptest#NewServer) and obtained a random port as the service address. Then it calls the directive [host](https://pkg.go.dev/github.com/goplus/yap/ytest#App.Host) to map the random service address to `foo.com`. This makes all other code no need to changed.

For more details, see [yaptest - Go+ HTTP Test Framework](ytest).

### ydb: Database Framework

This classfile has the file suffix `_ydb.gox`.

TODO
