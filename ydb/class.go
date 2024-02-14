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

package ydb

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"runtime/debug"

	"github.com/goplus/yap/test"
	"github.com/goplus/yap/test/logt"
)

var (
	ErrNoRows     = sql.ErrNoRows
	ErrDuplicated = errors.New("duplicated")
	ErrOutOfLimit = errors.New("out of limit")
)

// -----------------------------------------------------------------------------

type Class struct {
	Sql

	self reflect.Value
	tbl  string

	result []reflect.Value // result of an api call

	ret     func(args ...any) error
	onErr   func(err error)
	lastErr error
	test.Case

	query *query // query
}

func (p *Class) initClass(self any) {
	p.initSql()
	p.self = reflect.ValueOf(self)
}

func (p *Class) t() test.CaseT {
	if p.CaseT == nil {
		p.CaseT = logt.New()
	}
	return p.CaseT
}

// Use sets the default table used in following sql operations.
func (p *Class) Use(table string) {
	_, ok := p.tables[table]
	if !ok {
		log.Panicln("table not found:", table)
	}
	p.tbl = table
}

// OnErr sets error processing of a sql execution.
func (p *Class) OnErr(onErr func(error)) {
	p.onErr = onErr
}

func (p *Class) handleErr(prompt string, err error) {
	err = p.wrap(prompt, err)
	if p.onErr == nil {
		log.Panicln(prompt, err)
	}
	p.onErr(err)
}

// Ret checks a query or call result.
//
// For checking query result:
//   - ret <expr1>, &<var1>, <expr2>, &<var2>, ...
//   - ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...
//   - ret &<structVar>
//   - ret &<structSlice>
//
// For checking call result:
//   - ret <expr1>, <expr2>, ...
func (p *Class) Ret(args ...any) {
	if p.ret == nil {
		log.Panicln("please call `ret` after a `query` or `call` statement")
	}
	p.ret(args...)
}

// -----------------------------------------------------------------------------

func (p *Class) Gop_Exec(name string, args ...any) {
	vFn := p.method(name)
	p.call(name, vFn, args...)
}

func (p *Class) method(name string) reflect.Value {
	c := name[0]
	if c >= 'a' && c <= 'z' {
		name = string(c-('a'-'A')) + name[1:]
	}
	name = "API_" + name
	return p.self.MethodByName(name)
}

func (p *Class) call(name string, vFn reflect.Value, args ...any) {
	if debugExec {
		log.Println("==>", name, args)
	}

	vArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vArgs[i] = reflect.ValueOf(arg)
	}

	var old = p.onErr
	var errRet error
	p.onErr = func(err error) {
		errRet = err
		panic(err)
	}
	defer func() {
		p.onErr = old
		if p.result == nil { // set p.result to zero if panic
			fnt := vFn.Type()
			n := fnt.NumOut()
			p.result = make([]reflect.Value, n, n+1)
			for i := 0; i < n; i++ {
				p.result[i] = reflect.Zero(fnt.Out(i))
			}
		}
		if !hasRetErrType(p.result) {
			p.result = append(p.result, reflect.Zero(tyError))
		}
		if e := recover(); e != nil {
			if errRet == nil {
				errRet = recoverErr(e)
			}
			p.result[len(p.result)-1] = reflect.ValueOf(errRet)
			if debugExec {
				log.Println("PANIC:", e)
				debug.PrintStack()
			}
		}
		p.ret = p.callRet
	}()
	p.result = nil
	p.result = vFn.Call(vArgs)
}

func (p *Class) callRet(args ...any) error {
	t := p.t()
	result := p.result
	if len(result) != len(args) {
		if len(result) != len(args)+1 {
			t.Fatalf(
				"call ret: unmatched result parameters count - got %d, expected %d\n",
				len(args), len(result),
			)
		}
		args = append(args, nil)
	}
	for i, arg := range args {
		ret := result[i].Interface()
		test.Gopt_Case_MatchAny(t, arg, ret)
	}
	p.ret = nil
	return nil
}

var (
	tyError = reflect.TypeOf((*error)(nil)).Elem()
)

func hasRetErrType(result []reflect.Value) bool {
	if n := len(result); n > 0 {
		return result[n-1].Type() == tyError
	}
	return false
}

// Out returns the ith reuslt.
func (p *Class) Out(i int) any {
	return p.result[i].Interface()
}

// -----------------------------------------------------------------------------
