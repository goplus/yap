/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package yap

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/qiniu/x/http/fsx"
)

const (
	GopPackage = true
)

// App is project class of YAP classfile (old version).
type App struct {
	Engine
}

// Get is a shortcut for router.Route(http.MethodGet, path, handle)
func (p *App) Get(path string, handle func(ctx *Context)) {
	p.Route(http.MethodGet, path, handle)
}

// Head is a shortcut for router.Route(http.MethodHead, path, handle)
func (p *App) Head(path string, handle func(ctx *Context)) {
	p.Route(http.MethodHead, path, handle)
}

// Options is a shortcut for router.Route(http.MethodOptions, path, handle)
func (p *App) Options(path string, handle func(ctx *Context)) {
	p.Route(http.MethodOptions, path, handle)
}

// Post is a shortcut for router.Route(http.MethodPost, path, handle)
func (p *App) Post(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPost, path, handle)
}

// Put is a shortcut for router.Route(http.MethodPut, path, handle)
func (p *App) Put(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPut, path, handle)
}

// Patch is a shortcut for router.Route(http.MethodPatch, path, handle)
func (p *App) Patch(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPatch, path, handle)
}

// Delete is a shortcut for router.Route(http.MethodDelete, path, handle)
func (p *App) Delete(path string, handle func(ctx *Context)) {
	p.Route(http.MethodDelete, path, handle)
}

// Static serves static files from a dir (default is "$YapFS/static").
func (p *App) Static__0(pattern string, dir ...fs.FS) {
	p.Static(pattern, dir...)
}

// Static serves static files from a http file system scheme (url).
// See https://pkg.go.dev/github.com/qiniu/x/http/fsx for more information.
func (p *App) Static__1(pattern string, ctx context.Context, url string) (closer fsx.Closer, err error) {
	fs, closer, err := fsx.Open(ctx, url)
	if err == nil {
		p.StaticHttp(pattern, fs, false)
	}
	return
}

// Static serves static files from a http file system.
func (p *App) Static__2(pattern string, fs http.FileSystem, allowRedirect ...bool) {
	p.StaticHttp(pattern, fs, allowRedirect...)
}

// AppType represents an abstract of YAP applications.
type AppType interface {
	InitYap(fs ...fs.FS)
	SetLAS(listenAndServe func(addr string, handler http.Handler) error)
	Route(method, path string, handle func(ctx *Context))
	Handle(pattern string, f func(ctx *Context))
	Run(addr string, mws ...func(h http.Handler) http.Handler) error
}

var (
	_ AppType = (*App)(nil)
	_ AppType = (*Engine)(nil)
)

// Gopt_App_Main is required by Go+ compiler as the entry of a YAP project.
func Gopt_App_Main(app AppType) {
	app.InitYap()
	app.(interface{ MainEntry() }).MainEntry()
}

const (
	mimeText   = "text/plain"
	mimeHtml   = "text/html"
	mimeBinary = "application/octet-stream"
)

func (p *Context) Text__0(code int, mime string, text string) {
	p.TEXT(code, mime, text)
}

func (p *Context) Text__1(code int, text string) {
	p.TEXT(code, mimeText, text)
}

func (p *Context) Text__2(text string) {
	p.TEXT(200, mimeText, text)
}

func (p *Context) Text__3(code int, text []byte) {
	p.DATA(code, mimeText, text)
}

func (p *Context) Text__4(text []byte) {
	p.DATA(200, mimeText, text)
}

func (p *Context) Binary__0(code int, mime string, data []byte) {
	p.DATA(code, mime, data)
}

func (p *Context) Binary__1(code int, data []byte) {
	p.DATA(code, mimeBinary, data)
}

func (p *Context) Binary__2(data []byte) {
	p.DATA(200, mimeBinary, data)
}

func (p *Context) Binary__3(code int, data string) {
	p.TEXT(code, mimeBinary, data)
}

func (p *Context) Binary__4(data string) {
	p.TEXT(200, mimeBinary, data)
}

func (p *Context) Html__0(code int, text string) {
	p.TEXT(code, mimeHtml, text)
}

func (p *Context) Html__1(text string) {
	p.TEXT(200, mimeHtml, text)
}

func (p *Context) Html__2(code int, text []byte) {
	p.DATA(code, mimeHtml, text)
}

func (p *Context) Html__3(text []byte) {
	p.DATA(200, mimeHtml, text)
}

func (p *Context) Json__0(code int, data interface{}) {
	p.JSON(code, data)
}

func (p *Context) Json__1(data interface{}) {
	p.JSON(200, data)
}

func (p *Context) PrettyJson__0(code int, data interface{}) {
	p.PrettyJSON(code, data)
}

func (p *Context) PrettyJson__1(data interface{}) {
	p.PrettyJSON(200, data)
}

func (p *Context) Yap__0(code int, yapFile string, data interface{}) {
	p.YAP(code, yapFile, data)
}

func (p *Context) Yap__1(yapFile string, data interface{}) {
	p.YAP(200, yapFile, data)
}
