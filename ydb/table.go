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
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type (
	blob  []byte
	float float64

	date      time.Time
	datetime  time.Time
	timestamp time.Time
)

type basetype interface {
	string | int | bool | blob | date | datetime | timestamp | float
}

func colBaseType(v any) string {
	switch v.(type) {
	case *string:
		return "TEXT"
	case *int:
		return "INT"
	case *bool:
		return "BOOL"
	case *blob:
		return "BLOB"
	case *date:
		return "DATE"
	case *datetime:
		return "DATETIME"
	case *timestamp:
		return "TIMESTAMP"
	case *float:
		return "DOUBLE"
	default:
		panic("unknown column type: " + reflect.TypeOf(v).Elem().String())
	}
}

/*
type baseelem interface {
	byte | int | blob | datetime | timestamp | float
}
*/

func colArrType(v any, n int) string {
	switch v.(type) {
	case byte:
		if n <= 65535 {
			return "TEXT(" + strconv.Itoa(n) + ")"
		}
		if n <= 16777215 {
			return "MEDIUMTEXT"
		}
		return "LONGTEXT"
	case int:
		return "BIGINT(" + strconv.Itoa(n) + ")"
	case blob:
		if n <= 16777215 {
			return "MEDIUMBLOB"
		}
		return "LONGBLOB"
	case datetime:
		return "DATETIME(" + strconv.Itoa(n) + ")"
	case timestamp:
		return "TIMESTAMP(" + strconv.Itoa(n) + ")"
	case float:
		return "FLOAT(" + strconv.Itoa(n) + ")"
	default:
		panic(fmt.Sprintf("unknown column type: [%d]%T", n, v))
	}
}

// -----------------------------------------------------------------------------

type Table struct {
}

func newTable(name, ver string) *Table {
	return &Table{}
}

func (p *Table) Unique(name ...string) {
}

func (p *Table) Index(name ...string) {
}

// From migrates from old table because it's an incompatible change
func (p *Table) From(old string, migrate func()) {
}

func (p *Table) Insert(kvPair ...any) {
}

func (p *Table) Ret(kvPair ...any) {
}

func (p *Table) Query(query string) {
}

func (p *Table) Limit__0(n int) {
}

func (p *Table) Limit__1(n int, query string) {
}

// -----------------------------------------------------------------------------

func Gopt_Table_Gopx_Col__0[T basetype](p *Table, name string) {
	tcol := colBaseType((*T)(nil))
	_ = tcol
}

func Gopt_Table_Gopx_Col__1[Array any](p *Table, name string) {
	varr := (*Array)(nil)
	tarr := reflect.TypeOf(varr).Elem()
	if tarr.Kind() != reflect.Array {
		panic(fmt.Sprintf("unknown column type: %T", varr))
	}
	n := tarr.Len()
	elem := tarr.Elem()
	v := reflect.Zero(reflect.PointerTo(elem)).Interface()
	tcol := colArrType(v, n)
	_ = tcol
}

// -----------------------------------------------------------------------------
