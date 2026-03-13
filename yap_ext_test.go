/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package yap_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/goplus/yap"
)

// newEngine creates a new Engine with the workspace filesystem.
func newEngine() *yap.Engine {
	return yap.New(os.DirFS("."))
}

// newContext creates a Context for testing from an httptest recorder and request.
func newContext(method, path string, body io.Reader) (*yap.Engine, *httptest.ResponseRecorder, *yap.Context) {
	e := newEngine()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)
	ctx := e.NewContext(w, req)
	return e, w, ctx
}

// --- Context tests ---

func TestContextParam(t *testing.T) {
	_, w, ctx := newContext("GET", "/p/hello?name=world", nil)
	_ = w
	if got := ctx.Param("name"); got != "world" {
		t.Fatalf("Param: expected world, got %s", got)
	}
}

func TestContextParamInt(t *testing.T) {
	_, _, ctx := newContext("GET", "/?id=42", nil)
	if got := ctx.ParamInt("id", 0); got != 42 {
		t.Fatalf("ParamInt: expected 42, got %d", got)
	}
}

func TestContextParamIntDefault(t *testing.T) {
	_, _, ctx := newContext("GET", "/", nil)
	if got := ctx.ParamInt("id", 99); got != 99 {
		t.Fatalf("ParamInt default: expected 99, got %d", got)
	}
}

func TestContextParamIntInvalid(t *testing.T) {
	_, _, ctx := newContext("GET", "/?id=notanint", nil)
	if got := ctx.ParamInt("id", 5); got != 5 {
		t.Fatalf("ParamInt invalid: expected 5, got %d", got)
	}
}

func TestContextXGoEnv(t *testing.T) {
	_, _, ctx := newContext("GET", "/?key=val", nil)
	if got := ctx.XGo_Env("key"); got != "val" {
		t.Fatalf("XGo_Env: expected val, got %s", got)
	}
}

func TestContextAcceptMatch(t *testing.T) {
	_, _, ctx := newContext("GET", "/", nil)
	ctx.Request.Header.Set("Accept", "text/html, application/json")
	if got := ctx.Accept("text/html", "text/plain"); got != "text/html" {
		t.Fatalf("Accept: expected text/html, got %s", got)
	}
}

func TestContextAcceptNoMatch(t *testing.T) {
	_, _, ctx := newContext("GET", "/", nil)
	ctx.Request.Header.Set("Accept", "image/png")
	if got := ctx.Accept("text/html", "application/json"); got != "" {
		t.Fatalf("Accept no match: expected empty, got %s", got)
	}
}

func TestContextAcceptWithQuality(t *testing.T) {
	_, _, ctx := newContext("GET", "/", nil)
	ctx.Request.Header.Set("Accept", "text/html;q=0.9, application/json")
	if got := ctx.Accept("text/html"); got != "text/html" {
		t.Fatalf("Accept with quality: expected text/html, got %s", got)
	}
}

func TestContextAcceptEmpty(t *testing.T) {
	_, _, ctx := newContext("GET", "/", nil)
	if got := ctx.Accept("text/html"); got != "" {
		t.Fatalf("Accept empty header: expected empty, got %s", got)
	}
}

func TestContextRedirectDefault(t *testing.T) {
	_, w, ctx := newContext("GET", "/old", nil)
	ctx.Redirect("/new")
	if w.Code != http.StatusFound {
		t.Fatalf("Redirect: expected %d, got %d", http.StatusFound, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/new" {
		t.Fatalf("Redirect location: expected /new, got %s", loc)
	}
}

func TestContextRedirectCustomCode(t *testing.T) {
	_, w, ctx := newContext("GET", "/old", nil)
	ctx.Redirect("/new", http.StatusMovedPermanently)
	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("Redirect custom code: expected %d, got %d", http.StatusMovedPermanently, w.Code)
	}
}

func TestContextTEXT(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.TEXT(200, "text/plain", "hello world")
	if w.Code != 200 {
		t.Fatalf("TEXT: expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("TEXT content-type: expected text/plain, got %s", ct)
	}
	if body := w.Body.String(); body != "hello world" {
		t.Fatalf("TEXT body: expected 'hello world', got %s", body)
	}
}

func TestContextDATA(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	data := []byte{0x01, 0x02, 0x03}
	ctx.DATA(200, "application/octet-stream", data)
	if !bytes.Equal(w.Body.Bytes(), data) {
		t.Fatalf("DATA body mismatch")
	}
}

func TestContextJSON(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.JSON(200, map[string]string{"key": "value"})
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("JSON content-type: expected application/json, got %s", ct)
	}
	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal: %v", err)
	}
	if result["key"] != "value" {
		t.Fatalf("JSON value mismatch")
	}
}

