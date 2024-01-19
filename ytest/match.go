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
	"log"
	"net/url"
	"strconv"
)

// -----------------------------------------------------------------------------

func JsonEncode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Panicln("json.Marshal failed:", err)
	}
	return string(b)
}

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
	log.Panicf("formVal unexpected type: %T\n", val)
	return nil
}

// -----------------------------------------------------------------------------

func Equal__0[T basetype](a, b T) bool {
	return a == b
}

func Equal__1(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for key, got := range a {
		if expected, ok := b[key]; !ok || !Equal__4(got, expected) {
			return false
		}
	}
	return true
}

func Equal__2(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i, got := range a {
		if !Equal__4(got, b[i]) {
			return false
		}
	}
	return true
}

func Equal__3[T baseelem](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, got := range a {
		if !Equal__0(got, b[i]) {
			return false
		}
	}
	return true
}

func Equal__4(got, expected any) bool {
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
			return Equal__1(gv, ev)
		case *Var__1[map[string]any]:
			return Equal__1(gv, ev.Val())
		}
	case []any:
		switch ev := expected.(type) {
		case []any:
			return Equal__2(gv, ev)
		case *Var__2[[]any]:
			return Equal__2(gv, ev.Val())
		}
	case []string:
		switch ev := expected.(type) {
		case []string:
			return Equal__3(gv, ev)
		case *Var__3[[]string]:
			return Equal__3(gv, ev.Val())
		}
	case *Var__0[string]:
		switch ev := expected.(type) {
		case string:
			return gv.Equal(ev)
		case *Var__0[string]:
			return gv.Equal(ev.Val())
		}
	case *Var__0[int]:
		switch ev := expected.(type) {
		case int:
			return gv.Equal(ev)
		case *Var__0[int]:
			return gv.Equal(ev.Val())
		}
	case *Var__0[bool]:
		switch ev := expected.(type) {
		case bool:
			return gv.Equal(ev)
		case *Var__0[bool]:
			return gv.Equal(ev.Val())
		}
	case *Var__0[float64]:
		switch ev := expected.(type) {
		case float64:
			return gv.Equal(ev)
		case *Var__0[float64]:
			return gv.Equal(ev.Val())
		}
	case *Var__1[map[string]any]:
		switch ev := expected.(type) {
		case map[string]any:
			return gv.Equal(ev)
		case *Var__1[map[string]any]:
			return gv.Equal(ev.Val())
		}
	case *Var__2[[]any]:
		switch ev := expected.(type) {
		case []any:
			return gv.Equal(ev)
		case *Var__2[[]any]:
			return gv.Equal(ev.Val())
		}
	case *Var__3[[]string]:
		switch ev := expected.(type) {
		case []string:
			return gv.Equal(ev)
		case *Var__3[[]string]:
			return gv.Equal(ev.Val())
		}
	}
	log.Panicf("unsupported type to compare: %T == %T\n", got, expected)
	return false
}

// -----------------------------------------------------------------------------

func Match__0[T basetype](got, expected T) {
	if got != expected {
		log.Panicf("unmatched value! got: %v, expected: %v\n", got, expected)
	}
}

func Match__1(got, expected map[string]any) {
	for key, gv := range got {
		Match__4(gv, expected[key])
	}
}

func Match__2(got, expected []any) {
	if len(got) != len(expected) {
		log.Panicf("unmatched slice length! got: %d, expected: %d\n", len(got), len(expected))
	}
	for i, gv := range got {
		Match__4(gv, expected[i])
	}
}

func Match__3[T baseelem](got, expected []T) {
	if len(got) != len(expected) {
		log.Panicf("unmatched slice length! got: %d, expected: %d\n", len(got), len(expected))
	}
	for i, gv := range got {
		Match__0(gv, expected[i])
	}
}

