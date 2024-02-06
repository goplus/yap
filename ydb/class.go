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
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/goplus/gop/ast"
	"github.com/qiniu/x/ctype"
)

var (
	ErrNoRows     = sql.ErrNoRows
	ErrDuplicated = errors.New("duplicated")
)

// -----------------------------------------------------------------------------

type Class struct {
	name string
	tbl  string
	tobj *Table
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
	tblobj, ok := p.sql.tables[table]
	if !ok {
		log.Panicln("table not found:", table)
	}
	p.tbl = table
	p.tobj = tblobj
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
//   - insert <structOrPtrSlice>
func (p *Class) Insert__0(src ast.Node, args ...any) {
	if p.tbl == "" {
		log.Panicln("please call `use <tableName>` to specified current table")
	}
	nArg := len(args)
	if nArg == 1 {
		p.insertStruc(args[0])
	} else {
		p.insertKvPair(args...)
	}
}

// Insert inserts a new row.
//   - insert <structValOrPtr>
//   - insert <structOrPtrSlice>
func (p *Class) insertStruc(arg any) {
	vArg := reflect.ValueOf(arg)
	switch vArg.Kind() {
	case reflect.Slice:
		p.insertStrucRows(vArg)
	case reflect.Pointer:
		vArg = vArg.Elem()
		fallthrough
	default:
		p.insertStrucRow(vArg)
	}
}

func (p *Class) insertStrucRows(vSlice reflect.Value) {
	rows := vSlice.Len()
	if rows == 0 {
		return
	}
	hasPtr := false
	elem := vSlice.Type().Elem()
	kind := elem.Kind()
	if kind == reflect.Pointer {
		elem, hasPtr = elem.Elem(), true
		kind = elem.Kind()
	}
	if kind != reflect.Struct {
		log.Panicln("usage: insert <structOrPtrSlice>")
	}
	n := elem.NumField()
	names, cols := getCols(make([]string, 0, n), make([]field, 0, n), n, elem, 0)
	vals := make([]any, 0, len(names)*rows)
	for row := 0; row < rows; row++ {
		vElem := vSlice.Index(row)
		if hasPtr {
			vElem = vElem.Elem()
		}
		vals = getVals(vals, vElem, cols)
	}
	p.insertRowsVals(names, vals, rows)
}

func (p *Class) insertStrucRow(vArg reflect.Value) {
	if vArg.Kind() != reflect.Struct {
		log.Panicln("usage: insert <structValOrPtr>")
	}
	n := vArg.NumField()
	names, cols := getCols(make([]string, 0, n), make([]field, 0, n), n, vArg.Type(), 0)
	vals := getVals(make([]any, 0, len(cols)), vArg, cols)
	p.insertRow(names, vals)
}

const (
	valFlagNormal  = 1
	valFlagSlice   = 2
	valFlagInvalid = valFlagNormal | valFlagSlice
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
		p.insertSliceRows(names, vals, rows)
	default:
		log.Panicln("can't insert mix slice and normal value")
	}
}

// NOTE: len(args) == len(names)
func (p *Class) insertSliceRows(names []string, args []any, rows int) {
	vals := make([]any, 0, len(names)*rows)
	for i := 0; i < rows; i++ {
		for _, arg := range args {
			v := arg.(reflect.Value)
			vals = append(vals, v.Index(i).Interface())
		}
	}
	p.insertRowsVals(names, vals, rows)
}

// NOTE: len(vals) == len(names) * rows
func (p *Class) insertRowsVals(names []string, vals []any, rows int) {
	n := len(names)
	query := insertQuery(p.tbl, names)
	query = append(query, valParams(n, rows)...)

	result, err := p.db.ExecContext(context.TODO(), string(query), vals...)
	insertRet(result, err)
}

