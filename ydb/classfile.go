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
	"strings"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

type dbTable = Table
type dbClass = Class

type Sql struct {
	driver  string
	tables  map[string]*Table
	classes map[string]*Class
	*dbTable
	*dbClass
}

func (p *Sql) initSql() {
	p.tables = make(map[string]*Table)
	p.classes = make(map[string]*Class)
}

// Engine sets engine name of a sql driver.
func (p *Sql) Engine(name string) {
	p.driver = name
}

// Table creates a new table.
func (p *Sql) Table(nameVer string, creator func()) {
	pos := strings.IndexByte(nameVer, ' ') // user v0.1.0
	if pos < 0 {
		panic("table name should have a version: eg. `user v0.1.0`")
	}
	name, ver := nameVer[:pos], strings.TrimLeft(nameVer[pos+1:], " \t")
	p.dbTable = newTable(name, ver)
	p.tables[name] = p.dbTable
	creator()
	p.dbTable = nil
}

// Class creates a new class.
func (p *Sql) Class(name string, creator func()) {
	p.dbClass = newClass(name)
	p.classes[name] = p.dbClass
	creator()
	p.dbClass = nil
}

// Ret checks a query or call result.
func (p *Sql) Ret(args ...any) {
	if tbl := p.dbTable; tbl != nil {
		tbl.Ret(args...)
	} else if cls := p.dbClass; cls != nil {
		cls.Ret(args...)
	} else {
		panic("pelase use ret after query or call")
	}
}

// -----------------------------------------------------------------------------

type AppGen struct {
}

func (p *AppGen) initApp() {
}

func Gopt_AppGen_Main(app interface{ initApp() }, workers ...interface{ initSql() }) {
	app.initApp()
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	for _, worker := range workers {
		worker.initSql()
		worker.(interface{ Main() }).Main()
	}
}

// -----------------------------------------------------------------------------
