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

import "github.com/goplus/gop/ast"

// -----------------------------------------------------------------------------

// Delete deltes rows by cond.
//   - delete <cond>, <arg1>, <arg2>, ...
func (p *Class) Delete__0(src ast.Expr, cond string, args ...any) error {
	panic("todo")
}

// Delete deltes rows by cond.
//   - delete <cond>, <arg1>, <arg2>, ...
func (p *Class) Delete__1(cond string, args ...any) error {
	return p.Delete__0(nil, cond, args...)
}

// -----------------------------------------------------------------------------
