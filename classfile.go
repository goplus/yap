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
	"io/fs"
	"net/http"
	"os"
)

const (
	GopPackage = true
)

type App struct {
	*Engine
}

func (p *App) initApp() {
	p.Engine = New()
}

// Get is a shortcut for router.Route(http.MethodGet, path, handle)
func (p App) Get(path string, handle func(ctx *Context)) {
	p.Route(http.MethodGet, path, handle)
}

// Head is a shortcut for router.Route(http.MethodHead, path, handle)
func (p App) Head(path string, handle func(ctx *Context)) {
	p.Route(http.MethodHead, path, handle)
}

// Options is a shortcut for router.Route(http.MethodOptions, path, handle)
func (p App) Options(path string, handle func(ctx *Context)) {
	p.Route(http.MethodOptions, path, handle)
}

// Post is a shortcut for router.Route(http.MethodPost, path, handle)
func (p App) Post(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPost, path, handle)
}

// Put is a shortcut for router.Route(http.MethodPut, path, handle)
func (p App) Put(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPut, path, handle)
}

// Patch is a shortcut for router.Route(http.MethodPatch, path, handle)
func (p App) Patch(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPatch, path, handle)
}

// Delete is a shortcut for router.Route(http.MethodDelete, path, handle)
func (p App) Delete(path string, handle func(ctx *Context)) {
	p.Route(http.MethodDelete, path, handle)
}

// Run with specified `fsys` as yap template directory.
func (p App) Run__0(fsys fs.FS, addr string, mws ...func(h http.Handler) http.Handler) {
	p.initYapFS(fsys)
	p.Run(addr, mws...)
}

// Run with cwd as yap template directory.
func (p App) Run__1(addr string, mws ...func(h http.Handler) http.Handler) {
	p.Run__0(os.DirFS("."), addr, mws...)
}

// Gopt_App_Main is required by Go+ compiler as the entry of a YAP project.
func Gopt_App_Main(app interface{ initApp() }) {
	app.initApp()
	app.(interface{ MainEntry() }).MainEntry()
}

const (
	mimeText   = "text/plain"
	mimeHtml   = "text/html"
	mimeBinary = "application/octet-stream"
)

func (p *Context) Text__0(text string) {
	p.TEXT(200, mimeText, text)
}

func (p *Context) Text__1(text []byte) {
	p.DATA(200, mimeText, text)
}

func (p *Context) Text__2(code int, mime string, text string) {
	p.TEXT(code, mime, text)
}

func (p *Context) Text__3(code int, text string) {
	p.TEXT(code, mimeText, text)
}

func (p *Context) Text__4(code int, text []byte) {
	p.DATA(code, mimeText, text)
}

func (p *Context) Binary__0(data []byte) {
	p.DATA(200, mimeBinary, data)
}

func (p *Context) Binary__1(data string) {
	p.TEXT(200, mimeBinary, data)
}

func (p *Context) Binary__2(code int, mime string, data []byte) {
	p.DATA(code, mime, data)
}

func (p *Context) Binary__3(code int, data []byte) {
	p.DATA(code, mimeBinary, data)
}

func (p *Context) Binary__4(code int, data string) {
	p.TEXT(code, mimeBinary, data)
}

func (p *Context) Html__0(text string) {
	p.TEXT(200, mimeHtml, text)
}

func (p *Context) Html__1(text []byte) {
	p.DATA(200, mimeHtml, text)
}

func (p *Context) Html__2(code int, text string) {
	p.TEXT(code, mimeHtml, text)
}

func (p *Context) Html__3(code int, text []byte) {
	p.DATA(code, mimeHtml, text)
}

func (p *Context) Json__0(data interface{}) {
	p.JSON(200, data)
}

func (p *Context) Json__1(code int, data interface{}) {
	p.JSON(code, data)
}

func (p *Context) PrettyJson__0(data interface{}) {
	p.PrettyJSON(200, data)
}

func (p *Context) PrettyJson__1(code int, data interface{}) {
	p.PrettyJSON(code, data)
}

func (p *Context) Yap__0(yapFile string, data interface{}) {
	p.YAP(200, yapFile, data)
}

func (p *Context) Yap__1(code int, yapFile string, data interface{}) {
	p.YAP(code, yapFile, data)
}