func TestContextPrettyJSON(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.PrettyJSON(200, map[string]string{"key": "value"})
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("PrettyJSON content-type: expected application/json, got %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "\n") {
		t.Fatalf("PrettyJSON: expected indented output, got %s", body)
	}
}

func TestContextSTREAM(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	reader := strings.NewReader("stream data")
	ctx.STREAM(200, "text/plain", reader, nil)
	if w.Code != 200 {
		t.Fatalf("STREAM: expected 200, got %d", w.Code)
	}
	if body := w.Body.String(); body != "stream data" {
		t.Fatalf("STREAM body: expected 'stream data', got %s", body)
	}
}

func TestContextSTREAMNoMime(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	reader := strings.NewReader("raw data")
	ctx.STREAM(200, "", reader, nil)
	if body := w.Body.String(); body != "raw data" {
		t.Fatalf("STREAM no mime body: expected 'raw data', got %s", body)
	}
	_ = w
}

func TestContextSTREAMWithSmallBuffer(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	reader := strings.NewReader("buffered data")
	buf := make([]byte, 4)
	ctx.STREAM(200, "text/plain", reader, buf)
	if body := w.Body.String(); body != "buffered data" {
		t.Fatalf("STREAM with buffer: expected 'buffered data', got %s", body)
	}
	_ = w
}

func TestContextSTREAMLargeBuffer(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	reader := strings.NewReader("large buf data")
	buf := make([]byte, 64*1024)
	ctx.STREAM(200, "text/plain", reader, buf)
	if body := w.Body.String(); body != "large buf data" {
		t.Fatalf("STREAM large buffer body mismatch")
	}
	_ = w
}

// --- Engine tests ---

func TestEngineHandle(t *testing.T) {
	e := newEngine()
	e.Handle("/hello", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "hi")
	})
	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	e.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Handle: expected 200, got %d", w.Code)
	}
	if body := w.Body.String(); body != "hi" {
		t.Fatalf("Handle body: expected 'hi', got %s", body)
	}
}

func TestEngineHandlerWithMiddleware(t *testing.T) {
	e := newEngine()
	called := false
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}
	h := e.Handler(mw)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if !called {
		t.Fatal("Handler middleware: middleware was not called")
	}
}

func TestEngineSetDelims(t *testing.T) {
	e := newEngine()
	e.SetDelims("[[", "]]")
	// Just test it doesn't panic — template parsing would need actual template files
}

func TestEngineStaticHttp(t *testing.T) {
	e := newEngine()
	fsys := os.DirFS(".")
	e.StaticHttp("/static/", http.FS(fsys))
	req := httptest.NewRequest("GET", "/static/go.mod", nil)
	w := httptest.NewRecorder()
	e.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("StaticHttp: expected 200, got %d", w.Code)
	}
}

func TestEngineStaticHttpNoRedirect(t *testing.T) {
	e := newEngine()
	fsys := os.DirFS(".")
	e.StaticHttp("/files/", http.FS(fsys), false)
	req := httptest.NewRequest("GET", "/files/go.mod", nil)
	w := httptest.NewRecorder()
	e.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("StaticHttp no redirect: expected 200, got %d", w.Code)
	}
}

func TestEngineStaticHttpPatternNoSlash(t *testing.T) {
	e := newEngine()
	fsys := os.DirFS(".")
	// pattern without trailing slash — StaticHttp should append it
	e.StaticHttp("/assets", http.FS(fsys))
	req := httptest.NewRequest("GET", "/assets/go.mod", nil)
	w := httptest.NewRecorder()
	e.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("StaticHttp pattern no slash: expected 200, got %d", w.Code)
	}
}

func TestEngineStatic(t *testing.T) {
	e := yap.New(os.DirFS("."))
	// "static" directory doesn't exist in workspace but Static should not panic on registration
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Static panicked: %v", r)
		}
	}()
	// Use an explicit fs to avoid panic on missing "static" directory
	e.Static("/pub/", os.DirFS("."))
	req := httptest.NewRequest("GET", "/pub/go.mod", nil)
	w := httptest.NewRecorder()
	e.Mux.ServeHTTP(w, req)
}

