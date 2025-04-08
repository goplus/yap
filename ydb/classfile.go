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
	"flag"
	"log"
	"reflect"
	"strings"
)

const (
	GopPackage = "github.com/qiniu/x/test"
)

var (
	debugExec bool
)

// -----------------------------------------------------------------------------

type Sql struct {
	driver *Engine
	wrap   func(string, error) error

	tables  map[string]*Table
	classes map[string]*Class
	db      *sql.DB

	autodrop bool
}

func (p *Sql) initSql() {
	p.tables = make(map[string]*Table)
	p.classes = make(map[string]*Class)
}

// Engine initializes database by specified engine name.
func (p *Sql) Engine__0(name string) {
	driver, ok := engines[name]
	if !ok {
		log.Panicf("engine `%s` not found: please call ydb.Register first\n", name)
	}
	defaultDataSource := driver.TestSource
	dataSource, ok := defaultDataSource.(string)
	if !ok {
		dataSource = defaultDataSource.(func() string)()
	}
	const (
		autodropParam = "autodrop"
	)
	if strings.HasSuffix(dataSource, autodropParam) {
		dataSource = dataSource[:len(dataSource)-len(autodropParam)-1]
		p.autodrop = true
	}
	db, err := sql.Open(name, dataSource)
	if err != nil {
		log.Panicln("sql.Open:", err)
	}
	p.db = db
	p.driver = driver
	p.wrap = driver.WrapErr
}

// Engine returns engine name of the database.
func (p *Sql) Engine__1() string {
	return p.driver.Name
}

func (p *Sql) defineTable(nameVer string, zeroSchema any) {
	var name, ver string
	pos := strings.IndexByte(nameVer, ' ') // user v0.1.0
	if pos < 0 {
		ver = nameVer
	} else {
		name, ver = nameVer[:pos], strings.TrimLeft(nameVer[pos+1:], " \t")
	}
	schema := reflect.TypeOf(zeroSchema).Elem()
	if name == "" {
		name = dbName(schema.Name())
	}
	if _, ok := p.tables[name]; ok {
		log.Panicf("table `%s` exists\n", name)
	}
	tbl := newTable(name, ver, schema)
	p.tables[name] = tbl
	tbl.create(context.TODO(), p)
}

func dbName(fldName string) string {
	c := fldName[0]
	if c >= 'A' && c <= 'Z' {
		c += ('a' - 'A')
	}
	return string(c) + fldName[1:]
}

// Table creates a new table by specified Schema.
func Gopt_Sql_Gopx_Table[Schema any](sql interface{ defineTable(string, any) }, nameVer string) {
	sql.defineTable(nameVer, (*Schema)(nil))
}

// -----------------------------------------------------------------------------

type Engine struct {
	Name       string
	TestSource any // can be a `string` or a `func() string` object.
	WrapErr    func(prompt string, err error) error
}

var (
	engines = make(map[string]*Engine) // engineName => engine
)

// Register registers an engine.
func Register(e *Engine) {
	engines[e.Name] = e
}

// -----------------------------------------------------------------------------

type AppGen struct {
}

func (p *AppGen) initApp() {
}

func Gopt_AppGen_Main(app interface{ initApp() }, workers ...interface{ initClass(self any) }) {
	flag.BoolVar(&debugExec, "v", false, "verbose infromation")
	flag.Parse()
	app.initApp()
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	for _, worker := range workers {
		worker.initClass(worker)
		worker.(interface{ Main() }).Main()
	}
}

// -----------------------------------------------------------------------------