func Match__4(got, expected any) {
	switch gv := got.(type) {
	case string:
		switch ev := expected.(type) {
		case string:
			Match__0(gv, ev)
			return
		case *Var__0[string]:
			Match__0(gv, ev.Val())
			return
		}
	case int:
		switch ev := expected.(type) {
		case int:
			Match__0(gv, ev)
			return
		case *Var__0[int]:
			Match__0(gv, ev.Val())
			return
		}
	case bool:
		switch ev := expected.(type) {
		case bool:
			Match__0(gv, ev)
			return
		case *Var__0[bool]:
			Match__0(gv, ev.Val())
			return
		}
	case float64:
		switch ev := expected.(type) {
		case float64:
			Match__0(gv, ev)
			return
		case *Var__0[float64]:
			Match__0(gv, ev.Val())
			return
		}
	case map[string]any:
		switch ev := expected.(type) {
		case map[string]any:
			Match__1(gv, ev)
			return
		case *Var__1[map[string]any]:
			Match__1(gv, ev.Val())
			return
		}
	case []any:
		switch ev := expected.(type) {
		case []any:
			Match__2(gv, ev)
			return
		case *Var__2[[]any]:
			Match__2(gv, ev.Val())
			return
		}
	case []string:
		switch ev := expected.(type) {
		case []string:
			Match__3(gv, ev)
			return
		case *Var__3[[]string]:
			Match__3(gv, ev.Val())
			return
		}
	case *Var__0[string]:
		switch ev := expected.(type) {
		case string:
			gv.Match(ev)
			return
		case *Var__0[string]:
			gv.Match(ev.Val())
			return
		}
	case *Var__0[int]:
		switch ev := expected.(type) {
		case int:
			gv.Match(ev)
			return
		case *Var__0[int]:
			gv.Match(ev.Val())
			return
		}
	case *Var__0[bool]:
		switch ev := expected.(type) {
		case bool:
			gv.Match(ev)
			return
		case *Var__0[bool]:
			gv.Match(ev.Val())
			return
		}
	case *Var__0[float64]:
		switch ev := expected.(type) {
		case float64:
			gv.Match(ev)
			return
		case *Var__0[float64]:
			gv.Match(ev.Val())
			return
		}
	case *Var__1[map[string]any]:
		switch ev := expected.(type) {
		case map[string]any:
			gv.Match(ev)
			return
		case *Var__1[map[string]any]:
			gv.Match(ev.Val())
			return
		}
	case *Var__2[[]any]:
		switch ev := expected.(type) {
		case []any:
			gv.Match(ev)
			return
		case *Var__2[[]any]:
			gv.Match(ev.Val())
			return
		}
	case *Var__3[[]string]:
		switch ev := expected.(type) {
		case []string:
			gv.Match(ev)
			return
		case *Var__3[[]string]:
			gv.Match(ev.Val())
			return
		}
	}
	log.Panicf("unmatched type! got: %T, expected: %T\n", got, expected)
}

// -----------------------------------------------------------------------------

type baseelem interface {
	string
}

type baseslice interface {
	[]string
}

type basetype interface {
	string | int | bool | float64
}

type Var__0[T basetype] struct {
	val   T
	valid bool
}

func (p *Var__0[T]) check() {
	if !p.valid {
		log.Panicln("read variable value before initialization")
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

func (p *Var__0[T]) Equal(v T) bool {
	p.check()
	return p.val == v
}

func (p *Var__0[T]) Match(v T) {
	if !p.valid {
		p.val, p.valid = v, true
		return
	}
	Match__0(p.val, v)
}

// -----------------------------------------------------------------------------

type Var__1[T map[string]any] struct {
	val T
}

func (p *Var__1[T]) check() {
	if p.val == nil {
		log.Panicln("read variable value before initialization")
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

func (p *Var__1[T]) Equal(v T) bool {
	p.check()
	return Equal__1(p.val, v)
}

func (p *Var__1[T]) Match(v T) {
	if p.val == nil {
		p.val = v
		return
	}
	Match__1(p.val, v)
}

// -----------------------------------------------------------------------------

type Var__2[T []any] struct {
	val   T
	valid bool
}

func (p *Var__2[T]) check() {
	if !p.valid {
		log.Panicln("read variable value before initialization")
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

func (p *Var__2[T]) Equal(v T) bool {
	p.check()
	return Equal__2(p.val, v)
}

func (p *Var__2[T]) Match(v T) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	Match__2(p.val, v)
}

// -----------------------------------------------------------------------------

type Var__3[T baseslice] struct {
	val   T
	valid bool
}

func (p *Var__3[T]) check() {
	if !p.valid {
		log.Panicln("read variable value before initialization")
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

func (p *Var__3[T]) Equal(v T) bool {
	p.check()
	return Equal__3(p.val, v)
}

func (p *Var__3[T]) Match(v T) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	Match__3(p.val, v)
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
