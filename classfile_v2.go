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
	"net/http"
	"reflect"
	"strings"
)

// HandlerProto is the prototype of a YAP handler.
type HandlerProto interface {
	Main(ctx *Context)
	Classclone() any
}

// ProtoHandle registers a YAP handler with a prototype.
func (p *Engine) ProtoHandle(pattern string, proto HandlerProto) {
	p.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// ensure isolation of handler state per request
		h := proto.Classclone().(HandlerProto)
		h.Main(p.NewContext(w, r))
	})
}

// ProtoRoute registers a YAP handler with a prototype.
func (p *router) ProtoRoute(method, path string, proto HandlerProto) {
	p.Route(method, path, func(ctx *Context) {
		// ensure isolation of handler state per request
		h := proto.Classclone().(HandlerProto)
		h.Main(ctx)
	})
}

// Handler is worker class of YAP classfile (v2).
type Handler struct {
	Context
}

// Main is required by Go+ compiler as the entry of a YAP HTTP handler.
func (p *Handler) Main(ctx *Context) {
	p.Context = *ctx
}

var (
	repl = strings.NewReplacer("_", "/", "#", ":")
)

func parseClassfname(name string) (method, path string) {
	pos := strings.IndexByte(name, '_')
	if pos < 0 {
		return name, "/"
	}
	return name[:pos], repl.Replace(name[pos:])
}

// AppV2 is project class of YAP classfile (v2).
type AppV2 struct {
	App
}

type iHandlerProto interface {
	HandlerProto
	Classfname() string
}

// Gopt_AppV2_Main is required by Go+ compiler as the entry of a YAP project.
func Gopt_AppV2_Main(app AppType, handlers ...iHandlerProto) {
	app.InitYap()
	for _, h := range handlers {
		reflect.ValueOf(h).Elem().Field(1).Set(reflect.ValueOf(app)) // (*handler).AppV2 = app
		switch method, path := parseClassfname(h.Classfname()); method {
		case "handle":
			app.ProtoHandle(path, h)
		default:
			app.ProtoRoute(strings.ToUpper(method), path, h)
		}
	}
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	} else {
		app.Run("localhost:8080")
	}
}
