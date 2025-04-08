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
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/yap"
	"github.com/qiniu/x/mockhttp"
	"github.com/qiniu/x/test"
	"github.com/qiniu/x/test/logt"
)

const (
	GopPackage   = "github.com/qiniu/x/test"
	GopTestClass = true
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

// Gop_Env retrieves the value of the environment variable named by the key.
func (p *App) Gop_Env(key string) string {
	return os.Getenv(key)
}

// Mock runs a YAP server by mockhttp.
func (p *App) Mock(host string, app yap.AppType) {
	tr := mockhttp.NewTransport()
	p.transport = tr
	app.InitYap()
	app.SetLAS(func(addr string, h http.Handler) error {
		return tr.ListenAndServe(host, h)
	})
	app.(interface{ Main() }).Main()
}

// TestServer runs a YAP server by httptest.Server.
func (p *App) TestServer(host string, app yap.AppType) {
	app.InitYap()
	app.SetLAS(func(addr string, h http.Handler) error {
		svr := httptest.NewServer(h)
		p.Host("http://"+host, svr.URL)
		return nil
	})
	app.(interface{ Main() }).Main()
}

// RunMock runs a HTTP server by mockhttp.
func (p *App) RunMock(host string, h http.Handler) {
	tr := mockhttp.NewTransport()
	p.transport = tr
	tr.ListenAndServe(host, h)
}

// RunTestServer runs a HTTP server by httptest.Server.
func (p *App) RunTestServer(host string, h http.Handler) {
	svr := httptest.NewServer(h)
	p.Host("http://"+host, svr.URL)
}

// Gopt_App_TestMain is required by Go+ compiler as the TestMain entry of a YAP testing project.
func Gopt_App_TestMain(app interface{ initApp() *App }, m *testing.M) {
	app.initApp()
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	os.Exit(m.Run())
}

// Gopt_App_Main is required by Go+ compiler as the Main entry of a YAP testing project.
func Gopt_App_Main(app interface{ initApp() *App }, workers ...interface{ initCase(*App, CaseT) }) {
	a := app.initApp()
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	t := logt.New()
	for _, worker := range workers {
		worker.initCase(a, t)
		reflect.ValueOf(worker).Elem().Field(1).Set(reflect.ValueOf(app)) // (*worker).App = app
		worker.(interface{ Main() }).Main()
	}
}

// Host replaces a host into real. For example:
//
//	host "https://example.com", "http://localhost:8080"
//	host "http://example.com", "http://localhost:8888"
func (p *App) Host(host, real string) {
	if !strings.HasPrefix(host, "http") {
		test.Fatalf("invalid host `%s`: should start with http:// or https://\n", host)
	}
	p.hosts[host] = real
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
