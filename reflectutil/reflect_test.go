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

package reflectutil

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestSetZero(t *testing.T) {
	a := 2
	v := reflect.ValueOf(&a).Elem()
	SetZero(v)
	if a != 0 {
		t.Fatal("SetZero:", a)
	}
}

func TestUnsafeAddr(t *testing.T) {
	if unsafe.Sizeof(value{}) != unsafe.Sizeof(reflect.Value{}) {
		panic("unexpected sizeof reflect.Value")
	}
	v := reflect.ValueOf(0)
	if UnsafeAddr(v) == 0 {
		t.Fatal("UnsafeAddr")
	}
}
