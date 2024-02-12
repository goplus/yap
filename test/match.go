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

package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	GopPackage = true
)

type basetype interface {
	string | int | bool | float64
}

func toMapAny[T basetype](val map[string]T) map[string]any {
	ret := make(map[string]any, len(val))
	for k, v := range val {
		ret[k] = v
	}
	return ret
}

// -----------------------------------------------------------------------------

type baseelem interface {
	string
}

type baseslice interface {
	[]string
}

type TySet[T baseelem] []T

func Set__0[T baseelem](vals ...T) TySet[T] {
	return TySet[T](vals)
}

func Set__1[T []string](v *Var__3[T]) TySet[string] {
	return TySet[string](v.Val())
}

// -----------------------------------------------------------------------------

type Case struct {
	CaseT
}

func nameCtx(name []string) string {
	if name != nil {
		return " (" + strings.Join(name, ".") + ")"
	}
	return ""
}

const (
	Gopo_Gopt_Case_Match = "Gopt_Case_MatchTBase,Gopt_Case_MatchMap,Gopt_Case_MatchSlice,Gopt_Case_MatchBaseSlice,Gopt_Case_MatchSet,Gopt_Case_MatchAny"
)

func Gopt_Case_MatchTBase[T basetype](t CaseT, got, expected T, name ...string) {
	if got != expected {
		t.Helper()
		t.Fatalf("unmatched value%s - got: %v, expected: %v\n", nameCtx(name), got, expected)
	}
}

func Gopt_Case_MatchMap(t CaseT, got, expected map[string]any, name ...string) {
	t.Helper()
	idx := len(name)
	name = append(name, "")
	for key, gv := range got {
		name[idx] = key
		Gopt_Case_MatchAny(t, gv, expected[key], name...)
	}
}

func Gopt_Case_MatchSlice(t CaseT, got, expected []any, name ...string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("unmatched slice%s length - got: %d, expected: %d\n", nameCtx(name), len(got), len(expected))
	}
	idx := len(name) - 1
	if idx < 0 {
		idx, name = 0, []string{"$"}
	}
	for i, gv := range got {
		name[idx] = "[" + strconv.Itoa(i) + "]"
		Gopt_Case_MatchAny(t, gv, expected[i], name...)
	}
}

func Gopt_Case_MatchBaseSlice[T baseelem](t CaseT, got, expected []T, name ...string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("unmatched slice%s length - got: %d, expected: %d\n", nameCtx(name), len(got), len(expected))
	}
	idx := len(name) - 1
	if idx < 0 {
		idx, name = 0, []string{"$"}
	}
	for i, gv := range got {
		name[idx] = "[" + strconv.Itoa(i) + "]"
		Gopt_Case_MatchTBase(t, gv, expected[i], name...)
	}
}

func Gopt_Case_MatchSet[T baseelem](t CaseT, got []T, expected TySet[T], name ...string) {
	if len(got) != len(expected) {
		t.Fatalf("unmatched set%s length - got: %d, expected: %d\n", nameCtx(name), len(got), len(expected))
	}
	for _, ev := range expected {
		if !hasElem(ev, got) {
			t.Fatalf("unmatched set%s: got: %v, value %v doesn't exist in it\n", nameCtx(name), got, ev)
		}
	}
}

func hasElem[T baseelem](v T, got []T) bool {
	for _, gv := range got {
		if v == gv {
			return true
		}
	}
	return false
}

