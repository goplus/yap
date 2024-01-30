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
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type basetype interface {
	string | int | bool | float64
}

type baseelem interface {
	string
}

type baseslice interface {
	[]string
}

// -----------------------------------------------------------------------------

// JsonEncode encodes a value into string in json format.
func JsonEncode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		fatal("json.Marshal failed:", err)
	}
	return string(b)
}

// Form encodes a map value into string in form format.
func Form(m map[string]any) (ret url.Values) {
	ret = make(url.Values)
	for k, v := range m {
		ret[k] = formVal(v)
	}
	return ret
}

func formVal(val any) []string {
	switch v := val.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case int:
		return []string{strconv.Itoa(v)}
	case bool:
		if v {
			return []string{"true"}
		}
		return []string{"false"}
	case float64:
		return []string{strconv.FormatFloat(v, 'g', -1, 64)}
	}
	fatalf("formVal unexpected type: %T\n", val)
	return nil
}

// -----------------------------------------------------------------------------

func Gopt_Case_Equal__0[T basetype](t CaseT, a, b T) bool {
	return a == b
}

func Gopt_Case_Equal__1(t CaseT, a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	t.Helper()
	for key, got := range a {
		if expected, ok := b[key]; !ok || !Gopt_Case_Equal__4(t, got, expected) {
			return false
		}
	}
	return true
}

func Gopt_Case_Equal__2(t CaseT, a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	t.Helper()
	for i, got := range a {
		if !Gopt_Case_Equal__4(t, got, b[i]) {
			return false
		}
	}
	return true
}

func Gopt_Case_Equal__3[T baseelem](t CaseT, a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	t.Helper()
	for i, got := range a {
		if !Gopt_Case_Equal__0(t, got, b[i]) {
			return false
		}
	}
	return true
}

