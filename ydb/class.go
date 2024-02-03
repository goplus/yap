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
	"log"
	"reflect"

	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

type Class struct {
	name   string
	tbl    string
	apis   map[string]*dbApi
	api    *dbApi
	result []reflect.Value
	ret    func(args ...any)
}

func newClass(name string) *Class {
	apis := make(map[string]*dbApi)
	return &Class{name: name, apis: apis}
}

// Ret checks a query or call result.
func (p *Class) Ret(args ...any) {
	p.ret(args...)
}

// -----------------------------------------------------------------------------

func (p *Class) Use(table string) {
	p.tbl = table
}

func (p *Class) Insert(kvPair ...any) {
}

func (p *Class) queryRet(kvPair ...any) {
}

func (p *Class) Query(query string) {
	p.ret = p.queryRet
}

func (p *Class) Limit__0(n int) {
}

func (p *Class) Limit__1(n int, query string) {
}

// -----------------------------------------------------------------------------

type dbApi struct {
	name string
	spec any
}

func (p *Class) Api(name string, spec any, fnlit ...*ast.FuncLit) {
	api := &dbApi{name: name, spec: spec}
	p.api = api
	p.apis[name] = api
}

func (p *Class) Call(args ...any) {
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

func (p *Class) callRet(args ...any) {
}

// -----------------------------------------------------------------------------
