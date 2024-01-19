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
	"log"
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

func (p *Response) MatchCode(code any) {
	switch v := code.(type) {
	case int:
		Match__0(p.code, v)
	case *Var__0[int]:
		v.Match(p.code)
	default:
		log.Panicf("match status code failed! unexpected type: %T\n", code)
	}
}

func (p *Response) Header() http.Header {
	return p.header
}

func (p *Response) MatchHeader(key string, value any) {
	switch v := value.(type) {
	case string:
		Match__0(v, p.header.Get(key))
	case []string:
		Match__3(v, p.header[key])
	case *Var__0[string]:
		v.Match(p.header.Get(key))
	case *Var__3[[]string]:
		v.Match(p.header[key])
	default:
		log.Panicf("match header failed! unexpected value type: %T\n", value)
	}
}

func (p *Response) MatchBody(bodyType string, body any) {
}

func (p *Response) MatchJson(body any) {
}

func (p *Response) MatchForm(body any) {
}
