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
	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

type Class struct {
}

func newClass(name string) *Class {
	return &Class{}
}

func (p *Class) Use(table string) {
}

func (p *Class) Api(name string, creator func(), fnlit ...*ast.FuncLit) {
}

func (p *Class) Call(args ...any) {
}

func (p *Class) Ret(args ...any) {
}

// -----------------------------------------------------------------------------
