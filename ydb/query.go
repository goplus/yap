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
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/goplus/gop/ast"
	"github.com/goplus/yap/reflectutil"
	"github.com/qiniu/x/ctype"
)

// -----------------------------------------------------------------------------

// Query creates a new query.
//   - query <cond>, <arg1>, <arg2>, ...
func (p *Class) Query__0(src ast.Expr, cond string, args ...any) {
	p.query = &query{
		cond: cond, args: args,
	}
	p.ret = p.queryRet
}

// Query creates a new query.
//   - query <cond>, <arg1>, <arg2>, ...
func (p *Class) Query__1(cond string, args ...any) {
	p.Query__0(nil, cond, args...)
}

// -----------------------------------------------------------------------------

type query struct {
	cond  string // where
	args  []any  // one of query argument <argN> can be a slice
	limit int    // 0 means no limit
}

func (q *query) makeSelectExpr(tbl string, exprs []string) string {
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
	return string(query)
}

// For checking query result:
//   - ret <expr1>, &<var1>, <expr2>, &<var2>, ...
//   - ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...
//   - ret &<structVar>
//   - ret &<structSlice>
func (p *Class) queryRet(args ...any) (err error) {
	nArg := len(args)
	if nArg == 1 {
		err = p.queryRetPtr(args[0])
	} else {
		err = p.queryRetKvPair(args...)
	}
	p.query = nil
	p.ret = nil
	return
}

// For checking query result:
//   - ret &<structVar>
//   - ret &<structOrPtrSlice>
func (p *Class) queryRetPtr(ret any) error {
	vRet := reflect.ValueOf(ret)
	if vRet.Kind() != reflect.Pointer {
		log.Panicln("usage: ret &<structVar>")
	}

	switch vRet = vRet.Elem(); vRet.Kind() {
	case reflect.Slice:
		return p.queryStrucRows(vRet)
	default:
		return p.queryStrucRow(vRet)
	}
}

// For checking query result:
//   - ret &<structVar>
func (p *Class) queryStrucRow(vRet reflect.Value) error {
	if vRet.Kind() != reflect.Struct {
		log.Panicln("usage: ret &<structVar>")
	}

	n := vRet.NumField()
	names, cols := getCols(make([]string, 0, n), make([]field, 0, n), n, vRet.Type(), 0)
	rets := getVals(make([]any, 0, len(cols)), vRet, cols, false)

	q := p.query
	query := q.makeSelectExpr(p.tbl, names)
	return p.sqlQueryVals(context.TODO(), query, q.args, rets)
}

func (p *Class) queryStrucOne(
	ctx context.Context, query string, args []any,
	vSlice reflect.Value, elem dbType, cols []field, hasPtr bool) error {
	vRet := reflect.New(elem).Elem()
	rets := getVals(make([]any, 0, len(cols)), vRet, cols, false)
	err := p.sqlQueryVals(ctx, query, args, rets)
	if err != nil {
		return err
	}
	if hasPtr {
		vRet = vRet.Addr()
	}
	vSlice.Set(reflect.Append(vSlice, vRet))
	return nil
}

func (p *Class) queryStrucMulti(
	ctx context.Context, query string, args []any, iArgSlice int,
	vSlice reflect.Value, elem dbType, cols []field, hasPtr bool) error {
	argSlice := args[iArgSlice]
	defer func() {
		args[iArgSlice] = argSlice
	}()
	vArgSlice := reflect.ValueOf(argSlice)
	for i, n := 0, vArgSlice.Len(); i < n; i++ {
		arg := vArgSlice.Index(i).Interface()
		args[iArgSlice] = arg
		if err := p.queryStrucOne(ctx, query, args, vSlice, elem, cols, hasPtr); err != nil {
			return err
		}
	}
	return nil
}

// For checking query result:
//   - ret &<structOrPtrSlice>
func (p *Class) queryStrucRows(vSlice reflect.Value) error {
	hasPtr := false
	elem := vSlice.Type().Elem()
	kind := elem.Kind()
	if kind == reflect.Pointer {
		elem, hasPtr = elem.Elem(), true
		kind = elem.Kind()
	}
	if kind != reflect.Struct {
		log.Panicln("usage: ret &<structOrPtrSlice>")
	}

	n := elem.NumField()
	names, cols := getCols(make([]string, 0, n), make([]field, 0, n), n, elem, 0)

	q := p.query
	query := q.makeSelectExpr(p.tbl, names)

	args := q.args
	iArgSlice := checkArgSlice(args)
	if iArgSlice >= 0 {
		return p.queryStrucMulti(context.TODO(), query, args, iArgSlice, vSlice, elem, cols, hasPtr)
	}
	return p.queryStrucOne(context.TODO(), query, args, vSlice, elem, cols, hasPtr)
}

// sqlQueryVals NOTE:
//   - one of args maybe is a slice
func (p *Class) sqlQueryVals(ctx context.Context, query string, args, rets []any) error {
	iArgSlice := checkArgSlice(args)
	if iArgSlice >= 0 {
		log.Panicln("one of `query` arguments is a slice, but `ret` arguments are not")
	}

	if debugExec {
		log.Println("==>", query, args)
	}
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		p.handleErr("query:", err)
		return err
	}
	defer rows.Close()

	return p.sqlRetRow(rows, rets)
}

