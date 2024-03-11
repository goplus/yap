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
	"testing"

	"github.com/goplus/yap/test"
)

// -----------------------------------------------------------------------------

type CaseT = test.CaseT

type Case struct {
	*Request
	test.Case
	app *App

	DefaultHeader http.Header
}

func (p *Case) initCase(app *App, t CaseT) {
	p.app = app
	p.CaseT = t
	p.DefaultHeader = make(http.Header)
}

// T returns the testing object.
func (p *Case) T() CaseT {
	return p.CaseT
}

// Req create a new request given a method and url.
func (p *Case) Req__0(method, url string) *Request {
	req := newRequest(p, method, url)
	p.Request = req
	return req
}

// Req returns current request object.
func (p *Case) Req__1() *Request {
	return p.Request
}

// Get is a shortcut for Req(http.MethodGet, url)
func (p *Case) Get(url string) *Request {
	return p.Req__0(http.MethodGet, url)
}

// Post is a shortcut for Req(http.MethodPost, url)
func (p *Case) Post(url string) *Request {
	return p.Req__0(http.MethodPost, url)
}

// Head is a shortcut for Req(http.MethodHead, url)
func (p *Case) Head(url string) *Request {
	return p.Req__0(http.MethodHead, url)
}

// Put is a shortcut for Req(http.MethodPut, url)
func (p *Case) Put(url string) *Request {
	return p.Req__0(http.MethodPut, url)
}

// Options is a shortcut for Req(http.MethodOptions, url)
func (p *Case) Options(url string) *Request {
	return p.Req__0(http.MethodOptions, url)
}

// Patch is a shortcut for Req(http.MethodPatch, url)
func (p *Case) Patch(url string) *Request {
	return p.Req__0(http.MethodPatch, url)
}

// -----------------------------------------------------------------------------

// GET is a shortcut for Req(http.MethodGet, url)
func (p *Case) GET(url string) *Request {
	return p.Req__0(http.MethodGet, url)
}

// POST is a shortcut for Req(http.MethodPost, url)
func (p *Case) POST(url string) *Request {
	return p.Req__0(http.MethodPost, url)
}

// HEAD is a shortcut for Req(http.MethodHead, url)
func (p *Case) HEAD(url string) *Request {
	return p.Req__0(http.MethodHead, url)
}

// PUT is a shortcut for Req(http.MethodPut, url)
func (p *Case) PUT(url string) *Request {
	return p.Req__0(http.MethodPut, url)
}

// OPTIONS is a shortcut for Req(http.MethodOptions, url)
func (p *Case) OPTIONS(url string) *Request {
	return p.Req__0(http.MethodOptions, url)
}

// PATCH is a shortcut for Req(http.MethodPatch, url)
func (p *Case) PATCH(url string) *Request {
	return p.Req__0(http.MethodPatch, url)
}

// DELETE is a shortcut for Req(http.MethodDelete, url)
func (p *Case) DELETE(url string) *Request {
	return p.Req__0(http.MethodDelete, url)
}

// -----------------------------------------------------------------------------

type CaseApp struct {
	Case
	*App
}

// Gopt_CaseApp_TestMain is required by Go+ compiler as the entry of a YAP test case.
func Gopt_CaseApp_TestMain(c interface{ initCaseApp(*App, CaseT) }, t *testing.T) {
	app := new(App).initApp()
	c.initCaseApp(app, test.NewT(t))
	c.(interface{ Main() }).Main()
}

func (p *CaseApp) initCaseApp(app *App, t CaseT) {
	p.initCase(app, t)
	p.App = app
}

// -----------------------------------------------------------------------------
