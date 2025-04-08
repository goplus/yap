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

package ytest

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/goplus/yap/ytest/auth"
	"github.com/goplus/yap/ytest/auth/bearer"
	"github.com/qiniu/x/test"
)

// -----------------------------------------------------------------------------

// Bearer creates a Bearer Authorization by specified token.
func Bearer(token string) auth.RTComposer {
	return bearer.New(token)
}

// -----------------------------------------------------------------------------

// JsonEncode encodes a value into string in json format.
func JsonEncode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		test.Fatal("json.Marshal failed:", err)
	}
	return string(b)
}

// Form encodes a map value into string in form format.
func Form(m map[string]any) (ret url.Values) {
	ret = make(url.Values)
	for k, v := range m {
		ret[k] = formVal(v)
	}
	return ret
}

func formVal(val any) []string {
	switch v := val.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case int:
		return []string{strconv.Itoa(v)}
	case bool:
		if v {
			return []string{"true"}
		}
		return []string{"false"}
	case float64:
		return []string{strconv.FormatFloat(v, 'g', -1, 64)}
	}
	test.Fatalf("formVal unexpected type: %T\n", val)
	return nil
}

// -----------------------------------------------------------------------------