func Gopt_Case_MatchAny(t CaseT, got, expected any, name ...string) {
	t.Helper()
retry:
	switch gv := got.(type) {
	case string:
		switch ev := expected.(type) {
		case string:
			Gopt_Case_MatchTBase(t, gv, ev, name...)
			return
		case *Var__0[string]:
			Gopt_Case_MatchTBase(t, gv, ev.Val(), name...)
			return
		}
	case int:
		switch ev := expected.(type) {
		case int:
			Gopt_Case_MatchTBase(t, gv, ev, name...)
			return
		case *Var__0[int]:
			Gopt_Case_MatchTBase(t, gv, ev.Val(), name...)
			return
		}
	case bool:
		switch ev := expected.(type) {
		case bool:
			Gopt_Case_MatchTBase(t, gv, ev, name...)
			return
		case *Var__0[bool]:
			Gopt_Case_MatchTBase(t, gv, ev.Val(), name...)
			return
		}
	case float64:
		switch ev := expected.(type) {
		case float64:
			Gopt_Case_MatchTBase(t, gv, ev, name...)
			return
		case *Var__0[float64]:
			Gopt_Case_MatchTBase(t, gv, ev.Val(), name...)
			return
		}
	case map[string]any:
		switch ev := expected.(type) {
		case map[string]any:
			Gopt_Case_MatchMap(t, gv, ev, name...)
			return
		case *Var__1[map[string]any]:
			Gopt_Case_MatchMap(t, gv, ev.Val(), name...)
			return
		}
	case []any:
		switch ev := expected.(type) {
		case []any:
			Gopt_Case_MatchSlice(t, gv, ev, name...)
			return
		case *Var__2[[]any]:
			Gopt_Case_MatchSlice(t, gv, ev.Val(), name...)
			return
		}
	case []string:
		switch ev := expected.(type) {
		case []string:
			Gopt_Case_MatchBaseSlice(t, gv, ev, name...)
			return
		case TySet[string]:
			Gopt_Case_MatchSet(t, gv, ev, name...)
			return
		case *Var__3[[]string]:
			Gopt_Case_MatchBaseSlice(t, gv, ev.Val(), name...)
			return
		}
	case *Var__0[string]:
		switch ev := expected.(type) {
		case string:
			gv.Match(t, ev, name...)
			return
		case *Var__0[string]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__0[int]:
		switch ev := expected.(type) {
		case int:
			gv.Match(t, ev, name...)
			return
		case *Var__0[int]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__0[bool]:
		switch ev := expected.(type) {
		case bool:
			gv.Match(t, ev, name...)
			return
		case *Var__0[bool]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__0[float64]:
		switch ev := expected.(type) {
		case float64:
			gv.Match(t, ev, name...)
			return
		case *Var__0[float64]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__1[map[string]any]:
		switch ev := expected.(type) {
		case map[string]any:
			gv.Match(t, ev, name...)
			return
		case *Var__1[map[string]any]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__2[[]any]:
		switch ev := expected.(type) {
		case []any:
			gv.Match(t, ev, name...)
			return
		case *Var__2[[]any]:
			gv.Match(t, ev.Val(), name...)
			return
		}
	case *Var__3[[]string]:
		switch ev := expected.(type) {
		case []string:
			gv.Match__0(t, ev, name...)
			return
		case TySet[string]:
			gv.Match__1(t, ev, name...)
			return
		case *Var__3[[]string]:
			gv.Match__0(t, ev.Val(), name...)
			return
		}

	// fallback types:
	case map[string]string:
		got = toMapAny(gv)
		goto retry
	case map[string]int:
		got = toMapAny(gv)
		goto retry
	case map[string]bool:
		got = toMapAny(gv)
		goto retry
	case map[string]float64:
		got = toMapAny(gv)
		goto retry

	// other types:
	default:
		if reflect.DeepEqual(got, expected) {
			return
		}
	}
	t.Fatalf(
		"unmatched%s - got: %v (%T), expected: %v (%T)\n",
		nameCtx(name), got, got, expected, expected,
	)
}

// -----------------------------------------------------------------------------

type Var__0[T basetype] struct {
	val   T
	valid bool
}

func (p *Var__0[T]) check() {
	if !p.valid {
		Fatal("read variable value before initialization")
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

func (p *Var__0[T]) UnmarshalJSON(data []byte) error {
	p.valid = true
	return json.Unmarshal(data, &p.val)
}

func (p *Var__0[T]) Equal(t CaseT, v T) bool {
	p.check()
	return p.val == v
}

func (p *Var__0[T]) Match(t CaseT, v T, name ...string) {
	if !p.valid {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_MatchTBase(t, p.val, v, name...)
}

// -----------------------------------------------------------------------------

type Var__1[T map[string]any] struct {
	val T
}

func (p *Var__1[T]) check() {
	if p.val == nil {
		Fatal("read variable value before initialization")
	}
}

func (p *Var__1[T]) Valid() bool {
	return p.val != nil
}

func (p *Var__1[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__1[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__1[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.val)
}

func (p *Var__1[T]) Match(t CaseT, v T, name ...string) {
	if p.val == nil {
		p.val = v
		return
	}
	t.Helper()
	Gopt_Case_MatchMap(t, p.val, v, name...)
}

// -----------------------------------------------------------------------------

type Var__2[T []any] struct {
	val   T
	valid bool
}

func (p *Var__2[T]) check() {
	if !p.valid {
		Fatal("read variable value before initialization")
	}
}

func (p *Var__2[T]) Valid() bool {
	return p.valid
}

func (p *Var__2[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__2[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__2[T]) UnmarshalJSON(data []byte) error {
	p.valid = true
	return json.Unmarshal(data, &p.val)
}

func (p *Var__2[T]) Match(t CaseT, v T, name ...string) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_MatchSlice(t, p.val, v, name...)
}

// -----------------------------------------------------------------------------

type Var__3[T baseslice] struct {
	val   T
	valid bool
}

func (p *Var__3[T]) check() {
	if !p.valid {
		Fatal("read variable value before initialization")
	}
}

func (p *Var__3[T]) Valid() bool {
	return p.valid
}

func (p *Var__3[T]) Val() T {
	p.check()
	return p.val
}

func (p *Var__3[T]) MarshalJSON() ([]byte, error) {
	p.check()
	return json.Marshal(p.val)
}

func (p *Var__3[T]) UnmarshalJSON(data []byte) error {
	p.valid = true
	return json.Unmarshal(data, &p.val)
}

func (p *Var__3[T]) Match__0(t CaseT, v T, name ...string) {
	if p.val == nil {
		p.val, p.valid = v, true
		return
	}
	t.Helper()
	Gopt_Case_MatchBaseSlice(t, p.val, v, name...)
}

func (p *Var__3[T]) Match__1(t CaseT, v TySet[string], name ...string) {
	if p.val == nil {
		p.val, p.valid = T(v), true
		return
	}
	t.Helper()
	Gopt_Case_MatchSet(t, p.val, v, name...)
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

func Gopx_Var_Cast__3[T []string]() *Var__3[T] {
	return new(Var__3[T])
}

// -----------------------------------------------------------------------------
