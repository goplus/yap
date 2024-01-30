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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/goplus/yap/ytest/auth"
)

type RequestBody interface {
	io.Reader
	io.Seeker
	Size() int64
}

type Request struct {
	method   string
	url      string
	header   http.Header
	auth     auth.RTComposer
	bodyType string
	body     RequestBody
	resp     *Response
	ctx      *Case
}

func newRequest(ctx *Case, method, url string) *Request {
	return &Request{
		method: method,
		url:    url,
		header: make(http.Header),
		ctx:    ctx,
	}
}

func (p *Request) t() CaseT {
	return p.ctx.CaseT
}

// Auth sets an Authorization for this request.
func (p *Request) Auth(auth auth.RTComposer) *Request {
	p.auth = auth
	return p
}

// -----------------------------------------------------------------------------

// WithHeader sets a Header for this request.
func (p *Request) WithHeader(key string, value any) *Request {
	switch v := value.(type) {
	case string:
		p.header.Set(key, v)
	case []string:
		p.header[key] = v
	case *Var__0[string]:
		p.header.Set(key, v.Val())
	case *Var__3[[]string]:
		p.header[key] = v.Val()
	default:
		fatalf("set header failed! unexpected value type: %T\n", value)
	}
	return p
}

// Header sets a Header for this request (if request is not sended), or matches
// a Header for response of this request (after response is returned).
// Here value can be: string, []string, Var(string), Var([]string).
func (p *Request) Header(key string, value any) *Request {
	if p.resp == nil {
		return p.WithHeader(key, value)
	}
	t := p.t()
	t.Helper()
	p.resp.MatchHeader(t, key, value)
	return p
}

// -----------------------------------------------------------------------------

// Body sets request body for this request (if request is not sended), or matches
// response body of this request (after response is returned).
// Here body can be: string, Var(string), []byte, RequestBody.
func (p *Request) Body(bodyType string, body any) *Request {
	if p.resp == nil {
		return p.WithBodyEx(bodyType, body)
	}
	t := p.t()
	t.Helper()
	p.resp.MatchBody(t, bodyType, body)
	return p
}

func (p *Request) WithBodyEx(bodyType string, body any) *Request {
	switch v := body.(type) {
	case string:
		return p.WithText(bodyType, v)
	case *Var__0[string]:
		return p.WithText(bodyType, v.Val())
	case []byte:
		return p.WithBinary(bodyType, v)
	case RequestBody:
		return p.WithBody(bodyType, v)
	default:
		fatalf("set body failed! unexpected value type: %T\n", body)
	}
	return p
}

func (p *Request) WithText(bodyType, body string) *Request {
	return p.WithBody(bodyType, strings.NewReader(body))
}

func (p *Request) WithBinary(bodyType string, body []byte) *Request {
	return p.WithBody(bodyType, bytes.NewReader(body))
}

func (p *Request) WithBody(bodyType string, body RequestBody) *Request {
	p.bodyType = mimeType(bodyType)
	p.body = body
	return p
}

const (
	mimeForm   = "application/x-www-form-urlencoded"
	mimeJson   = "application/json"
	mimeBinary = "application/octet-stream"
	mimeText   = "text/plain"
)

func mimeType(ct string) string {
	if strings.Contains(ct, "/") {
		return ct
	}
	switch ct {
	case "form":
		return mimeForm
	case "binary":
		return mimeBinary
	}
	return "application/" + ct
}

func (p *Request) Binary(body any) *Request {
	p.t().Helper()
	return p.Body(mimeBinary, body)
}

func (p *Request) Text(body any) *Request {
	p.t().Helper()
	return p.Body(mimeText, body)
}

// -----------------------------------------------------------------------------

// body can be:
//   - map[string]any, Var(map[string]any), []any, Var([]any),
//   - []string, Var([]string), string, Var(string), int, Var(int),
//   - bool, Var(bool), float64, Var(float64).
func (p *Request) Json(body any) *Request {
	if p.resp == nil {
		return p.WithJson(body)
	}
	t := p.t()
	t.Helper()
	p.resp.MatchBody(t, mimeJson, body)
	return p
}

func (p *Request) WithJson(body any) *Request {
	b, err := json.Marshal(body)
	if err != nil {
		fatal("json.Marshal failed:", err)
	}
	return p.WithBinary(mimeJson, b)
}

// -----------------------------------------------------------------------------

func (p *Request) Form(body any) *Request {
	if p.resp == nil {
		return p.WithFormEx(body)
	}
	t := p.t()
	t.Helper()
	p.resp.MatchBody(t, mimeForm, body)
	return p
}

func (p *Request) WithForm(body url.Values) *Request {
	return p.WithText(mimeForm, body.Encode())
}

func (p *Request) WithFormEx(body any) *Request {
	var vals url.Values
	switch v := body.(type) {
	case map[string]any:
		vals = Form(v)
	case *Var__1[map[string]any]:
		vals = Form(v.Val())
	case url.Values:
		vals = v
	default:
		fatalf("request with form: unexpected type %T\n", body)
	}
	return p.WithText(mimeForm, vals.Encode())
}

// -----------------------------------------------------------------------------

func mergeHeader(to, from http.Header) {
	for k, v := range from {
		to[k] = v
	}
}

func (p *Request) doSend() (resp *http.Response, err error) {
	body := p.body
	req, err := p.ctx.newRequest(p.method, p.url, body)
	if err != nil {
		fatalf("newRequest(%s, %s) failed: %v\n", p.method, p.url, err)
	}

	mergeHeader(req.Header, p.ctx.DefaultHeader)
	mergeHeader(req.Header, p.header)

	if body != nil {
		if p.bodyType != "" {
			req.Header.Set("Content-Type", p.bodyType)
		}
		req.ContentLength = body.Size()
	}
	tr := p.ctx.transport
	if p.auth != nil {
		tr = p.auth.Compose(tr)
	}
	c := &http.Client{Transport: tr}
	return c.Do(req)
}

const (
	Gopo_Request_Ret = ".Send,.RetWith"
)

func (p *Request) Send() *Request {
	resp, err := p.doSend()
	if err != nil {
		fatalf("sendRequest(%v, %v) failed: %v\n", p.method, p.url, err)
	}
	p.resp = newResponse(resp)
	return p
}

func (p *Request) RetWith(code any) *Request {
	t := p.t()
	t.Helper()
	p.Send().resp.MatchCode(t, code)
	return p
}

// -----------------------------------------------------------------------------

func (p *Request) Resp() *Response {
	return p.resp
}

// -----------------------------------------------------------------------------
