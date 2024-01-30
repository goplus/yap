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

type Response struct {
	code   int
	header http.Header
}

func newResponse(resp *http.Response) *Response {
	return &Response{
		code:   resp.StatusCode,
		header: resp.Header,
	}
}

func (p *Response) Code() int {
	return p.code
}

func (p *Response) MatchCode(t CaseT, code any) {
	t.Helper()
	switch v := code.(type) {
	case int:
		Gopt_Case_Match__0(t, p.code, v)
	case *Var__0[int]:
		v.Match(t, p.code)
	default:
		t.Fatalf("match status code failed! unexpected type: %T\n", code)
	}
}

func (p *Response) Header() http.Header {
	return p.header
}

func (p *Response) MatchHeader(t CaseT, key string, value any) {
	t.Helper()
	switch v := value.(type) {
	case string:
		Gopt_Case_Match__0(t, v, p.header.Get(key))
	case []string:
		Gopt_Case_Match__3(t, v, p.header[key])
	case *Var__0[string]:
		v.Match(t, p.header.Get(key))
	case *Var__3[[]string]:
		v.Match(t, p.header[key])
	default:
		t.Fatalf("match header failed! unexpected value type: %T\n", value)
	}
}

func (p *Response) Body() any {
	return nil // TODO
}

func (p *Response) MatchBody(t CaseT, bodyType string, body any) {
	t.Helper()
	// t.Fatal("todo")
}
