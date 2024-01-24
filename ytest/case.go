/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

package ytest

import (
	"net/http"
)

type Case struct {
	*Request
	*App
	name string

	DefaultHeader http.Header
}

func New() *Case {
	return &Case{}
}

func (p *Case) initCase(app *App) {
	p.App = app
	p.DefaultHeader = make(http.Header)
}

func (p *Case) Run(name string, doSth func()) {
	p.name = name
	doSth()
}

func (p *Case) Req(method, url string) *Request {
	req := newRequest(p, method, url)
	p.Request = req
	return req
}

func (p *Case) Get(url string) *Request {
	return p.Req(http.MethodGet, url)
}

func (p *Case) Post(url string) *Request {
	return p.Req(http.MethodPost, url)
}
