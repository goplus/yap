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

package jwt

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goplus/yap/ytest/auth"
)

// -----------------------------------------------------------------------------

type Signaturer struct {
	method jwt.SigningMethod
	claims jwt.MapClaims
	secret any
}

func (p *Signaturer) Set(k string, v any) *Signaturer {
	p.claims[k] = v
	return p
}

func (p *Signaturer) Audience(aud ...string) *Signaturer {
	p.claims["aud"] = jwt.ClaimStrings(aud)
	return p
}

func (p *Signaturer) Expiration(exp time.Time) *Signaturer {
	p.claims["exp"] = jwt.NewNumericDate(exp)
	return p
}

func (p *Signaturer) NotBefore(nbf time.Time) *Signaturer {
	p.claims["nbf"] = jwt.NewNumericDate(nbf)
	return p
}

func (p *Signaturer) IssuedAt(iat time.Time) *Signaturer {
	p.claims["iat"] = jwt.NewNumericDate(iat)
	return p
}

func (p *Signaturer) Compose(rt http.RoundTripper) http.RoundTripper {
	token := jwt.NewWithClaims(p.method, p.claims)
	raw, err := token.SignedString(p.secret)
	if err != nil {
		log.Panicln("jwt token.SignedString:", err)
	}
	return auth.WithToken(rt, "Bearer "+raw)
}

func newSign(method jwt.SigningMethod, secret any) *Signaturer {
	return &Signaturer{
		method: method,
		claims: make(jwt.MapClaims),
		secret: secret,
	}
}

// HS256 creates a signing methods by using the HMAC-SHA256.
func HS256(key []byte) *Signaturer {
	return newSign(jwt.SigningMethodHS256, key)
}

// HS384 creates a signing methods by using the HMAC-SHA384.
func HS384(key []byte) *Signaturer {
	return newSign(jwt.SigningMethodHS384, key)
}

// HS512 creates a signing methods by using the HMAC-SHA512.
func HS512(key []byte) *Signaturer {
	return newSign(jwt.SigningMethodHS512, key)
}

// -----------------------------------------------------------------------------
