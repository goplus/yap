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

func parseFS(fsys fs.FS, delimLeft, delimRight, suffix string) (*template.Template, error) {
	if delimLeft == "" {
		delimLeft = "{{"
	}
	if delimRight == "" {
		delimRight = "}}"
	}
	t := template.New("").Delims(delimLeft, delimRight)

	pattern := "*" + suffix
	filenames, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	for _, filename := range filenames {
		content, err := fs.ReadFile(fsys, filename)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if templ.TranslateEx(&buf, string(content), delimLeft, delimRight) {
			content = buf.Bytes()
		}

		name := strings.TrimSuffix(path.Base(filename), suffix)
		if _, err := t.New(name).Parse(string(content)); err != nil {
			return nil, err
		}
	}
	return t, nil
}
