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

package auth

import (
	"net/http"
)

// RTComposer represents an abstract of http Authorization objects.
type RTComposer interface {
	Compose(base http.RoundTripper) http.RoundTripper
}

// -----------------------------------------------------------------------------

type tokenRT struct {
	rt    http.RoundTripper
	token string
}

func (p *tokenRT) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", p.token)
	return p.rt.RoundTrip(req)
}

func WithToken(rt http.RoundTripper, token string) http.RoundTripper {
	return &tokenRT{
		rt:    rt,
		token: token,
	}
}

// -----------------------------------------------------------------------------
