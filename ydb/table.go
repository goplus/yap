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

	"github.com/goplus/yap/stringutil"
)

type dbType = reflect.Type

type dbIndex struct {
	index  []*column
	col    *column
	params string
}

func (p *dbIndex) get(tbl *Table) []*column {
	if p.index == nil {
		p.index = tbl.makeIndex(p.col, p.params)
	}
	return p.index
}

type Table struct {
	name   string
	ver    string
	schema dbType
	cols   []*column
	uniqs  []*dbIndex
	idxs   []*dbIndex
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

func getVals(vals []any, v reflect.Value, cols []field, elem bool) []any {
	this := uintptr(v.Addr().UnsafePointer())
	for _, col := range cols {
		v := reflect.NewAt(col.typ, unsafe.Pointer(this+col.offset))
		if elem {
			v = v.Elem()
		}
		val := v.Interface()
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
							p.uniqs = append(p.uniqs, &dbIndex{nil, col, params})
						case `INDEX`:
							p.idxs = append(p.idxs, &dbIndex{nil, col, params})
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

func (p *Table) makeIndex(col *column, params string) []*column {
	if params == "" {
		return []*column{col}
	}
	pos := strings.IndexByte(params, ',')
	if pos < 0 {
		return []*column{col, p.getCol(params)}
	}
	ret := make([]*column, 1, 4)
	ret[0] = col
	for {
		ret = append(ret, p.getCol(params[:pos]))
		params = params[pos+1:]
		pos = strings.IndexByte(params, ',')
		if pos < 0 {
			break
		}
	}
	return append(ret, p.getCol(params))
}

func (p *Table) getCol(name string) *column {
	for _, col := range p.cols {
		if col.name == name {
			return col
		}
	}
	log.Panicf("table `%s` doesn't have column `%s`\n", p.name, name)
	return nil
}

// -----------------------------------------------------------------------------

func (p *Table) create(ctx context.Context, sql *Sql) {
	n := len(p.cols)
	if n == 0 {
		log.Panicln("empty table:", p.name, p.ver)
	}

	db := sql.db
	query := make([]byte, 0, 64)
	if sql.autodrop {
		query = append(query, "DROP TABLE "...)
		query = append(query, p.name...)
		db.ExecContext(ctx, string(query))
		query = query[:0]
	}

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

	q := string(query)
	_, err := db.ExecContext(ctx, q)
	if err != nil {
		log.Panicf("%s\ncreate table (%s): %v\n", q, p.name, err)
	}

	for _, uniq := range p.uniqs {
		cols := uniq.get(p)
		name := indexName(cols, "uniq_", p.name)
		createIndex(sql, db, ctx, "CREATE UNIQUE INDEX ", name, p.name, cols)
	}
	for _, idx := range p.idxs {
		cols := idx.get(p)
		name := indexName(cols, "idx_", p.name)
		createIndex(sql, db, ctx, "CREATE INDEX ", name, p.name, cols)
	}
}

// prefix_tbl_name1_name2_...
func indexName(cols []*column, prefix, tbl string) string {
	n := len(prefix) + len(tbl)
	for _, col := range cols {
		n += 1 + len(col.name)
	}
	b := make([]byte, 0, n)
	b = append(b, prefix...)
	b = append(b, tbl...)
	for _, col := range cols {
		b = append(b, '_')
		b = append(b, col.name...)
	}
	return stringutil.String(b)
}

func createIndex(sql *Sql, db *sql.DB, ctx context.Context, cmd string, name, tbl string, cols []*column) {
	parts := make([]string, 0, 5+2*len(cols))
	parts = append(parts, cmd, name, " ON ", tbl, "(")
	for _, col := range cols {
		parts = append(parts, col.name, ",")
	}
	parts[len(parts)-1] = ")"
	query := stringutil.Concat(parts...)
	if _, err := db.ExecContext(ctx, query); err != nil {
		log.Panicf("%s\ncreate index `%s`: %v\n", query, name, err)
	}
}

// -----------------------------------------------------------------------------
