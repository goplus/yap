/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package htmltempl

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"path"
	"strings"
	_ "unsafe"

	"github.com/goplus/yap/internal/templ"
)

// Template is the representation of a parsed template. The *parse.Tree
// field is exported only for use by html/template and should be treated
// as unexported by all other clients.
type Template struct {
	*template.Template
}

func (p *Template) InitTemplates(fsys fs.FS, delimLeft, delimRight, suffix string) {
	tpl, err := parseFS(fsys, delimLeft, delimRight, suffix)
	if err != nil {
		log.Panicln(err)
	}
	p.Template = tpl
}

//go:linkname parseFiles html/template.parseFiles
func parseFiles(t *template.Template, readFile func(string) (string, []byte, error), filenames ...string) (*template.Template, error)

func parseFS(fsys fs.FS, delimLeft, delimRight, suffix string) (*template.Template, error) {
	pattern := "*" + suffix
	filenames, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	if delimLeft == "" {
		delimLeft = "{{"
	}
	if delimRight == "" {
		delimRight = "}}"
	}
	t := template.New("").Delims(delimLeft, delimRight)
	return parseFiles(t, readFileFS(fsys, delimLeft, delimRight, suffix), filenames...)
}

func readFileFS(fsys fs.FS, delimLeft, delimRight, suffix string) func(string) (string, []byte, error) {
	return func(file string) (name string, b []byte, err error) {
		name = strings.TrimSuffix(path.Base(file), suffix)
		if b, err = fs.ReadFile(fsys, file); err != nil {
			return
		}
		var buf bytes.Buffer
		if templ.TranslateEx(&buf, string(b), delimLeft, delimRight) {
			b = buf.Bytes()
		}
		return
	}
}
