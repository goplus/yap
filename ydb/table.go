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
	"strconv"
	"strings"
	"time"
)

var (
	ErrDuplicated = errors.New("duplicated")
)

type (
	String = string
	Int    = int
	Bool   = bool
	Byte   = byte

	Blob  []byte
	Float float64

	Date      time.Time
	DateTime  time.Time
	Timestamp time.Time
)

type basetype interface {
	String | Int | Bool | Blob | Date | DateTime | Timestamp | Float
}

/*
type baseelem interface {
	Byte | Int | Blob | DateTime | Timestamp | Float
}
*/

func colBaseType(v any) string {
	switch v.(type) {
	case *String:
		return "TEXT"
	case *Int:
		return "INT"
	case *Bool:
		return "BOOL"
	case *Blob:
		return "BLOB"
	case *Date:
		return "DATE"
	case *DateTime:
		return "DATETIME"
	case *Timestamp:
		return "TIMESTAMP"
	case *Float:
		return "DOUBLE"
	default:
		panic("unknown column type: " + reflect.TypeOf(v).Elem().String())
	}
}

func colArrType(v any, n int) string {
	switch v.(type) {
	case *Byte:
		if n <= 65535 {
			return "TEXT(" + strconv.Itoa(n) + ")"
		}
		if n <= 16777215 {
			return "MEDIUMTEXT"
		}
		return "LONGTEXT"
	case *Int:
		return "BIGINT(" + strconv.Itoa(n) + ")"
	case *Blob:
		if n <= 16777215 {
			return "MEDIUMBLOB"
		}
		return "LONGBLOB"
	case *DateTime:
		return "DATETIME(" + strconv.Itoa(n) + ")"
	case *Timestamp:
		return "TIMESTAMP(" + strconv.Itoa(n) + ")"
	case *Float:
		return "FLOAT(" + strconv.Itoa(n) + ")"
	default:
		panic(fmt.Sprintf("unknown column type: [%d]%T", n, v))
	}
}

// -----------------------------------------------------------------------------

type Table struct {
	name  string
	ver   string
	cols  []*column
	uniqs [][]string
	idxs  [][]string
}

func newTable(name, ver string) *Table {
	return &Table{name: name, ver: ver}
}

// From migrates from old table because it's an incompatible change
func (p *Table) From(old string, migrate func()) {
}

// -----------------------------------------------------------------------------

func (p *Table) Unique(name ...string) {
	p.uniqs = append(p.uniqs, name)
}

func (p *Table) Index(name ...string) {
	p.idxs = append(p.idxs, name)
}

// -----------------------------------------------------------------------------

type column struct {
	typ  string // type in DB
	name string // column name
	link string // optional
	zero any    // zero is (*basetype)(nil) if n == 0; else zero is (*baseelem)(nil)
	n    int    // array if n != 0
}

func (p *Table) create(ctx context.Context, sql *Sql) {
	n := len(p.cols)
	if n == 0 {
		log.Panicln("empty table:", p.name, p.ver)
	}
	fldQuery := strings.Repeat("? ?,\n", n)
	query := "CREATE TABLE ? (" + fldQuery[:len(fldQuery)-2] + ")"
	args := make([]any, 1, n*2+1)
	args[0] = p.name
	for _, c := range p.cols {
		args = append(args, c.name, c.typ)
	}
	db := sql.db
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Panicf("create table (%s): %v\n", p.name, err)
	}

	for _, uniq := range p.uniqs {
		name := indexName(uniq, true)
		_, err = execWithStrArgs(db, ctx, "CREATE UNIQUE INDEX ? ON ? (", ")", name, p.name, uniq)
		if err != nil {
			log.Panicln("create unique index:", err)
		}
	}
	for _, idx := range p.idxs {
		name := indexName(idx, false)
		_, err = execWithStrArgs(db, ctx, "CREATE INDEX ? ON ? (", ")", name, p.name, idx)
		if err != nil {
			log.Panicln("create index:", err)
		}
	}
}

func indexName(name []string, uniq bool) string {
	panic("todo")
}

func execWithStrArgs(db *sql.DB, ctx context.Context, queryPrefix, querySuffix string, v1, v2 string, args []string) (sql.Result, error) {
	switch len(args) {
	case 1:
		return db.ExecContext(ctx, queryPrefix+"?"+querySuffix, v1, v2, args[0])
	case 2:
		return db.ExecContext(ctx, queryPrefix+"?,?"+querySuffix, v1, v2, args[0], args[1])
	default:
		fldQuery := strings.Repeat("?,", len(args))
		query := queryPrefix + fldQuery[:len(fldQuery)-1] + querySuffix
		vArgs := make([]any, 2+len(args))
		vArgs[0], vArgs[1] = v1, v2
		for i, arg := range args {
			vArgs[2+i] = arg
		}
		return db.ExecContext(ctx, query, vArgs...)
	}
}

func (p *Table) defineCol(c *column) {
	p.cols = append(p.cols, c)
}

func Gopt_Table_Gopx_Col__0[T basetype](tbl interface{ defineCol(c *column) }, name string, link ...string) {
	vcol := (*T)(nil)
	tcol := colBaseType(vcol)
	tbl.defineCol(&column{
		typ:  tcol,
		name: name,
		link: optString(link),
		zero: vcol,
	})
}

func Gopt_Table_Gopx_Col__1[Array any](tbl interface{ defineCol(c *column) }, name string, link ...string) {
	varr := (*Array)(nil)
	tarr := reflect.TypeOf(varr).Elem()
	if tarr.Kind() != reflect.Array {
		panic(fmt.Sprintf("unknown column type: %T", varr))
	}
	n := tarr.Len()
	elem := tarr.Elem()
	velem := reflect.Zero(reflect.PointerTo(elem)).Interface()
	tcol := colArrType(velem, n)
	tbl.defineCol(&column{
		typ:  tcol,
		name: name,
		link: optString(link),
		zero: velem,
		n:    n,
	})
}

func optString(v []string) string {
	if v != nil {
		return v[0]
	}
	return ""
}

// -----------------------------------------------------------------------------
