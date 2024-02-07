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

package sqlite3

import (
	"github.com/goplus/yap/ydb"
	"github.com/mattn/go-sqlite3"
)

func wrapErr(prompt string, err error) error {
	if prompt == "insert:" {
		if e, ok := err.(sqlite3.Error); ok && e.Code == sqlite3.ErrConstraint {
			return ydb.ErrDuplicated
		}
	}
	return err
}

func init() {
	ydb.Register(&ydb.Engine{
		Name:       "sqlite3",
		TestSource: "file:test.db?cache=shared&mode=memory",
		WrapErr:    wrapErr,
	})
}
