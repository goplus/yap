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

import "net/http"

type tokenAuth struct {
	// Type     string
	Token string
}

type tokenRounderTripper struct {
	RoundTripper http.RoundTripper
	Token        string
}

func (rt *tokenRounderTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", rt.Token)
	return rt.RoundTripper.RoundTrip(req)
}

func (p *tokenAuth) Compose(rt http.RoundTripper) http.RoundTripper {
	return &tokenRounderTripper{
		RoundTripper: rt,
		Token:        p.Token,
	}
}

// -----------------------------------------------------------------------------