func TestEngineSubFS(t *testing.T) {
	// internal package test via SubFS
	fsys := os.DirFS(".")
	sub := yap.SubFS(fsys, "ytest")
	if sub == nil {
		t.Fatal("SubFS returned nil")
	}
}

func TestEngineSubFSPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("SubFS with bad dir should panic")
		}
	}()
	yap.SubFS(os.DirFS("."), "nonexistent_dir_xyz")
}

// --- Router method tests ---

func TestRouterHEAD(t *testing.T) {
	e := newEngine()
	e.GET("/head-test", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "head content")
	})
	req := httptest.NewRequest("HEAD", "/head-test", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("HEAD: expected 200, got %d", w.Code)
	}
	// HEAD should have no body
	if w.Body.Len() != 0 {
		t.Fatalf("HEAD: body should be empty, got %s", w.Body.String())
	}
}

func TestRouterOPTIONS(t *testing.T) {
	e := newEngine()
	e.GET("/opt", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	req := httptest.NewRequest("OPTIONS", "/opt", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	allow := w.Header().Get("Allow")
	if allow == "" {
		t.Fatal("OPTIONS: expected Allow header to be set")
	}
}

func TestRouterGlobalOPTIONS(t *testing.T) {
	e := newEngine()
	e.GET("/resource", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	e.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "options")
		w.WriteHeader(204)
	})
	req := httptest.NewRequest("OPTIONS", "/resource", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Header().Get("X-Custom") != "options" {
		t.Fatal("GlobalOPTIONS: expected custom handler to be called")
	}
}

func TestRouterPOST(t *testing.T) {
	e := newEngine()
	e.POST("/submit", func(ctx *yap.Context) {
		ctx.TEXT(201, "text/plain", "created")
	})
	req := httptest.NewRequest("POST", "/submit", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("POST: expected 201, got %d", w.Code)
	}
}

func TestRouterPUT(t *testing.T) {
	e := newEngine()
	e.PUT("/item/:id", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", ctx.Param("id"))
	})
	req := httptest.NewRequest("PUT", "/item/42", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("PUT: expected 200, got %d", w.Code)
	}
}

func TestRouterPATCH(t *testing.T) {
	e := newEngine()
	e.PATCH("/update/:id", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "patched")
	})
	req := httptest.NewRequest("PATCH", "/update/1", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("PATCH: expected 200, got %d", w.Code)
	}
}

func TestRouterDELETE(t *testing.T) {
	e := newEngine()
	e.DELETE("/item/:id", func(ctx *yap.Context) {
		ctx.TEXT(204, "text/plain", "")
	})
	req := httptest.NewRequest("DELETE", "/item/5", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("DELETE: expected 204, got %d", w.Code)
	}
}

func TestRouterOPTIONSMethod(t *testing.T) {
	e := newEngine()
	e.OPTIONS("/opts", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "options")
	})
	req := httptest.NewRequest("OPTIONS", "/opts", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("OPTIONS route: expected 200, got %d", w.Code)
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	e := newEngine()
	e.GET("/only-get", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	req := httptest.NewRequest("POST", "/only-get", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("MethodNotAllowed: expected 405, got %d", w.Code)
	}
}

func TestRouterCustomMethodNotAllowed(t *testing.T) {
	e := newEngine()
	e.GET("/only-get", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	e.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	})
	req := httptest.NewRequest("POST", "/only-get", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 418 {
		t.Fatalf("CustomMethodNotAllowed: expected 418, got %d", w.Code)
	}
}

func TestRouterRedirectTrailingSlash(t *testing.T) {
	e := newEngine()
	e.GET("/trailingslash/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	req := httptest.NewRequest("GET", "/trailingslash", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("RedirectTrailingSlash: expected 301, got %d", w.Code)
	}
}

func TestRouterRedirectFixedPath(t *testing.T) {
	e := newEngine()
	e.GET("/fixed", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	req := httptest.NewRequest("GET", "/FIXED", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("RedirectFixedPath: expected 301, got %d", w.Code)
	}
}

func TestRouterPanicHandler(t *testing.T) {
	e := newEngine()
	e.GET("/panic", func(ctx *yap.Context) {
		panic("test panic")
	})
	recovered := false
	e.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv any) {
		recovered = true
		w.WriteHeader(500)
	}
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if !recovered {
		t.Fatal("PanicHandler was not called")
	}
	if w.Code != 500 {
		t.Fatalf("PanicHandler: expected 500, got %d", w.Code)
	}
}

func TestRouterNotFound(t *testing.T) {
	e := newEngine()
	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Fatalf("NotFound: expected 404, got %d", w.Code)
	}
}

