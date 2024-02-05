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
	"log"
	"reflect"
	"strings"

	"github.com/goplus/gop/ast"
)

var (
	ErrDuplicated = errors.New("duplicated")
)

// -----------------------------------------------------------------------------

type Class struct {
	name string
	tbl  string
	sql  *Sql
	db   *sql.DB
	apis map[string]*api

	query *query // query

	api    *api
	result []reflect.Value // result of an api call

	ret func(args ...any)
}

func newClass(name string, sql *Sql) *Class {
	if sql.db == nil {
		log.Panicln("please call `engine <sqlDriverName>` first")
	}
	return &Class{
		name: name,
		sql:  sql,
		db:   sql.db,
		apis: make(map[string]*api),
	}
}

func (p *Class) gen(ctx context.Context) {
}

// Use sets the default table used in following sql operations.
func (p *Class) Use(table string, src ...ast.Node) {
	if _, ok := p.sql.tables[table]; !ok {
		log.Panicln("table not found:", table)
	}
	p.tbl = table
}

// Ret checks a query or call result.
//
// For checking query result:
//   - ret <colName1>, &<var1>, <colName2>, &<var2>, ...
//   - ret <colName1>, &<varSlice1>, <colName2>, &<varSlice2>, ...
//   - ret &<structVar>
//   - ret &<structSlice>
//
// For checking call result:
//   - TODO
func (p *Class) Ret__0(src ast.Node, args ...any) {
	if p.ret == nil {
		log.Panicln("please call `ret` after a `query` or `call` statement")
	}
	p.ret(args...)
}

// Ret checks a query or call result.
func (p *Class) Ret__1(args ...any) {
	p.Ret__0(nil, args...)
}

// -----------------------------------------------------------------------------

// Insert inserts a new row.
//   - insert <colName1>, <val1>, <colName2>, <val2>, ...
//   - insert <colName1>, <valSlice1>, <colName2>, <valSlice2>, ...
//   - insert <structValOrPtr>
//   - insert <structSlice>
func (p *Class) Insert__0(src ast.Node, args ...any) {
	/* if p.tbl == "" {
		TODO:
	} */
	nArg := len(args)
	if nArg == 1 {
		p.insertStruc(args[0])
	} else {
		p.insertKvPair(args...)
	}
}

// Insert inserts a new row.
//   - insert <structValOrPtr>
//   - insert <structSlice>
func (p *Class) insertStruc(arg any) {
}

const (
	valFlagNormal = 1
	valFlagSlice  = 2
)

// Insert inserts a new row.
//   - insert <colName1>, <val1>, <colName2>, <val2>, ...
//   - insert <colName1>, <valSlice1>, <colName2>, <valSlice2>, ...
func (p *Class) insertKvPair(kvPair ...any) {
	nPair := len(kvPair)
	if nPair < 2 || nPair&1 != 0 {
		log.Panicln("usage: insert <colName1>, <val1>, <colName2>, <val2>, ...")
	}
	n := nPair >> 1
	names := make([]string, n)
	vals := make([]any, n)
	rows := -1
	kind := 0
	for i := 0; i < nPair; i += 2 {
		val := kvPair[i+1]
		names[i>>1] = kvPair[i].(string)
		vals[i>>1] = val
		switch v := reflect.ValueOf(val); v.Kind() {
		case reflect.Slice:
			vals = append(vals, v)
			if kind == 0 {
				rows = v.Len()
				kind = valFlagSlice
			} else if vlen := v.Len(); rows != vlen {
				log.Panicf("insert: unexpected slice length. got %d, expected %d\n", vlen, rows)
			}
		default:
			vals = append(vals, val)
			kind |= valFlagNormal
		}
	}
	switch kind {
	case valFlagNormal:
		p.insertRow(names, vals)
	case valFlagSlice:
		p.insertRows(names, vals, rows)
	default:
		log.Panicln("can't insert mix slice and normal value")
	}
}

func (p *Class) insertRows(names []string, args []any, rows int) {
	n := len(args)
	valparam := valParam(n)
	valparams := strings.Repeat(valparam+",", rows)
	valparams = valparams[:len(valparams)-1]

	query := insertQuery(p.tbl, names)
	query = append(query, valparams...)

	vals := make([]any, 0, n*rows)
	for i := 0; i < rows; i++ {
		for _, arg := range args {
			v := arg.(reflect.Value)
			vals = append(vals, v.Index(i).Interface())
		}
	}
	result, err := p.db.ExecContext(context.TODO(), string(query), vals...)
	insertRet(result, err)
}

