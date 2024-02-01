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
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/goplus/yap/noredirect"
)

type H map[string]interface{}

type Engine struct {
	router
	Mux *http.ServeMux

	tpls   map[string]Template
	fs     fs.FS
	las    func(addr string, handler http.Handler) error
	delims Delims
}

// New creates a YAP engine.
func New(fs ...fs.FS) *Engine {
	e := new(Engine)
	e.InitYap(fs...)
	e.SetDelims("{{", "}}")
	return e
}

// InitYap initialize a YAP application.
func (p *Engine) InitYap(fs ...fs.FS) {
	if p.Mux == nil {
		p.Mux = http.NewServeMux()
		p.las = http.ListenAndServe
		p.router.init()
	}
	if fs != nil {
		p.initYapFS(fs[0])
	}
}

func (p *Engine) initYapFS(fsys fs.FS) {
	const name = "yap"
	if f, e := fsys.Open(name); e == nil {
		f.Close()
		if sub, e := fs.Sub(fsys, name); e == nil {
			fsys = sub
		}
	}
	p.fs = fsys
	p.tpls = make(map[string]Template)
}

func (p *Engine) SetDelims(left, right string) {
	if !(len(left) == 2 && len(right) == 2) {
		log.Panicln("The length of the delimiter must be two")
	}
	p.delims = Delims{left, right}
}

func (p *Engine) yapFS() fs.FS {
	if p.fs == nil {
		p.initYapFS(os.DirFS("."))
	}
	return p.fs
}

func (p *Engine) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{ResponseWriter: w, Request: r, engine: p}
	return ctx
}

// ServeHTTP makes the router implement the http.Handler interface.
func (p *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.router.serveHTTP(w, req, p)
}

// FS returns a $YapFS sub filesystem by specified a dir.
func (p *Engine) FS(dir string) (ret fs.FS) {
	return SubFS(p.yapFS(), dir)
}

// Static serves static files from a dir (default is "$YapFS/static").
func (p *Engine) Static(pattern string, dir ...fs.FS) {
	var fsys fs.FS
	if dir != nil {
		fsys = dir[0]
	} else {
		fsys = p.FS("static")
	}
	p.StaticHttp(pattern, http.FS(fsys))
}

// StaticHttp serves static files from fsys (http.FileSystem).
func (p *Engine) StaticHttp(pattern string, fsys http.FileSystem, allowRedirect ...bool) {
	if !strings.HasSuffix(pattern, "/") {
		pattern += "/"
	}
	allow := true
	if allowRedirect != nil {
		allow = allowRedirect[0]
	}
	var server http.Handler
	if allow {
		server = http.FileServer(fsys)
	} else {
		server = noredirect.FileServer(fsys)
	}
	p.Mux.Handle(pattern, http.StripPrefix(pattern, server))
}

// Handle registers the handler function for the given pattern.
func (p *Engine) Handle(pattern string, f func(ctx *Context)) {
	p.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		f(p.NewContext(w, r))
	})
}

// Handler returns the main entry that responds to HTTP requests.
func (p *Engine) Handler(mws ...func(h http.Handler) http.Handler) http.Handler {
	h := http.Handler(p)
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// Run listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (p *Engine) Run(addr string, mws ...func(h http.Handler) http.Handler) error {
	h := p.Handler(mws...)
	return p.las(addr, h)
}

// SetLAS sets listenAndServe func to listens on the TCP network address addr
// and to handle requests on incoming connections.
func (p *Engine) SetLAS(listenAndServe func(addr string, handler http.Handler) error) {
	p.las = listenAndServe
}

func (p *Engine) templ(path string) (t Template, err error) {
	fsys := p.yapFS()
	if p.tpls == nil {
		return Template{}, os.ErrNotExist
	}
	t, ok := p.tpls[path]
	if !ok {
		t, err = ParseFSFile(fsys, path+"_yap.html", p.delims)
		if err != nil {
			return
		}
		p.tpls[path] = t
	}
	return
}

// SubFS returns a sub filesystem by specified a dir.
func SubFS(fsys fs.FS, dir string) (ret fs.FS) {
	f, err := fsys.Open(dir)
	if err == nil {
		f.Close()
		ret, err = fs.Sub(fsys, dir)
	}
	if err != nil {
		log.Panicln("Get $YapFS sub filesystem failed:", err)
	}
	return ret
}
