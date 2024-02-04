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
	"log"
	"reflect"

	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

type Class struct {
	name   string
	tbl    string
	apis   map[string]*api
	api    *api
	result []reflect.Value
	ret    func(args ...any)
}

func newClass(name string) *Class {
	apis := make(map[string]*api)
	return &Class{name: name, apis: apis}
}

func (p *Class) create(ctx context.Context, sql *Sql) {
}

// Ret checks a query or call result.
func (p *Class) Ret__0(src ast.Node, args ...any) {
	p.ret(args...)
}

// Ret checks a query or call result.
func (p *Class) Ret__1(args ...any) {
	p.ret(args...)
}

// -----------------------------------------------------------------------------

// Use sets the default table used in following sql operations.
func (p *Class) Use(table string, src ...ast.Node) {
	p.tbl = table
}

// Insert inserts a new row.
func (p *Class) Insert__0(src ast.Node, kvPair ...any) {
}

// Insert inserts a new row.
func (p *Class) Insert__1(kvPair ...any) {
}

func (p *Class) queryRet(kvPair ...any) {
}

// Query creates a new query.
func (p *Class) Query(query string, src ...ast.Node) {
	p.ret = p.queryRet
}

// Limit sets query result rows limit.
func (p *Class) Limit__0(n int, src ...ast.Node) {
}

// Limit checks if query result rows is < n or not.
func (p *Class) Limit__1(n int, query string, src ...ast.Node) {
}

// -----------------------------------------------------------------------------

type api struct {
	name string
	spec any
}

// Api creates a new api by a spec.
func (p *Class) Api(name string, spec any, src ...*ast.FuncLit) {
	api := &api{name: name, spec: spec}
	p.api = api
	p.apis[name] = api
}

// Call calls an api with specified args.
func (p *Class) Call__0(src ast.Node, args ...any) {
	if p.api == nil {
		log.Panicln("please call after an api definition")
	}
	vArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vArgs[i] = reflect.ValueOf(arg)
	}
	p.result = reflect.ValueOf(p.api).Call(vArgs)
	p.ret = p.callRet
}

// Call calls an api with specified args.
func (p *Class) Call__1(args ...any) {
	p.Call__0(nil, args...)
}

func (p *Class) callRet(args ...any) {
}

// -----------------------------------------------------------------------------