func Gopt_Case_Equal__4(t CaseT, got, expected any) bool {
	t.Helper()
	switch gv := got.(type) {
	case string:
		switch ev := expected.(type) {
		case string:
			return gv == ev
		case *Var__0[string]:
			return gv == ev.Val()
		}
	case int:
		switch ev := expected.(type) {
		case int:
			return gv == ev
		case *Var__0[int]:
			return gv == ev.Val()
		}
	case bool:
		switch ev := expected.(type) {
		case bool:
			return gv == ev
		case *Var__0[bool]:
			return gv == ev.Val()
		}
	case float64:
		switch ev := expected.(type) {
		case float64:
			return gv == ev
		case *Var__0[float64]:
			return gv == ev.Val()
		}
	case map[string]any:
		switch ev := expected.(type) {
		case map[string]any:
			return Gopt_Case_Equal__1(t, gv, ev)
		case *Var__1[map[string]any]:
			return Gopt_Case_Equal__1(t, gv, ev.Val())
		}
	case []any:
		switch ev := expected.(type) {
		case []any:
			return Gopt_Case_Equal__2(t, gv, ev)
		case *Var__2[[]any]:
			return Gopt_Case_Equal__2(t, gv, ev.Val())
		}
	case []string:
		switch ev := expected.(type) {
		case []string:
			return Gopt_Case_Equal__3(t, gv, ev)
		case *Var__3[[]string]:
			return Gopt_Case_Equal__3(t, gv, ev.Val())
		}
	case *Var__0[string]:
		switch ev := expected.(type) {
		case string:
			return gv.Equal(t, ev)
		case *Var__0[string]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__0[int]:
		switch ev := expected.(type) {
		case int:
			return gv.Equal(t, ev)
		case *Var__0[int]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__0[bool]:
		switch ev := expected.(type) {
		case bool:
			return gv.Equal(t, ev)
		case *Var__0[bool]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__0[float64]:
		switch ev := expected.(type) {
		case float64:
			return gv.Equal(t, ev)
		case *Var__0[float64]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__1[map[string]any]:
		switch ev := expected.(type) {
		case map[string]any:
			return gv.Equal(t, ev)
		case *Var__1[map[string]any]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__2[[]any]:
		switch ev := expected.(type) {
		case []any:
			return gv.Equal(t, ev)
		case *Var__2[[]any]:
			return gv.Equal(t, ev.Val())
		}
	case *Var__3[[]string]:
		switch ev := expected.(type) {
		case []string:
			return gv.Equal(t, ev)
		case *Var__3[[]string]:
			return gv.Equal(t, ev.Val())
		}
	}
	t.Fatalf("unsupported type to compare: %T == %T\n", got, expected)
	return false
}

// -----------------------------------------------------------------------------

func Gopt_Case_Match__0[T basetype](t CaseT, got, expected T) {
	if got != expected {
		t.Helper()
		t.Fatalf("unmatched value - got: %v, expected: %v\n", got, expected)
	}
}

func Gopt_Case_Match__1(t CaseT, got, expected map[string]any) {
	t.Helper()
	for key, gv := range got {
		Gopt_Case_Match__4(t, gv, expected[key])
	}
}

func Gopt_Case_Match__2(t CaseT, got, expected []any) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("unmatched slice length - got: %d, expected: %d\n", len(got), len(expected))
	}
	for i, gv := range got {
		Gopt_Case_Match__4(t, gv, expected[i])
	}
}

func Gopt_Case_Match__3[T baseelem](t CaseT, got, expected []T) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("unmatched slice length - got: %d, expected: %d\n", len(got), len(expected))
	}
	for i, gv := range got {
		Gopt_Case_Match__0(t, gv, expected[i])
	}
}

func Gopt_Case_Match__4(t CaseT, got, expected any) {
	t.Helper()
	switch gv := got.(type) {
	case string:
		switch ev := expected.(type) {
		case string:
			Gopt_Case_Match__0(t, gv, ev)
			return
		case *Var__0[string]:
			Gopt_Case_Match__0(t, gv, ev.Val())
			return
		}
	case int:
		switch ev := expected.(type) {
		case int:
			Gopt_Case_Match__0(t, gv, ev)
			return
		case *Var__0[int]:
			Gopt_Case_Match__0(t, gv, ev.Val())
			return
		}
	case bool:
		switch ev := expected.(type) {
		case bool:
			Gopt_Case_Match__0(t, gv, ev)
			return
		case *Var__0[bool]:
			Gopt_Case_Match__0(t, gv, ev.Val())
			return
		}
	case float64:
		switch ev := expected.(type) {
		case float64:
			Gopt_Case_Match__0(t, gv, ev)
			return
		case *Var__0[float64]:
			Gopt_Case_Match__0(t, gv, ev.Val())
			return
		}
	case map[string]any:
		switch ev := expected.(type) {
		case map[string]any:
			Gopt_Case_Match__1(t, gv, ev)
			return
		case *Var__1[map[string]any]:
			Gopt_Case_Match__1(t, gv, ev.Val())
			return
		}
	case []any:
		switch ev := expected.(type) {
		case []any:
			Gopt_Case_Match__2(t, gv, ev)
			return
		case *Var__2[[]any]:
			Gopt_Case_Match__2(t, gv, ev.Val())
			return
		}
	case []string:
		switch ev := expected.(type) {
		case []string:
			Gopt_Case_Match__3(t, gv, ev)
			return
		case *Var__3[[]string]:
			Gopt_Case_Match__3(t, gv, ev.Val())
			return
		}
	case *Var__0[string]:
		switch ev := expected.(type) {
		case string:
			gv.Match(t, ev)
			return
		case *Var__0[string]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__0[int]:
		switch ev := expected.(type) {
		case int:
			gv.Match(t, ev)
			return
		case *Var__0[int]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__0[bool]:
		switch ev := expected.(type) {
		case bool:
			gv.Match(t, ev)
			return
		case *Var__0[bool]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__0[float64]:
		switch ev := expected.(type) {
		case float64:
			gv.Match(t, ev)
			return
		case *Var__0[float64]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__1[map[string]any]:
		switch ev := expected.(type) {
		case map[string]any:
			gv.Match(t, ev)
			return
		case *Var__1[map[string]any]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__2[[]any]:
		switch ev := expected.(type) {
		case []any:
			gv.Match(t, ev)
			return
		case *Var__2[[]any]:
			gv.Match(t, ev.Val())
			return
		}
	case *Var__3[[]string]:
		switch ev := expected.(type) {
		case []string:
			gv.Match(t, ev)
			return
		case *Var__3[[]string]:
			gv.Match(t, ev.Val())
			return
		}
	}
	t.Fatalf("unmatched type - got: %T, expected: %T\n", got, expected)
}

// -----------------------------------------------------------------------------

type Var__0[T basetype] struct {
	val   T
	valid bool
}

func (p *Var__0[T]) check() {
	if !p.valid {
		fatal("read variable value before initialization")
	}
}

func (p *Var__0[T]) Valid() bool {
	return p.valid
}

func (p *Var__0[T]) String() string {
	p.check()
	return fmt.Sprint(p.val) // TODO: optimization
}

func (p *Var__0[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__0[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__0[T]) Equal(t CaseT, v T) bool {
	p.check()
	return p.val == v
}

func (p *Var__0[T]) Match(t CaseT, v T) {
	if !p.valid {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_Match__0(t, p.val, v)
}

// -----------------------------------------------------------------------------

type Var__1[T map[string]any] struct {
	val T
}

func (p *Var__1[T]) check() {
	if p.val == nil {
		fatal("read variable value before initialization")
	}
}

func (p *Var__1[T]) Valid() bool {
	return p.val != nil
}

func (p *Var__1[T]) Form() string {
	p.check()
	return Form(p.val).Encode()
}

func (p *Var__1[T]) Json() string {
	p.check()
	return JsonEncode(p.val)
}

func (p *Var__1[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__1[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__1[T]) Equal(t CaseT, v T) bool {
	p.check()
	t.Helper()
	return Gopt_Case_Equal__1(t, p.val, v)
}

func (p *Var__1[T]) Match(t CaseT, v T) {
	if p.val == nil {
		p.val = v
		return
	}
	t.Helper()
	Gopt_Case_Match__1(t, p.val, v)
}

// -----------------------------------------------------------------------------

type Var__2[T []any] struct {
	val   T
	valid bool
}

func (p *Var__2[T]) check() {
	if !p.valid {
		fatal("read variable value before initialization")
	}
}

func (p *Var__2[T]) Valid() bool {
	return p.valid
}

func (p *Var__2[T]) Json() string {
	p.check()
	return JsonEncode(p.val)
}

func (p *Var__2[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__2[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__2[T]) Equal(t CaseT, v T) bool {
	p.check()
	t.Helper()
	return Gopt_Case_Equal__2(t, p.val, v)
}

func (p *Var__2[T]) Match(t CaseT, v T) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_Match__2(t, p.val, v)
}

// -----------------------------------------------------------------------------

type Var__3[T baseslice] struct {
	val   T
	valid bool
}

func (p *Var__3[T]) check() {
	if !p.valid {
		fatal("read variable value before initialization")
	}
}

func (p *Var__3[T]) Valid() bool {
	return p.valid
}

func (p *Var__3[T]) Json() string {
	p.check()
	return JsonEncode(p.val)
}

func (p *Var__3[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__3[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__3[T]) Equal(t CaseT, v T) bool {
	p.check()
	t.Helper()
	return Gopt_Case_Equal__3(t, p.val, v)
}

func (p *Var__3[T]) Match(t CaseT, v T) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_Match__3(t, p.val, v)
}

// -----------------------------------------------------------------------------

func Gopx_Var_Cast__0[T basetype]() *Var__0[T] {
	return new(Var__0[T])
}

func Gopx_Var_Cast__1[T map[string]any]() *Var__1[T] {
	return new(Var__1[T])
}

func Gopx_Var_Cast__2[T []any]() *Var__2[T] {
	return new(Var__2[T])
}

func Gopx_Var_Cast__3[T baseslice]() *Var__3[T] {
	return new(Var__3[T])
}

// -----------------------------------------------------------------------------