func TestRouterPostRedirectTrailingSlash(t *testing.T) {
	e := newEngine()
	e.POST("/submit/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "ok")
	})
	req := httptest.NewRequest("POST", "/submit", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	// POST should use 308 Permanent Redirect
	if w.Code != http.StatusPermanentRedirect {
		t.Fatalf("POST TrailingSlash redirect: expected 308, got %d", w.Code)
	}
}

// --- App classfile method tests ---

func TestAppGet(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Get("/app-get", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "app-get")
	})
	req := httptest.NewRequest("GET", "/app-get", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Get: expected 200, got %d", w.Code)
	}
}

func TestAppHead(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Head("/app-head", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "app-head")
	})
	req := httptest.NewRequest("HEAD", "/app-head", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Head: expected 200, got %d", w.Code)
	}
}

func TestAppOptions(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Options("/app-options", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "app-options")
	})
	req := httptest.NewRequest("OPTIONS", "/app-options", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Options: expected 200, got %d", w.Code)
	}
}

func TestAppPost(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Post("/app-post", func(ctx *yap.Context) {
		ctx.TEXT(201, "text/plain", "app-post")
	})
	req := httptest.NewRequest("POST", "/app-post", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("App.Post: expected 201, got %d", w.Code)
	}
}

func TestAppPut(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Put("/app-put/:id", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", ctx.Param("id"))
	})
	req := httptest.NewRequest("PUT", "/app-put/7", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Put: expected 200, got %d", w.Code)
	}
}

func TestAppPatch(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Patch("/app-patch/:id", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/plain", "patched")
	})
	req := httptest.NewRequest("PATCH", "/app-patch/3", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Patch: expected 200, got %d", w.Code)
	}
}

func TestAppDelete(t *testing.T) {
	a := new(yap.App)
	a.InitYap()
	a.Delete("/app-delete/:id", func(ctx *yap.Context) {
		ctx.TEXT(204, "text/plain", "")
	})
	req := httptest.NewRequest("DELETE", "/app-delete/9", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("App.Delete: expected 204, got %d", w.Code)
	}
}

func TestAppStaticFS(t *testing.T) {
	a := new(yap.App)
	a.InitYap(os.DirFS("."))
	a.Static__0("/static/", os.DirFS("."))
	req := httptest.NewRequest("GET", "/static/go.mod", nil)
	w := httptest.NewRecorder()
	a.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Static__0: expected 200, got %d", w.Code)
	}
}

func TestAppStaticHTTP(t *testing.T) {
	a := new(yap.App)
	a.InitYap(os.DirFS("."))
	a.Static__2("/statichttp/", http.FS(os.DirFS(".")))
	req := httptest.NewRequest("GET", "/statichttp/go.mod", nil)
	w := httptest.NewRecorder()
	a.Mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("App.Static__2: expected 200, got %d", w.Code)
	}
}

// --- Context response shorthand method tests ---

func TestContextText0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Text__0(200, "text/plain", "text0")
	if w.Body.String() != "text0" {
		t.Fatalf("Text__0: expected 'text0', got %s", w.Body.String())
	}
}

func TestContextText1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Text__1(200, "text1")
	if w.Body.String() != "text1" {
		t.Fatalf("Text__1: expected 'text1', got %s", w.Body.String())
	}
}

func TestContextText2(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Text__2("text2")
	if w.Code != 200 {
		t.Fatalf("Text__2: expected 200, got %d", w.Code)
	}
	if w.Body.String() != "text2" {
		t.Fatalf("Text__2: expected 'text2', got %s", w.Body.String())
	}
}

func TestContextText3(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Text__3(200, []byte("text3"))
	if w.Body.String() != "text3" {
		t.Fatalf("Text__3: expected 'text3', got %s", w.Body.String())
	}
}