func isSlice(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

func checkArgSlice(args []any) int {
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
	return iArgSlice
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

func (p *Class) sqlRetRow(rows *sql.Rows, rets []any) error {
	if !rows.Next() {
		err := rows.Err()
		if err == nil {
			err = ErrNoRows
		}
		p.handleErr("ret:", err)
		return err
	}
	err := rows.Scan(rets...)
	if err != nil {
		p.handleErr("ret:", err)
	}
	return err
}

func (p *Class) sqlRetRows(rows *sql.Rows, vRets []reflect.Value, oneRet []any, needInit bool) error {
	for rows.Next() {
		if needInit {
			for _, ret := range oneRet {
				reflectutil.SetZero(reflect.ValueOf(ret).Elem())
			}
		} else {
			needInit = true
		}
		err := rows.Scan(oneRet...)
		if err != nil {
			p.handleErr("ret:", err)
			return err
		}
		for i, vRet := range vRets {
			v := reflect.ValueOf(oneRet[i])
			vRet.Set(reflect.Append(vRet, v.Elem()))
		}
	}
	err := rows.Err()
	if err != nil {
		p.handleErr("ret:", err)
	}
	return err
}

// sqlQueryRows NOTE:
//   - one of args maybe is a slice
func (p *Class) sqlQueryRows(ctx context.Context, query string, args, rets []any) error {
	iArgSlice := checkArgSlice(args)
	if iArgSlice >= 0 {
		return p.sqlMultiQuery(ctx, query, iArgSlice, args, rets)
	}

	if debugExec {
		log.Println("==>", query, args)
	}
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		p.handleErr("query:", err)
		return err
	}
	defer rows.Close()

	vRets, oneRet := makeSliceRets(rets)
	return p.sqlRetRows(rows, vRets, oneRet, false)
}

func makeSliceRets(rets []any) (vRets []reflect.Value, oneRet []any) {
	vRets = make([]reflect.Value, len(rets))
	oneRet = make([]any, len(rets))
	for i, ret := range rets {
		slice := reflect.ValueOf(ret).Elem()
		vRets[i] = slice

		elem := slice.Type().Elem()
		oneRet[i] = reflect.New(elem).Interface()
	}
	return
}

func (p *Class) sqlQueryOne(ctx context.Context, query string, args, oneRet []any, vRets []reflect.Value) error {
	if debugExec {
		log.Println("==>", query, args)
	}
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		p.handleErr("query:", err)
		return err
	}
	defer rows.Close()

	return p.sqlRetRows(rows, vRets, oneRet, true)
}

func (p *Class) sqlMultiQuery(ctx context.Context, query string, iArgSlice int, args, rets []any) error {
	argSlice := args[iArgSlice]
	defer func() {
		args[iArgSlice] = argSlice
	}()
	vRets, oneRet := makeSliceRets(rets)
	vArgSlice := reflect.ValueOf(argSlice)
	for i, n := 0, vArgSlice.Len(); i < n; i++ {
		arg := vArgSlice.Index(i).Interface()
		args[iArgSlice] = arg
		if err := p.sqlQueryOne(ctx, query, args, oneRet, vRets); err != nil {
			return err
		}
	}
	return nil
}

// For checking query result:
//   - ret <expr1>, &<var1>, <expr2>, &<var2>, ...
//   - ret <expr1>, &<varSlice1>, <expr2>, &<varSlice2>, ...
func (p *Class) queryRetKvPair(kvPair ...any) error {
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

	query := q.makeSelectExpr(tbl, exprs)
	if kind == valFlagNormal {
		return p.sqlQueryVals(context.TODO(), query, q.args, rets)
	}
	return p.sqlQueryRows(context.TODO(), query, q.args, rets)
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

// -----------------------------------------------------------------------------

// Limit sets query result rows limit.
func (p *Class) Limit__0(n int, src ...ast.Expr) {
	if p.query == nil {
		log.Panicln("please call `limit` after a query statement")
	}
	p.query.limit = n
}

// Limit checks if query result rows is < n or not.
func (p *Class) Limit__1(src ast.Expr, n int, cond string, args ...any) error {
	ret, err := p.Count__0(src, cond, args...)
	if err != nil {
		return err
	}
	if ret >= n {
		if p.onErr == nil {
			log.Panicf("limit %s: got %d, expected <%d\n", cond, ret, n)
		}
		err = ErrOutOfLimit
		p.onErr(err)
	}
	return err
}

// Limit checks if query result rows is < n or not.
func (p *Class) Limit__2(n int, cond string, args ...any) error {
	return p.Limit__1(nil, n, cond, args...)
}

// -----------------------------------------------------------------------------

// Count returns rows of a query result.
func (p *Class) Count__0(src ast.Expr, cond string, args ...any) (n int, err error) {
	if p.tbl == "" {
		log.Panicln("please call `use <tableName>` to specified a table name")
	}
	row := p.db.QueryRowContext(context.TODO(), "SELECT COUNT(*) FROM "+p.tbl+" WHERE "+cond, args...)
	if err = row.Scan(&n); err != nil {
		p.handleErr("query:", err)
	}
	return
}

// Count returns rows of a query result.
func (p *Class) Count__1(cond string, args ...any) (n int, err error) {
	return p.Count__0(nil, cond, args...)
}

// -----------------------------------------------------------------------------
