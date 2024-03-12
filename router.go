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
	"net/http"
	"strings"

	"github.com/goplus/yap/internal/url"
)

// router is a http rounter which can be used to dispatch requests to different
// handler functions via configurable routes
type router struct {
	trees map[string]*node

	// An optional http.Handler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS http.Handler

	// Cached value of global (*) allowed methods
	globalAllowed string

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool
}

func (p *router) init() {
	p.RedirectTrailingSlash = true
	p.RedirectFixedPath = true
	p.HandleMethodNotAllowed = true
	p.HandleOPTIONS = true
}

// GET is a shortcut for router.Route(http.MethodGet, path, handle)
func (p *router) GET(path string, handle func(ctx *Context)) {
	p.Route(http.MethodGet, path, handle)
}

// HEAD is a shortcut for router.Route(http.MethodHead, path, handle)
func (p *router) HEAD(path string, handle func(ctx *Context)) {
	p.Route(http.MethodHead, path, handle)
}

// OPTIONS is a shortcut for router.Route(http.MethodOptions, path, handle)
func (p *router) OPTIONS(path string, handle func(ctx *Context)) {
	p.Route(http.MethodOptions, path, handle)
}

// POST is a shortcut for router.Route(http.MethodPost, path, handle)
func (p *router) POST(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPost, path, handle)
}

// PUT is a shortcut for router.Route(http.MethodPut, path, handle)
func (p *router) PUT(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPut, path, handle)
}

// PATCH is a shortcut for router.Route(http.MethodPatch, path, handle)
func (p *router) PATCH(path string, handle func(ctx *Context)) {
	p.Route(http.MethodPatch, path, handle)
}

// DELETE is a shortcut for router.Route(http.MethodDelete, path, handle)
func (p *router) DELETE(path string, handle func(ctx *Context)) {
	p.Route(http.MethodDelete, path, handle)
}

// Route registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (p *router) Route(method, path string, handle func(ctx *Context)) {
	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handle == nil {
		panic("handle must not be nil")
	}

	if p.trees == nil {
		p.trees = make(map[string]*node)
	}

	root := p.trees[method]
	if root == nil {
		root = new(node)
		p.trees[method] = root

		p.globalAllowed = p.allowed("*", "")
	}

	root.addRoute(path, handle)
}

func (p *router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		p.PanicHandler(w, req, rcv)
	}
}

func (p *router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range p.trees {
				if method == http.MethodOptions {
					continue
				}
				// Route request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return p.globalAllowed
		}
	} else { // specific path
		for method := range p.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _ := p.trees[method].getValue(path, nil)
			if handle != nil {
				// Route request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Route request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}

	return allow
}

func (p *router) serveHTTP(w http.ResponseWriter, req *http.Request, e *Engine) {
	if p.PanicHandler != nil {
		defer p.recv(w, req)
	}

	path := req.URL.Path
	root := p.trees[req.Method]
	if root != nil {
		ctx := e.NewContext(w, req)
		if handle, tsr := root.getValue(path, ctx); handle != nil {
			handle(ctx)
			return
		} else if req.Method != http.MethodConnect && path != "/" {
			// Moved Permanently, request with GET method
			code := http.StatusMovedPermanently
			if req.Method != http.MethodGet {
				// Permanent Redirect, request with same method
				code = http.StatusPermanentRedirect
			}

			if tsr && p.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if p.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					url.CleanPath(path),
					p.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = fixedPath
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	} else if req.Method == http.MethodHead {
		p.head(w, req, e)
		return
	}

	if req.Method == http.MethodOptions && p.HandleOPTIONS {
		// Route OPTIONS requests
		if allow := p.allowed(path, http.MethodOptions); allow != "" {
			w.Header().Set("Allow", allow)
			if p.GlobalOPTIONS != nil {
				p.GlobalOPTIONS.ServeHTTP(w, req)
			}
			return
		}
	} else if p.HandleMethodNotAllowed { // Route 405
		if allow := p.allowed(path, req.Method); allow != "" {
			w.Header().Set("Allow", allow)
			if p.MethodNotAllowed != nil {
				p.MethodNotAllowed.ServeHTTP(w, req)
			} else {
				http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}
			return
		}
	}

	e.Mux.ServeHTTP(w, req)
}

func (p *router) head(w http.ResponseWriter, req *http.Request, e *Engine) {
	req.Method = http.MethodGet
	p.serveHTTP(&headWriter{w}, req, e)
}

type headWriter struct {
	http.ResponseWriter
}

func (p *headWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