func TestContextText4(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Text__4([]byte("text4"))
	if w.Code != 200 {
		t.Fatalf("Text__4: expected 200, got %d", w.Code)
	}
}

func TestContextBinary0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Binary__0(200, "application/octet-stream", []byte{1, 2})
	if !bytes.Equal(w.Body.Bytes(), []byte{1, 2}) {
		t.Fatalf("Binary__0: body mismatch")
	}
}

func TestContextBinary1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Binary__1(200, []byte{3, 4})
	if !bytes.Equal(w.Body.Bytes(), []byte{3, 4}) {
		t.Fatalf("Binary__1: body mismatch")
	}
}

func TestContextBinary2(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Binary__2([]byte{5})
	if w.Code != 200 {
		t.Fatalf("Binary__2: expected 200, got %d", w.Code)
	}
}

func TestContextBinary3(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Binary__3(200, "bin3")
	if w.Body.String() != "bin3" {
		t.Fatalf("Binary__3: expected 'bin3', got %s", w.Body.String())
	}
}

func TestContextBinary4(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Binary__4("bin4")
	if w.Code != 200 {
		t.Fatalf("Binary__4: expected 200, got %d", w.Code)
	}
}

func TestContextHtml0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Html__0(200, "<p>html0</p>")
	if ct := w.Header().Get("Content-Type"); ct != "text/html" {
		t.Fatalf("Html__0: expected text/html, got %s", ct)
	}
}

func TestContextHtml1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Html__1("<p>html1</p>")
	if w.Code != 200 {
		t.Fatalf("Html__1: expected 200, got %d", w.Code)
	}
}

func TestContextHtml2(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Html__2(200, []byte("<p>html2</p>"))
	if w.Body.String() != "<p>html2</p>" {
		t.Fatalf("Html__2: body mismatch")
	}
}

func TestContextHtml3(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Html__3([]byte("<p>html3</p>"))
	if w.Code != 200 {
		t.Fatalf("Html__3: expected 200, got %d", w.Code)
	}
}

func TestContextJson0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Json__0(200, map[string]int{"n": 1})
	var result map[string]int
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["n"] != 1 {
		t.Fatalf("Json__0: value mismatch")
	}
}

func TestContextJson1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Json__1(map[string]int{"n": 2})
	if w.Code != 200 {
		t.Fatalf("Json__1: expected 200, got %d", w.Code)
	}
}

func TestContextPrettyJson0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.PrettyJson__0(200, map[string]int{"n": 3})
	body := w.Body.String()
	if !strings.Contains(body, "\"n\"") {
		t.Fatalf("PrettyJson__0: body mismatch, got %s", body)
	}
}

func TestContextPrettyJson1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.PrettyJson__1(map[string]int{"n": 4})
	if w.Code != 200 {
		t.Fatalf("PrettyJson__1: expected 200, got %d", w.Code)
	}
}

func TestContextStream0(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Stream__0(200, "text/plain", strings.NewReader("stream0"), nil)
	if w.Body.String() != "stream0" {
		t.Fatalf("Stream__0: expected 'stream0', got %s", w.Body.String())
	}
}

func TestContextStream1(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Stream__1(200, "text/plain", strings.NewReader("stream1"))
	if w.Body.String() != "stream1" {
		t.Fatalf("Stream__1: expected 'stream1', got %s", w.Body.String())
	}
}

func TestContextStream2(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Stream__2(strings.NewReader("stream2"), nil)
	if w.Code != 200 {
		t.Fatalf("Stream__2: expected 200, got %d", w.Code)
	}
}

func TestContextStream3(t *testing.T) {
	_, w, ctx := newContext("GET", "/", nil)
	ctx.Stream__3(strings.NewReader("stream3"))
	if w.Code != 200 {
		t.Fatalf("Stream__3: expected 200, got %d", w.Code)
	}
}

// --- Engine FS test ---

func TestEngineFS(t *testing.T) {
	e := yap.New(os.DirFS("."))
	sub := e.FS("ytest")
	if sub == nil {
		t.Fatal("Engine.FS returned nil")
	}
}

// --- InitYap idempotent ---

func TestInitYapIdempotent(t *testing.T) {
	e := newEngine()
	e.InitYap() // calling twice should be a no-op
	e.InitYap(os.DirFS("."))
}
