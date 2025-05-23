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
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Context struct {
	*http.Request
	http.ResponseWriter

	engine *Engine
}

func (p *Context) setParam(name, val string) {
	p.ParseForm()
	p.Form.Set(name, val)
}

// Gop_Env returns the value associated with the name.
// If the name exists in URL query, it returns the first value for the name.
func (p *Context) Gop_Env(name string) string {
	return p.Param(name)
}

// Param returns the value associated with the name.
// If the name exists in URL query, it returns the first value for the name.
func (p *Context) Param(name string) string {
	return p.FormValue(name)
}

func (p *Context) ParamInt(name string, defval int) int {
	ret := p.Param(name)
	if ret != "" {
		if v, err := strconv.Atoi(ret); err == nil {
			return v
		}
	}
	return defval
}

// Accept header specifies:
// Accept: <MIME_type>/<MIME_subtype>
// Accept: <MIME_type>/*
// Accept: */*
// Multiple types, weighted with the quality value syntax:
// Accept: text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8
// FIXME: 1. quality value not supported, 2. don't need parse all, just find the first match with a spliter iterator
func (p *Context) Accept(mime ...string) string {
	accept := p.Request.Header.Get("Accept")
	for _, m := range mime {
		if acceptMime(accept, m) {
			return m
		}
	}
	return ""
}

func acceptMime(accept string, mime string) bool {
	for accept != "" {
		item, left := acceptNext(accept)
		if item == mime || (strings.HasPrefix(item, mime) && item[len(mime)] == ';') {
			return true
		}
		accept = left
	}
	return false
}

func acceptNext(accept string) (item, left string) {
	item = strings.TrimLeft(accept, " ")
	if before, after, found := strings.Cut(item, ","); found {
		return before, after
	}
	left = ""
	return
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
func (p *Context) Redirect(url string, code ...int) {
	statusCode := http.StatusFound
	if code != nil {
		statusCode = code[0]
	}
	http.Redirect(p.ResponseWriter, p.Request, url, statusCode)
}

func (p *Context) TEXT(code int, mime string, text string) {
	w := p.ResponseWriter
	h := w.Header()
	h.Set("Content-Length", strconv.Itoa(len(text)))
	h.Set("Content-Type", mime)
	w.WriteHeader(code)
	io.WriteString(w, text)
}

func (p *Context) DATA(code int, mime string, data []byte) {
	w := p.ResponseWriter
	h := w.Header()
	h.Set("Content-Length", strconv.Itoa(len(data)))
	h.Set("Content-Type", mime)
	w.WriteHeader(code)
	w.Write(data)
}

func (p *Context) PrettyJSON(code int, data any) {
	msg, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	p.DATA(code, "application/json", msg)
}

func (p *Context) JSON(code int, data any) {
	msg, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	p.DATA(code, "application/json", msg)
}

func (p *Context) YAP(code int, yapFile string, data any) {
	w := p.ResponseWriter
	t := p.engine.templ(yapFile)
	if t == nil {
		log.Panicln("YAP: not find template:", yapFile)
	}
	h := w.Header()
	h.Set("Content-Type", "text/html")
	err := t.Execute(w, data)
	if err != nil {
		log.Panicln("YAP:", err)
	}
}

func (p *Context) STREAM(code int, mime string, read io.Reader, buf []byte) {
	w := p.ResponseWriter
	h := w.Header()
	if mime != "" {
		h.Set("Content-Type", mime)
	}
	w.WriteHeader(code)

	// Auto flush if the buffer is small
	if buf != nil && cap(buf) < 32*1024 {
		if f, ok := w.(http.Flusher); ok {
			w = newAutoFlushWriter(w, f)
		}
	}
	io.CopyBuffer(w, read, buf)
}

// autoFlushWriter wraps http.ResponseWriter and implements automatic flushing
type autoFlushWriter struct {
	http.ResponseWriter // Embeds the original interface
	flusher             http.Flusher
}

// newAutoFlushWriter creates a new AutoFlushWriter instance
func newAutoFlushWriter(w http.ResponseWriter, f http.Flusher) *autoFlushWriter {
	return &autoFlushWriter{
		ResponseWriter: w,
		flusher:        f,
	}
}

// Write implements the write method with automatic flushing
func (w *autoFlushWriter) Write(p []byte) (n int, err error) {
	// Call the original write method
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		return
	}

	// Immediately flush the buffer
	w.flusher.Flush()
	return
}
