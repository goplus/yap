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
	"strings"

	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

// Insert inserts new rows.
//   - insert <colName1>, <val1>, <colName2>, <val2>, ...
//   - insert <colName1>, <valSlice1>, <colName2>, <valSlice2>, ...
//   - insert <structValOrPtr>
//   - insert <structOrPtrSlice>
func (p *Class) Insert__0(src ast.Expr, args ...any) (sql.Result, error) {
	if p.tbl == "" {
		log.Panicln("please call `use <tableName>` to specified current table")
	}
	nArg := len(args)
	if nArg == 1 {
		return p.insertStruc(args[0])
	}
	return p.insertKvPair(args...)
}

// Insert inserts new rows.
//   - insert <colName1>, <val1>, <colName2>, <val2>, ...
//   - insert <colName1>, <valSlice1>, <colName2>, <valSlice2>, ...
//   - insert <structValOrPtr>
//   - insert <structOrPtrSlice>
func (p *Class) Insert__1(kvPair ...any) (sql.Result, error) {
	return p.Insert__0(nil, kvPair...)
}

// Insert inserts a new row.
//   - insert <structValOrPtr>
//   - insert <structOrPtrSlice>
func (p *Class) insertStruc(arg any) (sql.Result, error) {
	vArg := reflect.ValueOf(arg)
	switch vArg.Kind() {
	case reflect.Slice:
		return p.insertStrucRows(vArg)
	case reflect.Pointer:
		vArg = vArg.Elem()
		fallthrough
	default:
		return p.insertStrucRow(vArg)
	}
}

func (p *Class) insertStrucRows(vSlice reflect.Value) (sql.Result, error) {
	rows := vSlice.Len()
	if rows == 0 {
		return nil, nil
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
		vals = getVals(vals, vElem, cols, true)
	}
	return p.insertRowsVals(p.tbl, names, vals, rows)
}

func (p *Class) insertStrucRow(vArg reflect.Value) (sql.Result, error) {
	if vArg.Kind() != reflect.Struct {
		log.Panicln("usage: insert <structValOrPtr>")
	}
	n := vArg.NumField()
	names, cols := getCols(make([]string, 0, n), make([]field, 0, n), n, vArg.Type(), 0)
	vals := getVals(make([]any, 0, len(cols)), vArg, cols, true)
	return p.insertRow(p.tbl, names, vals)
}

const (
	valFlagNormal  = 1
	valFlagSlice   = 2
	valFlagInvalid = valFlagNormal | valFlagSlice
)

// Insert inserts a new row.
//   - insert <colName1>, <val1>, <colName2>, <val2>, ...
//   - insert <colName1>, <valSlice1>, <colName2>, <valSlice2>, ...
func (p *Class) insertKvPair(kvPair ...any) (sql.Result, error) {
	nPair := len(kvPair)
	if nPair < 2 || nPair&1 != 0 {
		log.Panicln("usage: insert <colName1>, <val1>, <colName2>, <val2>, ...")
	}
	n := nPair >> 1
	names := make([]string, n)
	vals := make([]any, n)
	rows := -1      // slice length
	iArgSlice := -1 // -1: no slice, or index of first slice arg
	kind := 0
	for iPair := 0; iPair < nPair; iPair += 2 {
		i := iPair >> 1
		names[i] = kvPair[iPair].(string)
		val := kvPair[iPair+1]
		switch v := reflect.ValueOf(val); v.Kind() {
		case reflect.Slice:
			vlen := v.Len()
			if iArgSlice == -1 {
				iArgSlice = i
				rows = vlen
			} else if rows != vlen {
				log.Panicf("insert: unexpected slice length. got %d, expected %d\n", vlen, rows)
			} else {
				kind |= valFlagSlice
			}
			vals[i] = v
		default:
			kind |= valFlagNormal
			vals[i] = val
		}
	}
	if kind == valFlagInvalid {
		log.Panicln("insert: can't mix multiple slice arguments and normal value")
	}
	tbl := p.tblFromNames(names)
	if kind == valFlagSlice {
		return p.insertSlice(tbl, names, vals, rows)
	}
	if iArgSlice == -1 {
		return p.insertRow(tbl, names, vals)
	}
	return p.insertMulti(tbl, names, iArgSlice, vals)
}

// NOTE: len(args) == len(names)
func (p *Class) insertMulti(tbl string, names []string, iArgSlice int, args []any) (sql.Result, error) {
	argSlice := args[iArgSlice]
	defer func() {
		args[iArgSlice] = argSlice
	}()
	vArgSlice := argSlice.(reflect.Value)
	rows := vArgSlice.Len()
	vals := make([]any, 0, len(names)*rows)
	for i := 0; i < rows; i++ {
		args[iArgSlice] = vArgSlice.Index(i).Interface()
		vals = append(vals, args...)
	}
	return p.insertRowsVals(tbl, names, vals, rows)
}

// NOTE: len(args) == len(names)
func (p *Class) insertSlice(tbl string, names []string, args []any, rows int) (sql.Result, error) {
	vals := make([]any, 0, len(names)*rows)
	for i := 0; i < rows; i++ {
		for _, arg := range args {
			v := arg.(reflect.Value)
			vals = append(vals, v.Index(i).Interface())
		}
	}
	return p.insertRowsVals(tbl, names, vals, rows)
}

// NOTE: len(vals) == len(names) * rows
func (p *Class) insertRowsVals(tbl string, names []string, vals []any, rows int) (sql.Result, error) {
	n := len(names)
	query := makeInsertExpr(tbl, names)
	query = append(query, valParams(n, rows)...)

	q := string(query)
	if debugExec {
		log.Println("==>", q, vals)
	}
	result, err := p.db.ExecContext(context.TODO(), q, vals...)
	return p.insertRet(result, err)
}

func (p *Class) insertRow(tbl string, names []string, vals []any) (sql.Result, error) {
	if len(names) == 0 {
		log.Panicln("insert: nothing to insert")
	}
	query := makeInsertExpr(tbl, names)
	query = append(query, valParam(len(vals))...)

	q := string(query)
	if debugExec {
		log.Println("==>", q, vals)
	}
	result, err := p.db.ExecContext(context.TODO(), q, vals...)
	return p.insertRet(result, err)
}

func (p *Class) insertRet(result sql.Result, err error) (sql.Result, error) {
	if err != nil {
		p.handleErr("insert:", err)
	}
	return result, err
}

func makeInsertExpr(tbl string, names []string) []byte {
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

// -----------------------------------------------------------------------------