func (p *Class) insertRow(names []string, vals []any) {
	if len(names) == 0 {
		log.Panicln("insert: nothing to insert")
	}
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

func valParams(n, rows int) string {
	valparam := valParam(n)
	valparams := strings.Repeat(valparam+",", rows)
	valparams = valparams[:len(valparams)-1]
	return valparams
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
//   - ret <expr1>, &<var1>, <expr2>, &<var2>, ...
//   - ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...
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

func isSlice(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

func retKind(ret any) int {
	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Pointer {
		log.Panicln("usage: ret <expr1>, &<var1>, <expr2>, &<var2>, ...")
	}
	if v.Elem().Kind() == reflect.Slice {
		return valFlagSlice
	}
	return valFlagNormal
}

func sqlRetRow(rows *sql.Rows, rets []any) {
	if !rows.Next() {
		err := rows.Err()
		if err == nil {
			err = ErrNoRows
		}
		log.Panicln("ret:", err)
	}
	err := rows.Scan(rets...)
	if err != nil {
		log.Panicln("ret:", err)
	}
}

func sqlRetRows(rows *sql.Rows, vRets []reflect.Value, oneRet []any, needInit bool) {
	for rows.Next() {
		if needInit {
			for _, ret := range oneRet {
				reflect.ValueOf(ret).Elem().SetZero()
			}
		} else {
			needInit = true
		}
		err := rows.Scan(oneRet...)
		if err != nil {
			log.Panicln("ret:", err)
		}
		for i, vRet := range vRets {
			v := reflect.ValueOf(oneRet[i])
			vRet.Set(reflect.Append(vRet, v.Elem()))
		}
	}
	if err := rows.Err(); err != nil {
		log.Panicln("ret:", err)
	}
}

// sqlQuery NOTE:
//   - one of args maybe is a slice
func sqlQuery(db *sql.DB, ctx context.Context, query string, args, rets []any, retSlice bool) {
	iArgSlice := -1
	for i, arg := range args {
		if isSlice(arg) {
			if iArgSlice >= 0 {
				log.Panicf(
					"query: multiple arguments (%dth, %dth) are slices (only one can be)\n",
					iArgSlice+1, i+1,
				)
			}
			iArgSlice = i
		}
	}
	if iArgSlice >= 0 {
		if !retSlice {
			log.Panicln("one of `query` arguments is a slice, but `ret` arguments are not")
		}
		sqlMultiQuery(db, ctx, query, iArgSlice, args, rets)
		return
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Panicln("query:", err)
	}
	defer rows.Close()

	if retSlice {
		vRets, oneRet := makeSliceRets(rets)
		sqlRetRows(rows, vRets, oneRet, false)
		return
	}
	sqlRetRow(rows, rets)
}

func makeSliceRets(rets []any) (vRets []reflect.Value, oneRet []any) {
	vRets = make([]reflect.Value, len(rets))
	oneRet = make([]any, len(rets))
	for i, ret := range rets {
		slice := reflect.ValueOf(ret).Elem()
		slice.SetZero()
		vRets[i] = slice

		elem := slice.Type().Elem()
		oneRet[i] = reflect.New(elem).Interface()
	}
	return
}

func sqlMultiQuery(db *sql.DB, ctx context.Context, query string, iArgSlice int, args, rets []any) {
}

// For checking query result:
//   - ret <expr1>, &<var1>, <expr2>, &<var2>, ...
//   - ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...
func (p *Class) queryRetKvPair(kvPair ...any) {
	nPair := len(kvPair)
	if nPair < 2 || nPair&1 != 0 {
		log.Panicln("usage: ret <expr1>, &<var1>, <expr2>, &<var2>, ...")
	}

	q := p.query
	tbl := p.exprTblname(q.cond)

	n := nPair >> 1
	exprs := make([]string, n)
	rets := make([]any, n)
	kind := 0
	for i := 0; i < nPair; i += 2 {
		expr := kvPair[i].(string)
		if etbl := p.exprTblname(expr); etbl != tbl {
			log.Panicf(
				"query currently doesn't support multiple tables: `query` use `%s` but `ret` use `%s`\n",
				tbl, etbl,
			)
		}
		ret := kvPair[i+1]
		kind |= retKind(ret)
		exprs[i>>1] = expr
		rets[i>>1] = ret
	}
	if kind == valFlagInvalid {
		log.Panicln(`all ret arguments should be address of slices or address of normal variable:
	ret <expr1>, &<var1>, <expr2>, &<var2>, ...
	ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...`)
	}

	query := make([]byte, 0, 128)
	query = append(query, "SELECT "...)
	query = append(query, strings.Join(exprs, ",")...)
	query = append(query, " FROM "...)
	query = append(query, tbl...)
	query = append(query, " WHERE "...)
	query = append(query, q.cond...)
	if q.limit > 0 {
		query = append(query, " LIMIT "...)
		query = append(query, strconv.Itoa(q.limit)...)
	}
	sqlQuery(p.db, context.TODO(), string(query), q.args, rets, kind == valFlagSlice)
}

func (p *Class) exprTblname(cond string) string {
	tbls := exprTblnames(cond)
	tbl := ""
	switch len(tbls) {
	case 0:
	case 1:
		tbl = tbls[0]
	default:
		log.Panicln("query currently doesn't support multiple tables")
	}
	if tbl == "" {
		tbl = p.tbl
	}
	return tbl
}

func exprTblnames(expr string) (tbls []string) {
	for expr != "" {
		pos := ctype.ScanCSymbol(expr)
		if pos != 0 {
			name := ""
			if pos > 0 {
				switch expr[pos] {
				case '.':
					name = expr[:pos]
					expr = ctype.SkipCSymbol(expr[pos+1:])
				case '(': // function call, eg. SUM(...)
					expr = expr[pos+1:]
					continue
				default:
					expr = expr[pos:]
				}
			} else {
				expr = ""
			}
			switch name {
			case "AND", "OR":
			default:
				tbls = addTblname(tbls, name)
			}
			continue
		}
		pos = ctype.ScanTypeEx(ctype.FLOAT_FIRST_CHAT, ctype.CSYMBOL_NEXT_CHAR, expr)
		if pos == 0 {
			c, size := utf8.DecodeRuneInString(expr)
			switch c {
			case '\'':
				expr = skipStringConst(expr[1:], '\'')
			default:
				expr = expr[size:]
			}
		} else if pos < 0 {
			break
		} else {
			expr = expr[pos:]
		}
	}
	return
}

func skipStringConst(next string, quot rune) string {
	skip := false
	for i, c := range next {
		if skip {
			skip = false
		} else if c == '\\' {
			skip = true
		} else if c == quot {
			return next[i+1:]
		}
	}
	return ""
}

func addTblname(tbls []string, tbl string) []string {
	for _, v := range tbls {
		if v == tbl {
			return tbls
		}
	}
	return append(tbls, tbl)
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
