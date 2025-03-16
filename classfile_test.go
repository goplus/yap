/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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
	"testing"

	"github.com/goplus/yap"
)

type AppV2 struct {
	yap.AppV2
}

func (p *AppV2) Main() {
	yap.Gopt_AppV2_Main(p, &handlerV2{fname: "get_p_#id"}, &handlerV2{fname: "handle"})
}

type handlerV2 struct {
	yap.Handler
	*AppV2
	fname string
}

func (p *handlerV2) Main(ctx *yap.Context) {
	p.Handler.Main(ctx)
	ctx.Json__1(yap.H{"msg": "Hello, Go+!"})
}

func (p *handlerV2) Classfname() string {
	return p.fname
}

func (p *handlerV2) Classclone() yap.HandlerProto {
	ret := *p
	return &ret
}

func TestClassfile(t *testing.T) {
	tr := mock("example.com", new(AppV2))

	c := http.Client{Transport: tr}

	resp, err := c.Get("http://example.com/p/123")
	if err != nil {
		t.Fatal("GET /p/123 failed:", err)
	}
	defer resp.Body.Close()

	resp, err = c.Get("http://example.com/")
	if err != nil {
		t.Fatal("GET / failed:", err)
	}
	defer resp.Body.Close()
}