func (p *Class) insertRow(names []string, vals []any) {
	query := insertQuery(p.tbl, names)
	query = append(query, valParam(len(vals))...)
	result, err := p.db.ExecContext(context.TODO(), string(query), vals...)
	insertRet(result, err)
}

func insertRet(result sql.Result, err error) {
	if err != nil {
		log.Panicln("insert:", err)
	}
}

func insertQuery(tbl string, names []string) []byte {
	query := make([]byte, 0, 128)
	query = append(query, "INSERT INTO "...)
	query = append(query, tbl...)
	query = append(query, ' ', '(')
	query = append(query, strings.Join(names, ",")...)
	query = append(query, ") VALUES "...)
	return query
}

func valParam(n int) string {
	valparam := strings.Repeat("?,", n)
	valparam = "(" + valparam[:len(valparam)-1] + ")"
	return valparam
}

// Insert inserts a new row.
func (p *Class) Insert__1(kvPair ...any) {
	p.Insert__0(nil, kvPair...)
}

// Count returns rows of a query result.
func (p *Class) Count__0(src ast.Node, cond string, args ...any) (n int) {
	if p.tbl == "" {
		log.Panicln("please call `use <tableName>` to specified a table name")
	}
	row := p.db.QueryRowContext(context.TODO(), "SELECT COUNT(*) FROM "+p.tbl+" WHERE "+cond, args...)
	if err := row.Scan(&n); err != nil {
		log.Panicln("count:", err)
	}
	return
}

// Count returns rows of a query result.
func (p *Class) Count__1(cond string, args ...any) (n int) {
	return p.Count__0(nil, cond, args...)
}

// -----------------------------------------------------------------------------

type query struct {
	cond  string // where
	args  []any  // one of query argument <argN> can be a slice
	limit int    // 0 means no limit
}

// For checking query result:
//   - ret <colName1>, &<var1>, <colName2>, &<var2>, ...
//   - ret <colName1>, &<varSlice1>, <colName2>, &<varSlice2>, ...
//   - ret &<structVar>
//   - ret &<structSlice>
func (p *Class) queryRet(args ...any) {
	nArg := len(args)
	if nArg == 1 {
		p.queryRetPtr(args[0])
	} else {
		p.queryRetKvPair(args...)
	}
	p.query = nil
	p.ret = nil
}

// For checking query result:
//   - ret &<structVar>
//   - ret &<structSlice>
func (p *Class) queryRetPtr(arg any) {
}

// For checking query result:
//   - ret <colName1>, &<var1>, <colName2>, &<var2>, ...
//   - ret <colName1>, &<varSlice1>, <colName2>, &<varSlice2>, ...
func (p *Class) queryRetKvPair(kvPair ...any) {
	nPair := len(kvPair)
	if nPair < 2 || nPair&1 != 0 {
		log.Panicln("usage: ret <colName1>, &<var1>, <colName2>, &<var2>, ...")
	}
	n := nPair >> 1
	names := make([]string, n)
	rets := make([]any, n)
	for i := 0; i < nPair; i += 2 {
		names[i>>1] = kvPair[i].(string)
		rets[i>>1] = kvPair[i+1]
	}
}

// Query creates a new query.
//   - query <cond>, <arg1>, <arg2>, ...
func (p *Class) Query__0(src ast.Node, cond string, args ...any) {
	p.query = &query{
		cond: cond, args: args,
	}
	p.ret = p.queryRet
}

// Query creates a new query.
func (p *Class) Query__1(cond string, args ...any) {
	p.Query__0(nil, cond, args...)
}

// Limit sets query result rows limit.
func (p *Class) Limit__0(n int, src ...ast.Node) {
	if p.query == nil {
		log.Panicln("please call `limit` after a query statement")
	}
	p.query.limit = n
}

// Limit checks if query result rows is < n or not.
func (p *Class) Limit__1(src ast.Node, n int, cond string, args ...any) {
	ret := p.Count__0(src, cond, args...)
	if ret >= n {
		log.Panicf("limit %s: got %d, expected <%d\n", cond, ret, n)
	}
}

// Limit checks if query result rows is < n or not.
func (p *Class) Limit__2(n int, cond string, args ...any) {
	p.Limit__1(nil, n, cond, args...)
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
		log.Panicln("please call `call` after an api definition")
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
	p.ret = nil
}

// -----------------------------------------------------------------------------
