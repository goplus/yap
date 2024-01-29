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
	"time"
)

// -----------------------------------------------------------------------------

type caseT interface {
	// Name returns the name of the running (sub-) test or benchmark.
	//
	// The name will include the name of the test along with the names of
	// any nested sub-tests. If two sibling sub-tests have the same name,
	// Name will append a suffix to guarantee the returned name is unique.
	Name() string

	// Fail marks the function as having failed but continues execution.
	Fail()

	// Failed reports whether the function has failed.
	Failed() bool

	// FailNow marks the function as having failed and stops its execution
	// by calling runtime.Goexit (which then runs all deferred calls in the
	// current goroutine).
	// Execution will continue at the next test or benchmark.
	// FailNow must be called from the goroutine running the
	// test or benchmark function, not from other goroutines
	// created during the test. Calling FailNow does not stop
	// those other goroutines.
	FailNow()

	// Log formats its arguments using default formatting, analogous to Println,
	// and records the text in the error log. For tests, the text will be printed only if
	// the test fails or the -test.v flag is set. For benchmarks, the text is always
	// printed to avoid having performance depend on the value of the -test.v flag.
	Log(args ...any)

	// Logf formats its arguments according to the format, analogous to Printf, and
	// records the text in the error log. A final newline is added if not provided. For
	// tests, the text will be printed only if the test fails or the -test.v flag is
	// set. For benchmarks, the text is always printed to avoid having performance
	// depend on the value of the -test.v flag.
	Logf(format string, args ...any)

	// Errorln is equivalent to Log followed by Fail.
	Errorln(args ...any)

	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...any)

	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...any)

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...any)

	// Skip is equivalent to Log followed by SkipNow.
	Skip(args ...any)

	// Skipf is equivalent to Logf followed by SkipNow.
	Skipf(format string, args ...any)

	// SkipNow marks the test as having been skipped and stops its execution
	// by calling runtime.Goexit.
	// If a test fails (see Error, Errorf, Fail) and is then skipped,
	// it is still considered to have failed.
	// Execution will continue at the next test or benchmark. See also FailNow.
	// SkipNow must be called from the goroutine running the test, not from
	// other goroutines created during the test. Calling SkipNow does not stop
	// those other goroutines.
	SkipNow()

	// Skipped reports whether the test was skipped.
	Skipped() bool

	// Helper marks the calling function as a test helper function.
	// When printing file and line information, that function will be skipped.
	// Helper may be called simultaneously from multiple goroutines.
	Helper()

	// Cleanup registers a function to be called when the test (or subtest) and all its
	// subtests complete. Cleanup functions will be called in last added,
	// first called order.
	Cleanup(f func())

	// TempDir returns a temporary directory for the test to use.
	// The directory is automatically removed by Cleanup when the test and
	// all its subtests complete.
	// Each subsequent call to t.TempDir returns a unique directory;
	// if the directory creation fails, TempDir terminates the test by calling Fatal.
	TempDir() string

	// Run runs f as a subtest of t called name.
	//
	// Run may be called simultaneously from multiple goroutines, but all such calls
	// must return before the outer test function for t returns.
	Run(name string, f func()) bool

	// Deadline reports the time at which the test binary will have
	// exceeded the timeout specified by the -timeout flag.
	//
	// The ok result is false if the -timeout flag indicates “no timeout” (0).
	Deadline() (deadline time.Time, ok bool)
}

type testingT struct {
	*testing.T
}

// Errorln is equivalent to Log followed by Fail.
func (p testingT) Errorln(args ...any) {
	p.T.Error(args...)
}

// Run runs f as a subtest of t called name.
//
// Run may be called simultaneously from multiple goroutines, but all such calls
// must return before the outer test function for t returns.
func (p testingT) Run(name string, f func()) bool {
	return p.T.Run(name, func(t *testing.T) { f() })
}

// -----------------------------------------------------------------------------

type Case struct {
	*Request
	*App
	caseT

	DefaultHeader http.Header
}

func New() *Case {
	return &Case{}
}

// Gopt_Case_TestMain is required by Go+ compiler as the entry of a YAP test case.
func Gopt_Case_TestMain(c interface{ initCase(*App, caseT) }, t *testing.T) {
	app := new(App).initApp()
	c.initCase(app, testingT{t})
	c.(interface{ Main() }).Main()
}

func (p *Case) initCase(app *App, t caseT) {
	p.App = app
	p.caseT = t
	p.DefaultHeader = make(http.Header)
}

// Req create a new request given a method and url.
func (p *Case) Req(method, url string) *Request {
	req := newRequest(p, method, url)
	p.Request = req
	return req
}

// Get is a shortcut for Req(http.MethodGet, url)
func (p *Case) Get(url string) *Request {
	return p.Req(http.MethodGet, url)
}

// Post is a shortcut for Req(http.MethodPost, url)
func (p *Case) Post(url string) *Request {
	return p.Req(http.MethodPost, url)
}

// Head is a shortcut for Req(http.MethodHead, url)
func (p *Case) Head(url string) *Request {
	return p.Req(http.MethodHead, url)
}

// Put is a shortcut for Req(http.MethodPut, url)
func (p *Case) Put(url string) *Request {
	return p.Req(http.MethodPut, url)
}

// Options is a shortcut for Req(http.MethodOptions, url)
func (p *Case) Options(url string) *Request {
	return p.Req(http.MethodOptions, url)
}

// Patch is a shortcut for Req(http.MethodPatch, url)
func (p *Case) Patch(url string) *Request {
	return p.Req(http.MethodPatch, url)
}

// -----------------------------------------------------------------------------

// GET is a shortcut for Req(http.MethodGet, url)
func (p *Case) GET(url string) *Request {
	return p.Req(http.MethodGet, url)
}

// POST is a shortcut for Req(http.MethodPost, url)
func (p *Case) POST(url string) *Request {
	return p.Req(http.MethodPost, url)
}

// HEAD is a shortcut for Req(http.MethodHead, url)
func (p *Case) HEAD(url string) *Request {
	return p.Req(http.MethodHead, url)
}

// PUT is a shortcut for Req(http.MethodPut, url)
func (p *Case) PUT(url string) *Request {
	return p.Req(http.MethodPut, url)
}

// OPTIONS is a shortcut for Req(http.MethodOptions, url)
func (p *Case) OPTIONS(url string) *Request {
	return p.Req(http.MethodOptions, url)
}

// PATCH is a shortcut for Req(http.MethodPatch, url)
func (p *Case) PATCH(url string) *Request {
	return p.Req(http.MethodPatch, url)
}

// DELETE is a shortcut for Req(http.MethodDelete, url)
func (p *Case) DELETE(url string) *Request {
	return p.Req(http.MethodDelete, url)
}

// -----------------------------------------------------------------------------
