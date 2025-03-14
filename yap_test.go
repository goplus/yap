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

package yap_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/goplus/yap"
)

func TestBasic(t *testing.T) {
	y := yap.New(os.DirFS("."))

	y.GET("/", func(ctx *yap.Context) {
		ctx.TEXT(200, "text/html", `<html><body>Hello, <a href="/p/123">YAP</a>!</body></html>`)
	})
	y.GET("/p/:id", func(ctx *yap.Context) {
		ctx.YAP(200, "article", yap.H{
			"id": ctx.Param("id"),
		})
	})
	y.SetLAS(func(addr string, h http.Handler) error {
		return nil
	})
	y.Run(":8888")
}

type handler struct{}

func (p *handler) Main(ctx *yap.Context) {
	ctx.JSON(200, yap.H{
		"msg": "Hello, YAP!",
	})
}

func (p *handler) Classclone() any {
	ret := *p
	return &ret
}

func TestProto(t *testing.T) {
	y := yap.New(os.DirFS("."))

	y.ProtoHandle("/", new(handler))
	y.ProtoRoute("GET", "/p/:id", new(handler))
	y.SetLAS(func(addr string, h http.Handler) error {
		return nil
	})
	y.Run(":8888")
}
