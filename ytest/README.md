yaptest - Go+ HTTP Test Framework
=====

yaptest is a web server testing framework. This classfile has the file suffix `_ytest.gox`.

Before using `yaptest`, you need to add `github.com/goplus/yap` to `go.mod`:

```
gop get github.com/goplus/yap@latest
```

Suppose we have a web server ([foo/get_p_#id.yap](demo/foo/get_p_%23id.yap)):

```go
json {
	"id": ${id},
}
```

Then we create a yaptest file ([foo/foo_ytest.gox](demo/foo/foo_ytest.gox)):

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

You can change the directive `mock` to `testServer` (see [foo/bar_ytest.gox](demo/foo/bar_ytest.gox)), and keep everything else unchanged:

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


### match

This is almost the core concept in `yaptest`. It matches two objects.

Let’s look at [a simple example](demo/match/simple/simple_yapt.gox) first:

```go
id := Var(int)
match id, 1+2
echo id
```

Here we define a variable called `id` and match it with expression `1+2`. If the variable is unbound, it is assigned the value of the expression. In this way the value of `id` becomes `3`.

So far, you've seen `match` like the assignment side. But you cannot assign a different value to a variable that has been bound:

```go
id := Var(int)
match id, 1+2
match id, 3
echo id

match id, 5  // unmatched value - expected: 3, got: 5
```

In the second `match` statement, the variable `id` has been bound. At this time, it will be compared with the expression value. If it is equal, it will succeed, otherwise an error will be reported (such as the third `match` statement above).

The `match` statement [can be complex](demo/match/complex/complex_yapt.gox), such as:

```go
d := Var(string)

match {
    "c": {"d": d},
}, {
    "a": 1,
    "b": 3.14,
    "c": {"d": "hello", "e": "world"},
    "f": 1,
}

echo d
match d, "hello"
```

Generally, the syntax of the match command is:

```go
match <ExpectedObject> <SourceObject>
```

Unbound variables are allowed in `<ExpectedObject>`, but cannot appear in `<SourceObject>`. `<ExpectedObject>` and `<SourceObject>` do not have to be exactly the same, but what appears in `<ExpectedObject>` must also appear in `<SourceObject>`. That is, it is required to be a subset relationship (`<ExpectedObject>` is a subset of `<SourceObject>`).

If a variable in `<ExpectedObject>` has not been bound, it will be bound according to the value of the corresponding `<SourceObject>`; if the variable has been bound, the values on both sides must match.

The cornerstone of `yaptest` is matching grammar. Let's look at the next example you saw at the beginning:

```go
ret 200
json {
	"id": id,
}
```
