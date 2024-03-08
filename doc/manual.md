yap - Yet Another Go/Go+ HTTP Web Framework
======

This repo contains three [Go+ classfiles](https://github.com/goplus/gop/blob/main/doc/classfile.md): `yap` (a HTTP Web Framework), `yaptest` (a HTTP Test Framework) and `ydb` (a Go+ Database Framework).

The classfile `yap` has the file suffix `.yap`. The classfile `yaptest` has the file suffix `_ytest.gox`. And the classfile `ydb` has the file suffix `_ydb.gox`.

Before using `yap`, `yaptest` or `ydb`, you need to add `github.com/goplus/yap` to `go.mod`:

```sh
gop get github.com/goplus/yap@latest
```


### Router and Parameters

demo in Go ([hello.go](../demo/hello/hello.go)):

```go
import "github.com/goplus/yap"

y := yap.New()
y.GET("/", func(ctx *yap.Context) {
	ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`)
})
y.GET("/p/:id", func(ctx *yap.Context) {
	ctx.JSON(200, yap.H{
		"id": ctx.Param("id"),
	})
})
y.Run(":8080")
```

demo in Go+ classfile ([main.yap](../demo/classfile_hello/main.yap)):

```go
get "/", ctx => {
	ctx.html `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`
}
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run "localhost:8080"
```


### Static files

Static files server demo in Go:

```go
y := yap.New(os.DirFS("."))
y.Static("/foo", y.FS("public"))
y.Static("/") // means: y.Static("/", y.FS("static"))
y.Run(":8080")
```

Static files server demo in Go+ classfile ([staticfile_yap.gox](../demo/classfile_static/staticfile_yap.gox)):

```go
static "/foo", FS("public")
static "/"
run ":8080"
```

Static files server also can use a `http.FileSystem` instead of `fs.FS` object (See [yapserve](https://github.com/xushiwei/yapserve) for details):

```go
import "github.com/qiniu/x/http/fs"

static "/", fs.http("https://goplus.org"), false // false means not allow to redirect
run ":8888"
```

### YAP Template

demo in Go ([blog.go](../demo/blog/blog.go), [article_yap.html](../demo/blog/yap/article_yap.html)):

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

demo in Go+ classfile ([blog_yap.gox](../demo/classfile_blog/blog_yap.gox), [article_yap.html](../demo/classfile_blog/yap/article_yap.html)):

```go
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

### YAP Test Framework

Suppose we have a web server named `foo` ([demo/foo/foo_yap.gox](../ytest/demo/foo/foo_yap.gox)):

```go
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

Then we create a yaptest file ([demo/foo/foo_ytest.gox](../ytest/demo/foo/foo_ytest.gox)):

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

You can change the directive `mock` to `testServer` (see [demo/foo/bar_ytest.gox](../ytest/demo/foo/bar_ytest.gox)), and keep everything else unchanged:

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

We can change this example more complicated:

```coffee
host "https://example.com", "http://localhost:8080"
testauth := oauth2("...")

run "urlWithVar", => {
	id := "123"
	get "https://example.com/p/${id}"
	ret
	echo "code:", resp.code
	echo "body:", resp.body
}

run "matchWithVar", => {
	code := Var(int)
	id := "123"
	get "https://example.com/p/${id}"
	ret code
	echo "code:", code
	match code, 200
}

run "postWithAuth", => {
	id := "123"
	title := "title"
	author := "author"
	post "https://example.com/p/${id}"
	auth testauth
	json {
		"title":  title,
		"author": author,
	}
	ret 200 # match resp.code, 200
	echo "body:", resp.body
}

run "matchJsonObject", => {
	title := Var(string)
	author := Var(string)
	id := "123"
	get "https://example.com/p/${id}"
	ret 200
	json {
		"title":  title,
		"author": author,
	}
	echo "title:", title
	echo "author:", author
}
```

For more details, see [yaptest - Go+ HTTP Test Framework](../ytest).
