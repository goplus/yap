yap - Yet Another Go/XGo HTTP Web Framework
======

[![Build Status](https://github.com/goplus/yap/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/yap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/yap)](https://goreportcard.com/report/github.com/goplus/yap)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/yap.svg?label=release)](https://github.com/goplus/yap/releases)
[![Coverage Status](https://codecov.io/gh/goplus/yap/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/yap)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/yap.svg)](https://pkg.go.dev/github.com/goplus/yap)
[![Language](https://img.shields.io/badge/language-XGo-blue.svg)](https://github.com/goplus/xgo)

This repo contains three [XGo classfiles](https://github.com/goplus/xgo/blob/main/doc/classfile.md). They are [yap](#yap-http-web-framework) (a HTTP Web Framework), [yaptest](ytest) (a HTTP Test Framework) and [ydb](ydb) (a XGo Database Framework).

The classfile [yap](#yap-http-web-framework) has the file suffix `_yap.gox`. The classfile [yaptest](ytest) has the file suffix `_ytest.gox`. And the classfile [ydb](ydb) has the file suffix `_ydb.gox`.

Before using `yap`, `yaptest` or `ydb`, you need to add `github.com/goplus/yap` to `go.mod`:

```sh
xgo get github.com/goplus/yap@latest
```

For more details, see [YAP Framework Manual](doc/manual.md).


### How to use in XGo

First let us initialize a hello project:

```sh
xgo mod init hello
```

Then we have it reference a classfile called `yap` as the HTTP Web Framework:

```sh
xgo get github.com/goplus/yap@latest
```

Create a file named [get_yap.gox](demo/classfile2_hello/get_yap.gox) with the following content:

```go
html `<html><body>Hello, YAP!</body></html>`
```

Execute the following commands:

```sh
xgo mod tidy
xgo run .
```

A simplest web program is running now. At this time, if you visit http://localhost:8080, you will get:

```
Hello, YAP!
```


### yap: HTTP Web Framework

This classfile has the file suffix `_yap.gox`.


#### Router and Parameters

YAP uses filenames to define routes. `get_yap.gox`'s route is `get "/"` (GET homepage), and `get_p_#id_yap.gox`'s route is `get "/p/:id"` (In fact, the filename can also be `get_p_:id_yap.gox`, but it is not recommended because `:` is not allowed to exist in filenames under Windows).

Let's create a file named [get_p_#id_yap.gox](demo/classfile2_hello/get_p_%23id_yap.gox) with the following content:

```coffee
json {
	"id": ${id},
}
```

Execute `xgo run .` and visit http://localhost:8080/p/123, you will get:

```
{"id": "123"}
```


#### YAP Template

In most cases, we don't use the `html` directive to generate html pages, but use the `yap` template engine. See [get_p_#id_yap.gox](demo/classfile2_blog/get_p_%23id_yap.gox):

```coffee
yap "article", {
	"id": ${id},
}
```

It means finding a template called `article` to render. See [yap/article_yap.html](demo/classfile2_blog/yap/article_yap.html):

```html
<html>
<head><meta charset="utf-8"/></head>
<body>Article {{.id}}</body>
</html>
```

#### Run at specified address

By default the YAP server runs on `localhost:8080`, but you can change it in [main_yap.gox](demo/classfile2_blog/main_yap.gox) file:

```coffee
run ":8888"
```


#### Static files

Static files server demo ([main_yap.gox](demo/classfile2_static/main_yap.gox)):

```coffee
static "/foo", FS("public")
static "/"

run ":8080"
```


### yaptest: HTTP Test Framework

[yaptest](ytest) is a web server testing framework. This classfile has the file suffix `_ytest.gox`.

Suppose we have a web server ([foo/get_p_#id_yap.gox](ytest/demo/foo/get_p_%23id_yap.gox)):

```go
json {
	"id": ${id},
}
```

Then we create a yaptest file ([foo/foo_ytest.gox](ytest/demo/foo/foo_ytest.gox)):

```go
mock "foo.com", new(AppV2)  // name of any YAP v2 web server is `AppV2`

id := "123"
get "http://foo.com/p/${id}"
ret 200
json {
	"id": id,
}
```

The directive `mock` creates the web server by [mockhttp](https://pkg.go.dev/github.com/qiniu/x/mockhttp). Then we write test code directly.

You can change the directive `mock` to `testServer` (see [foo/bar_ytest.gox](ytest/demo/foo/bar_ytest.gox)), and keep everything else unchanged:

```go
testServer "foo.com", new(AppV2)

id := "123"
get "http://foo.com/p/${id}"
ret 200
json {
	"id": id,
}
```

The directive `testServer` creates the web server by [net/http/httptest](https://pkg.go.dev/net/http/httptest#NewServer) and obtained a random port as the service address. Then it calls the directive [host](https://pkg.go.dev/github.com/goplus/yap/ytest#App.Host) to map the random service address to `foo.com`. This makes all other code no need to changed.

For more details, see [yaptest - XGo HTTP Test Framework](ytest).


### ydb: Database Framework

This classfile has the file suffix `_ydb.gox`.

TODO
