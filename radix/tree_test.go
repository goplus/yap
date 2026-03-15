/*
 * Copyright (c) 2023 The XGo Authors (xgo.dev). All rights reserved.
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

package radix

import (
	"testing"
)

// testContext implements the context interface for testing.
type testContext struct {
	params map[string]string
}

func (c *testContext) UnderlyingSetPathParam(name, val string) {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	c.params[name] = val
}

func newTestCtx() *testContext {
	return &testContext{params: make(map[string]string)}
}

// -----------------------------------------------------------------------------
// AddRoute / Route: static routes

func TestAddRouteGetValueSingleStatic(t *testing.T) {
	var root Node[func(*testContext)]
	called := false
	root.AddRoute("/users", func(ctx *testContext) { called = true })

	ctx := newTestCtx()
	h, _, tsr := Route(&root, "/users", ctx)
	if h == nil {
		t.Fatal("expected handle for /users")
	}
	if tsr {
		t.Fatal("unexpected TSR for /users")
	}
	h(ctx)
	if !called {
		t.Fatal("handle was not called")
	}
}

func TestAddRouteGetValueMultipleStatic(t *testing.T) {
	var root Node[func(*testContext)]
	results := map[string]bool{}

	routes := []string{"/", "/users", "/users/list", "/posts", "/posts/list"}
	for _, r := range routes {
		r := r
		root.AddRoute(r, func(ctx *testContext) { results[r] = true })
	}

	for _, r := range routes {
		results[r] = false
		ctx := newTestCtx()
		h, _, tsr := Route(&root, r, ctx)
		if h == nil {
			t.Fatalf("expected handle for %q", r)
		}
		if tsr {
			t.Fatalf("unexpected TSR for %q", r)
		}
		h(ctx)
		if !results[r] {
			t.Fatalf("handle not called for %q", r)
		}
	}
}

func TestGetValueNotFound(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users", func(ctx *testContext) {})

	h, _, _ := Route(&root, "/nonexistent", newTestCtx())
	if h != nil {
		t.Fatal("expected no handle for /nonexistent")
	}
}

func TestGetValueRootOnly(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/", func(ctx *testContext) {})

	h, _, _ := Route(&root, "/", newTestCtx())
	if h == nil {
		t.Fatal("expected handle for /")
	}

	h2, _, _ := Route(&root, "/other", newTestCtx())
	if h2 != nil {
		t.Fatal("expected no handle for /other")
	}
}

// -----------------------------------------------------------------------------
// AddRoute / Route: param routes

func TestGetValueSingleParam(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id", func(ctx *testContext) {})

	ctx := newTestCtx()
	h, _, tsr := Route(&root, "/users/42", ctx)
	if h == nil {
		t.Fatal("expected handle for /users/42")
	}
	if tsr {
		t.Fatal("unexpected TSR")
	}
	if ctx.params["id"] != "42" {
		t.Fatalf("expected param id=42, got %q", ctx.params["id"])
	}
}

func TestGetValueMultipleParams(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id/posts/:postId", func(ctx *testContext) {})

	ctx := newTestCtx()
	h, _, _ := Route(&root, "/users/123/posts/456", ctx)
	if h == nil {
		t.Fatal("expected handle")
	}
	if ctx.params["id"] != "123" {
		t.Fatalf("expected id=123, got %q", ctx.params["id"])
	}
	if ctx.params["postId"] != "456" {
		t.Fatalf("expected postId=456, got %q", ctx.params["postId"])
	}
}

func TestGetValueParamWithTrailingSlash(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id/", func(ctx *testContext) {})

	ctx := newTestCtx()
	h, _, _ := Route(&root, "/users/42/", ctx)
	if h == nil {
		t.Fatal("expected handle")
	}
	if ctx.params["id"] != "42" {
		t.Fatalf("expected id=42, got %q", ctx.params["id"])
	}
}

func TestGetValueParamNilCtx(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id", func(ctx *testContext) {})

	// Pass nil context (zero value for *testContext) - should not panic
	h, _, _ := root.Route("/users/42")
	if h == nil {
		t.Fatal("expected handle even with nil ctx")
	}
}

func TestGetValueParamNotFound(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id", func(ctx *testContext) {})

	// /users alone has no handler
	h, _, _ := Route(&root, "/users", newTestCtx())
	if h != nil {
		t.Fatal("expected no handle for /users without param")
	}
}

// -----------------------------------------------------------------------------
// AddRoute / Route: catchAll routes

func TestGetValueCatchAll(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/files/*filepath", func(ctx *testContext) {})

	ctx := newTestCtx()
	h, _, _ := Route(&root, "/files/path/to/file.txt", ctx)
	if h == nil {
		t.Fatal("expected handle")
	}
	if ctx.params["filepath"] != "/path/to/file.txt" {
		t.Fatalf("expected filepath=/path/to/file.txt, got %q", ctx.params["filepath"])
	}
}

func TestGetValueCatchAllNilCtx(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/files/*filepath", func(ctx *testContext) {})

	h, _, _ := root.Route("/files/a/b/c")
	if h == nil {
		t.Fatal("expected handle with nil ctx")
	}
}

// -----------------------------------------------------------------------------
// TSR (Trailing Slash Redirect)

func TestGetValueTSRMissingTrailingSlash(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/", func(ctx *testContext) {})

	h, _, tsr := Route(&root, "/users", newTestCtx())
	if h != nil {
		t.Fatal("expected no handle for /users (without trailing slash)")
	}
	if !tsr {
		t.Fatal("expected TSR recommendation for /users → /users/")
	}
}

func TestGetValueTSRExtraTrailingSlash(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users", func(ctx *testContext) {})

	h, _, tsr := Route(&root, "/users/", newTestCtx())
	if h != nil {
		t.Fatal("expected no handle for /users/ (with extra trailing slash)")
	}
	_ = tsr
}

func TestGetValueTSRWithParam(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id/", func(ctx *testContext) {})

	h, _, tsr := Route(&root, "/users/42", newTestCtx())
	if h != nil {
		t.Fatal("expected no handle")
	}
	if !tsr {
		t.Fatal("expected TSR recommendation for /users/42 → /users/42/")
	}
}

// -----------------------------------------------------------------------------
// Mixed routes (static + param + catchAll)

func TestMixedRoutes(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/", func(ctx *testContext) {})
	root.AddRoute("/users", func(ctx *testContext) {})
	root.AddRoute("/users/:id", func(ctx *testContext) {})
	root.AddRoute("/users/:id/posts", func(ctx *testContext) {})
	root.AddRoute("/files/*filepath", func(ctx *testContext) {})

	tests := []struct {
		path      string
		wantMatch bool
		wantParam map[string]string
	}{
		{"/", true, nil},
		{"/users", true, nil},
		{"/users/123", true, map[string]string{"id": "123"}},
		{"/users/123/posts", true, map[string]string{"id": "123"}},
		{"/files/img/logo.png", true, map[string]string{"filepath": "/img/logo.png"}},
		{"/notfound", false, nil},
	}

	for _, tt := range tests {
		ctx := newTestCtx()
		h, _, _ := Route(&root, tt.path, ctx)
		if (h != nil) != tt.wantMatch {
			t.Fatalf("Route(%q): got match=%v, want match=%v", tt.path, h != nil, tt.wantMatch)
		}
		for k, v := range tt.wantParam {
			if ctx.params[k] != v {
				t.Fatalf("Route(%q): param %q = %q, want %q", tt.path, k, ctx.params[k], v)
			}
		}
	}
}

// -----------------------------------------------------------------------------
// Priority / child ordering

func TestRoutePriorityOrdering(t *testing.T) {
	var root Node[func(*testContext)]
	// Add routes with shared prefixes
	root.AddRoute("/a", func(ctx *testContext) {})
	root.AddRoute("/b", func(ctx *testContext) {})
	root.AddRoute("/a/b", func(ctx *testContext) {})
	root.AddRoute("/b/c", func(ctx *testContext) {})

	for _, path := range []string{"/a", "/b", "/a/b", "/b/c"} {
		h, _, _ := Route(&root, path, newTestCtx())
		if h == nil {
			t.Fatalf("expected handle for %q after priority reordering", path)
		}
	}
}

// -----------------------------------------------------------------------------
// FindCaseInsensitivePath

func TestFindCaseInsensitivePath(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/Users/List", func(ctx *testContext) {})

	fixedPath, found := root.FindCaseInsensitivePath("/users/list", true)
	if !found {
		t.Fatal("expected case-insensitive match for /users/list")
	}
	if fixedPath != "/Users/List" {
		t.Fatalf("expected /Users/List, got %q", fixedPath)
	}
}

func TestFindCaseInsensitivePathNotFound(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users", func(ctx *testContext) {})

	_, found := root.FindCaseInsensitivePath("/nonexistent", false)
	if found {
		t.Fatal("expected not found for /nonexistent")
	}
}

func TestFindCaseInsensitivePathTrailingSlashFix(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/", func(ctx *testContext) {})

	// Should match with trailing slash added
	fixedPath, found := root.FindCaseInsensitivePath("/USERS", true)
	if !found {
		t.Fatal("expected match with trailing slash fix")
	}
	if fixedPath != "/users/" {
		t.Fatalf("expected /users/, got %q", fixedPath)
	}
}

func TestFindCaseInsensitivePathNoFix(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users/", func(ctx *testContext) {})

	// Without fixTrailingSlash, should not find when trailing slash is missing
	_, found := root.FindCaseInsensitivePath("/USERS", false)
	if found {
		t.Fatal("expected not found without fixTrailingSlash")
	}
}

func TestFindCaseInsensitivePathWithParam(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/Users/:id", func(ctx *testContext) {})

	fixedPath, found := root.FindCaseInsensitivePath("/users/123", true)
	if !found {
		t.Fatal("expected match for /users/123")
	}
	if fixedPath != "/Users/123" {
		t.Fatalf("expected /Users/123, got %q", fixedPath)
	}
}

func TestFindCaseInsensitivePathWithCatchAll(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/Files/*filepath", func(ctx *testContext) {})

	fixedPath, found := root.FindCaseInsensitivePath("/files/path/to/file", true)
	if !found {
		t.Fatal("expected match for /files/path/to/file")
	}
	if fixedPath != "/Files/path/to/file" {
		t.Fatalf("expected /Files/path/to/file, got %q", fixedPath)
	}
}

func TestFindCaseInsensitivePathExactMatch(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/users", func(ctx *testContext) {})

	// Exact match (already lowercase)
	fixedPath, found := root.FindCaseInsensitivePath("/users", true)
	if !found {
		t.Fatal("expected match for /users")
	}
	if fixedPath != "/users" {
		t.Fatalf("expected /users, got %q", fixedPath)
	}
}

// -----------------------------------------------------------------------------
// Panic tests

func TestAddRoutePanicDuplicateHandle(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate handle")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/users", func(ctx *testContext) {})
	root.AddRoute("/users", func(ctx *testContext) {})
}

func TestAddRoutePanicWildcardConflict(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for wildcard conflict")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id", func(ctx *testContext) {})
	root.AddRoute("/users/:name", func(ctx *testContext) {}) // conflicts with :id
}

func TestAddRoutePanicEmptyWildcardName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty wildcard name")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/users/:", func(ctx *testContext) {})
}

func TestAddRoutePanicCatchAllNotAtEnd(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for catch-all not at end of path")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/files/*filepath/extra", func(ctx *testContext) {})
}

func TestAddRoutePanicInvalidWildcard(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid wildcard (nested)")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/users/:id:name/posts", func(ctx *testContext) {})
}

// -----------------------------------------------------------------------------
// Static + param sibling routes (proposal #178)

// TestStaticAndParamSiblings verifies the core scenario from the proposal:
// static routes and a :param wildcard can coexist under the same parent.
func TestStaticAndParamSiblings(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/model/videoModels", func(ctx *testContext) {})
	root.AddRoute("/model/imageModels", func(ctx *testContext) {})
	root.AddRoute("/model/:model", func(ctx *testContext) {})

	tests := []struct {
		path       string
		wantMatch  bool
		wantParam  map[string]string
		wantStatic bool // true if result should NOT set :model param
	}{
		{"/model/videoModels", true, nil, true},
		{"/model/imageModels", true, nil, true},
		{"/model/textModels", true, map[string]string{"model": "textModels"}, false},
		{"/model/anything", true, map[string]string{"model": "anything"}, false},
		{"/model/", false, nil, false},
	}

	for _, tt := range tests {
		ctx := newTestCtx()
		h, _, _ := Route(&root, tt.path, ctx)
		if (h != nil) != tt.wantMatch {
			t.Fatalf("Route(%q): got match=%v, want match=%v", tt.path, h != nil, tt.wantMatch)
		}
		if tt.wantStatic && ctx.params["model"] != "" {
			t.Fatalf("Route(%q): expected static match, but :model param was set to %q", tt.path, ctx.params["model"])
		}
		for k, v := range tt.wantParam {
			if ctx.params[k] != v {
				t.Fatalf("Route(%q): param %q = %q, want %q", tt.path, k, ctx.params[k], v)
			}
		}
	}
}

// TestStaticAndParamSiblingsReverseOrder registers param first, then static children.
func TestStaticAndParamSiblingsReverseOrder(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/model/:model", func(ctx *testContext) {})
	root.AddRoute("/model/videoModels", func(ctx *testContext) {})
	root.AddRoute("/model/imageModels", func(ctx *testContext) {})

	tests := []struct {
		path      string
		wantMatch bool
		wantParam map[string]string
	}{
		{"/model/videoModels", true, nil},
		{"/model/imageModels", true, nil},
		{"/model/textModels", true, map[string]string{"model": "textModels"}},
	}

	for _, tt := range tests {
		ctx := newTestCtx()
		h, _, _ := Route(&root, tt.path, ctx)
		if (h != nil) != tt.wantMatch {
			t.Fatalf("Route(%q): got match=%v, want match=%v", tt.path, h != nil, tt.wantMatch)
		}
		for k, v := range tt.wantParam {
			if ctx.params[k] != v {
				t.Fatalf("Route(%q): param %q = %q, want %q", tt.path, k, ctx.params[k], v)
			}
		}
		if len(tt.wantParam) == 0 && ctx.params["model"] != "" {
			t.Fatalf("Route(%q): expected static match, but :model param was set to %q", tt.path, ctx.params["model"])
		}
	}
}

// TestStaticAndParamSiblingsFromIssue reproduces the exact example from the issue body.
func TestStaticAndParamSiblingsFromIssue(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/dev/document-api/apiReference/model/videoModels", func(ctx *testContext) {})
	root.AddRoute("/dev/document-api/apiReference/model/imageModels", func(ctx *testContext) {})
	root.AddRoute("/dev/document-api/apiReference/model/:model", func(ctx *testContext) {})

	ctx := newTestCtx()
	h, _, _ := Route(&root, "/dev/document-api/apiReference/model/videoModels", ctx)
	if h == nil {
		t.Fatal("expected match for /dev/.../videoModels")
	}
	if ctx.params["model"] != "" {
		t.Fatalf("expected static match, got :model=%q", ctx.params["model"])
	}

	ctx = newTestCtx()
	h, _, _ = Route(&root, "/dev/document-api/apiReference/model/textModels", ctx)
	if h == nil {
		t.Fatal("expected match for /dev/.../textModels via :model param")
	}
	if ctx.params["model"] != "textModels" {
		t.Fatalf("expected :model=textModels, got %q", ctx.params["model"])
	}
}

// TestStaticAndParamSiblingsPanicDuplicateParam ensures two :param children at the same level still panic.
func TestStaticAndParamSiblingsPanicDuplicateParam(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate param wildcard at same level")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/model/:id", func(ctx *testContext) {})
	root.AddRoute("/model/:name", func(ctx *testContext) {}) // second param – should panic
}

// TestStaticAndParamSiblingsPanicCatchAllWithSiblings ensures catchAll still panics when siblings exist.
func TestStaticAndParamSiblingsPanicCatchAllWithSiblings(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for catchAll with existing children")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/files/readme.txt", func(ctx *testContext) {})
	root.AddRoute("/files/*filepath", func(ctx *testContext) {}) // catchAll with static sibling – should panic
}

// TestFindCaseInsensitivePathWithStaticParamSiblings verifies case-insensitive lookup
// with mixed static + param children.
func TestFindCaseInsensitivePathWithStaticParamSiblings(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/model/VideoModels", func(ctx *testContext) {})
	root.AddRoute("/model/:model", func(ctx *testContext) {})

	// Case-insensitive lookup should prefer the static child.
	fixedPath, found := root.FindCaseInsensitivePath("/model/videomodels", true)
	if !found {
		t.Fatal("expected case-insensitive match for /model/videomodels")
	}
	if fixedPath != "/model/VideoModels" {
		t.Fatalf("expected /model/VideoModels, got %q", fixedPath)
	}

	// Unknown value should fall through to :model param.
	fixedPath, found = root.FindCaseInsensitivePath("/model/someother", true)
	if !found {
		t.Fatal("expected match for /model/someother via :model param")
	}
	if fixedPath != "/model/someother" {
		t.Fatalf("expected /model/someother, got %q", fixedPath)
	}
}

// TestStaticAndParamSiblingsFirstByteCollision verifies that a param value whose
// first byte collides with a registered static sibling still routes correctly via
// the param child (regression test for the no-backtracking bug).
func TestStaticAndParamSiblingsFirstByteCollision(t *testing.T) {
	var root Node[func(*testContext)]
	root.AddRoute("/model/videoModels", func(ctx *testContext) {})
	root.AddRoute("/model/imageModels", func(ctx *testContext) {})
	root.AddRoute("/model/:model", func(ctx *testContext) {})

	tests := []struct {
		path      string
		wantMatch bool
		wantParam string // expected :model value; empty string means static match
	}{
		// Exact static matches must still work.
		{"/model/videoModels", true, ""},
		{"/model/imageModels", true, ""},
		// Param values that share the first byte with a static sibling.
		{"/model/videoModelsHD", true, "videoModelsHD"},
		{"/model/imageModels2", true, "imageModels2"},
		// Param values with no first-byte collision.
		{"/model/textModels", true, "textModels"},
		{"/model/anything", true, "anything"},
	}

	for _, tt := range tests {
		ctx := newTestCtx()
		h, _, _ := Route(&root, tt.path, ctx)
		if (h != nil) != tt.wantMatch {
			t.Fatalf("Route(%q): got match=%v, want match=%v", tt.path, h != nil, tt.wantMatch)
		}
		if tt.wantParam == "" {
			if ctx.params["model"] != "" {
				t.Fatalf("Route(%q): expected static match, but :model param = %q", tt.path, ctx.params["model"])
			}
		} else {
			if ctx.params["model"] != tt.wantParam {
				t.Fatalf("Route(%q): param model = %q, want %q", tt.path, ctx.params["model"], tt.wantParam)
			}
		}
	}
}

// TestStaticAndParamSiblingsPanicDuplicateParamAfterStatic ensures the new
// duplicate-param guard in insertChild fires when a static sibling is registered
// before the first param and a second param is then attempted.
func TestStaticAndParamSiblingsPanicDuplicateParamAfterStatic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate param wildcard at same level")
		}
	}()
	var root Node[func(*testContext)]
	root.AddRoute("/model/static", func(ctx *testContext) {}) // static sibling registered first
	root.AddRoute("/model/:id", func(ctx *testContext) {})    // first param — OK
	root.AddRoute("/model/:name", func(ctx *testContext) {})  // second param — must panic
}
