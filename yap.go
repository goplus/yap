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

type H map[string]interface{}

type Engine struct {
	router
	Mux *http.ServeMux

	tpls map[string]Template
	fs   fs.FS
}

func New(fs ...fs.FS) *Engine {
	e := &Engine{
		Mux: http.NewServeMux(),
	}
	e.router.init()
	if fs != nil {
		e.initYapFS(fs[0])
	}
	return e
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

func (p *Engine) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{ResponseWriter: w, Request: r, engine: p}
	return ctx
}

// ServeHTTP makes the router implement the http.Handler interface.
func (p *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.router.serveHTTP(w, req, p)
}

func (p *Engine) Handle(pattern string, f func(ctx *Context)) {
	p.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		f(p.NewContext(w, r))
	})
}

func (p *Engine) Run(addr string, mws ...func(h http.Handler) http.Handler) {
	h := http.Handler(p)
	for _, mw := range mws {
		h = mw(h)
	}
	http.ListenAndServe(addr, h)
}

func (p *Engine) templ(path string) (t Template, err error) {
	if p.tpls == nil {
		return Template{}, os.ErrNotExist
	}
	t, ok := p.tpls[path]
	if !ok {
		t, err = ParseFSFile(p.fs, path+".yap")
		if err != nil {
			return
		}
		p.tpls[path] = t
	}
	return
}
