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
	"time"
	"unsafe"
)

type dbType = reflect.Type

type Table struct {
	name   string
	ver    string
	schema dbType
	cols   []*column
	uniqs  [][]string
	idxs   [][]string
}

type column struct {
	typ  string // type in DB
	name string // column name
	fld  field
}

type field struct {
	typ    dbType  // field type
	offset uintptr // offset within struct, in bytes
}

func newTable(name, ver string, schema dbType) *Table {
	n := schema.NumField()
	cols := make([]*column, 0, n)
	p := &Table{name: name, ver: ver, schema: schema, cols: cols}
	p.defineCols(n, schema, 0)
	return p
}

func getVals(vals []any, v reflect.Value, cols []field) []any {
	this := uintptr(v.Addr().UnsafePointer())
	for _, col := range cols {
		val := reflect.NewAt(col.typ, unsafe.Pointer(this+col.offset)).Interface()
		vals = append(vals, val)
	}
	return vals
}

func getCols(names []string, cols []field, n int, t dbType, base uintptr) ([]string, []field) {
	for i := 0; i < n; i++ {
		fld := t.Field(i)
		if fld.Anonymous {
			fldType := fld.Type
			names, cols = getCols(names, cols, fldType.NumField(), fldType, base+fld.Offset)
			continue
		}
		if fld.IsExported() {
			name := ""
			if tag := string(fld.Tag); tag != "" {
				if c := tag[0]; c >= 'a' && c <= 'z' { // suppose a column name is lower case
					if pos := strings.IndexByte(tag, ' '); pos > 0 {
						tag = tag[:pos]
					}
					name = tag
				}
			}
			if name == "" {
				name = dbName(fld.Name)
			}
			names = append(names, name)
			cols = append(cols, field{fld.Type, base + fld.Offset})
		}
	}
	return names, cols
}

func (p *Table) defineCols(n int, t dbType, base uintptr) {
	for i := 0; i < n; i++ {
		fld := t.Field(i)
		if fld.Anonymous {
			fldType := fld.Type
			p.defineCols(fldType.NumField(), fldType, base+fld.Offset)
			continue
		}
		if fld.IsExported() {
			col := &column{fld: field{fld.Type, base + fld.Offset}}
			if tag := string(fld.Tag); tag != "" {
				if parts := strings.Fields(tag); len(parts) > 0 {
					if c := parts[0][0]; c >= 'a' && c <= 'z' { // suppose a column name is lower case
						col.name = parts[0]
						parts = parts[1:]
					} else {
						col.name = dbName(fld.Name)
					}
					for _, part := range parts {
						cmd, params := part, "" // cmd(params)
						if pos := strings.IndexByte(part, '('); pos > 0 && part[len(part)-1] == ')' {
							cmd, params = part[:pos], part[pos+1:len(part)-1]
						}
						switch cmd {
						case `UNIQUE`:
							p.uniqs = append(p.uniqs, makeIndex(col.name, params))
						case `INDEX`:
							p.idxs = append(p.idxs, makeIndex(col.name, params))
						default:
							if col.typ != "" {
								log.Panicf("invalid tag `%s`: multiple column types?\n", tag)
							}
							col.typ = part
						}
					}
				}
			}
			if col.name == "" {
				col.name = dbName(fld.Name)
			}
			if col.typ == "" {
				col.typ = columnType(fld.Type)
			}
			p.cols = append(p.cols, col)
		}
	}
}

var (
	tyString  = reflect.TypeOf("")
	tyInt     = reflect.TypeOf(0)
	tyBool    = reflect.TypeOf(false)
	tyBlob    = reflect.TypeOf([]byte(nil))
	tyTime    = reflect.TypeOf(time.Time{})
	tyFloat64 = reflect.TypeOf(float64(0))
	tyFloat32 = reflect.TypeOf(float32(0))
)

func columnType(fldType dbType) string {
	switch fldType {
	case tyString:
		return "TEXT"
	case tyInt:
		return "INT"
	case tyBool:
		return "BOOL"
	case tyBlob:
		return "BLOB"
	case tyTime:
		return "DATETIME"
	case tyFloat64:
		return "DOUBLE"
	case tyFloat32:
		return "FLOAT"
	}
	panic("unknown column type: " + fldType.String())
}

func makeIndex(name string, params string) []string {
	if params == "" {
		return []string{name}
	}
	pos := strings.IndexByte(params, ',')
	if pos < 0 {
		return []string{name, params}
	}
	ret := make([]string, 1, 4)
	ret[0] = name
	for {
		ret = append(ret, params[:pos])
		params = params[pos+1:]
		pos = strings.IndexByte(params, ',')
		if pos < 0 {
			break
		}
	}
	return append(ret, params)
}

// -----------------------------------------------------------------------------

func (p *Table) create(ctx context.Context, sql *Sql) {
	n := len(p.cols)
	if n == 0 {
		log.Panicln("empty table:", p.name, p.ver)
	}

	query := make([]byte, 0, 64)
	query = append(query, "CREATE TABLE "...)
	query = append(query, p.name...)
	query = append(query, ' ', '(')
	for _, c := range p.cols {
		query = append(query, c.name...)
		query = append(query, ' ')
		query = append(query, c.typ...)
		query = append(query, ',')
	}
	query[len(query)-1] = ')'

	db := sql.db
	_, err := db.ExecContext(ctx, string(query))
	if err != nil {
		log.Panicf("create table (%s): %v\n", p.name, err)
	}

	for _, uniq := range p.uniqs {
		name := indexName(uniq, "uniq_")
		_, err = execWithStrArgs(db, ctx, "CREATE UNIQUE INDEX ? ON ? (", ")", name, p.name, uniq)
		if err != nil {
			log.Panicln("create unique index:", err)
		}
	}
	for _, idx := range p.idxs {
		name := indexName(idx, "idx_")
		_, err = execWithStrArgs(db, ctx, "CREATE INDEX ? ON ? (", ")", name, p.name, idx)
		if err != nil {
			log.Panicln("create index:", err)
		}
	}
}

func indexName(name []string, prefix string) string {
	return prefix + strings.Join(name, "_")
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

// -----------------------------------------------------------------------------
