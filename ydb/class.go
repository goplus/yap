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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/goplus/gop/ast"
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
	name string
	tbl  string
	sql  *Sql
	db   *sql.DB
	wrap func(string, error) error
	apis map[string]*api

	query *query // query

	result []reflect.Value // result of an api call

	ret   func(args ...any) error
	onErr func(err error)

	test.Case
}

func newClass(name string, sql *Sql) *Class {
	if sql.db == nil {
		log.Panicln("please call `engine <sqlDriverName>` first")
	}
	return &Class{
		name: name,
		sql:  sql,
		db:   sql.db,
		wrap: sql.driver.WrapErr,
		apis: make(map[string]*api),
	}
}

func (p *Class) t() test.CaseT {
	if p.CaseT == nil {
		p.CaseT = logt.New()
	}
	return p.CaseT
}

func (p *Class) gen(ctx context.Context) {
}

// Use sets the default table used in following sql operations.
func (p *Class) Use(table string, src ...ast.Expr) {
	_, ok := p.sql.tables[table]
	if !ok {
		log.Panicln("table not found:", table)
	}
	p.tbl = table
}

// OnErr sets error processing of a sql execution.
func (p *Class) OnErr(onErr func(error), src ...ast.Expr) {
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
func (p *Class) Ret__0(src ast.Expr, args ...any) {
	if p.ret == nil {
		log.Panicln("please call `ret` after a `query` or `call` statement")
	}
	if src == nil {
		if len(args) == 0 {
			p.ret(nil)
			return
		}
		args = append(make([]any, 1, len(args)+1), args...)
	}
	p.ret(args...)
}

// Ret checks a query or call result.
func (p *Class) Ret__1(args ...any) {
	p.ret(args...)
}

// -----------------------------------------------------------------------------

type classApi struct {
	cls *Class
	api *api
}

type api struct {
	name string
	spec any
}

// Api creates a new api by a spec.
func (p *Class) Api(name string, spec any, src ...ast.Expr) func(args ...any) {
	api := &api{name: name, spec: spec}
	p.apis[name] = api
	return classApi{p, api}.call
}

var (
	tyError = reflect.TypeOf((*error)(nil)).Elem()
)

func setRetErr(result []reflect.Value, errRet error) {
	if n := len(result); n > 0 {
		if result[n-1].Type() == tyError {
			result[n-1] = reflect.ValueOf(errRet)
		}
	}
}

func (ca classApi) call(args ...any) {
	if debugExec {
		log.Println("==>", ca.api.name, args)
	}

	vArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vArgs[i] = reflect.ValueOf(arg)
	}
	vFn := reflect.ValueOf(ca.api.spec)

	var p = ca.cls
	var old = p.onErr
	var errRet error
	p.onErr = func(err error) {
		errRet = err
		panic(err)
	}
	defer func() {
		p.onErr = old
		if e := recover(); e != nil {
			if p.result == nil { // set p.result to zero if panic
				fnt := vFn.Type()
				n := fnt.NumOut()
				p.result = make([]reflect.Value, n)
				for i := 0; i < n; i++ {
					p.result[i] = reflect.Zero(fnt.Out(i))
				}
			}
			if errRet == nil {
				errRet = fmt.Errorf("%v", e)
			}
			setRetErr(p.result, errRet)
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
		t.Fatalf(
			"call ret: unmatched result parameters count - got %d, expected %d\n",
			len(args), len(result),
		)
	}
	for i, arg := range args {
		ret := result[i].Interface()
		test.Gopt_Case_MatchAny(t, arg, ret)
	}
	p.ret = nil
	return nil
}

// Out returns the ith reuslt.
func (p *Class) Out(i int) any {
	return p.result[i].Interface()
}

// -----------------------------------------------------------------------------
