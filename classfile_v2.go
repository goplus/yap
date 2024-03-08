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
	"strings"
)

// Handler is worker class of YAP classfile (v2).
type Handler struct {
	Context
}

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

type iHandler interface {
	Main(ctx *Context)
	Classfname() string
}

// Gopt_AppV2_Main is required by Go+ compiler as the entry of a YAP project.
func Gopt_AppV2_Main(app AppType, handlers ...iHandler) {
	app.InitYap()
	for _, h := range handlers {
		switch method, path := parseClassfname(h.Classfname()); method {
		case "handle":
			app.Handle(path, h.Main)
		default:
			app.Route(strings.ToUpper(method), path, h.Main)
		}
	}
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	} else {
		app.Run(":8080")
	}
}
