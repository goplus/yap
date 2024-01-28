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
	"io"
	"net/http"
	"strings"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

type App struct {
	hosts     map[string]string
	transport http.RoundTripper
}

func (p *App) initApp() *App {
	p.hosts = make(map[string]string)
	p.transport = http.DefaultTransport
	return p
}

// Gopt_App_Main is required by Go+ compiler as the entry of a YAP testing project.
func Gopt_App_Main(app interface{ initApp() *App }, workers ...interface{ initCase(*App) }) {
	a := app.initApp()
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	for _, worker := range workers {
		worker.initCase(a)
		worker.(interface{ Main() }).Main()
	}
}

// Host replaces a host into real. For example:
//
//	host "https://example.com" "http://localhost:8080"
//	host "http://example.com" "http://localhost:8888"
func (p *App) Host(host, real string) {
}

func (p *App) hostOf(url string) (host string, url2 string, ok bool) {
	// http://host/xxx or https://host/xxx
	var istart int
	if strings.HasPrefix(url[4:], "://") {
		istart = 7
	} else if strings.HasPrefix(url[4:], "s://") {
		istart = 8
	} else {
		return
	}

	next := url[istart:]
	n := strings.IndexByte(next, '/')
	if n < 1 {
		return
	}

	host = next[:n]
	portal, ok := p.hosts[url[:istart+n]]
	if ok {
		url2 = portal + url[istart+n:]
	}
	return
}

func (p *App) newRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	if host, url2, ok := p.hostOf(url); ok {
		req, err = http.NewRequest(method, url2, body)
		req.Host = host
		return
	}
	return http.NewRequest(method, url, body)
}

// -----------------------------------------------------------------------------

func Oauth2(auth string) RTComposer {
	return nil
}

func Jwt(auth string) RTComposer {
	return &JwtAuth{
		JwtToken: auth,
	}
}

// -----------------------------------------------------------------------------
