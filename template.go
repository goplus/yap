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

package yap

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/yap/internal/templ"
)

type Delims struct {
	Left  string
	Right string
}
type Template struct {
	*template.Template
	fsys   fs.FS
	delims Delims
}

func NewTemplate(name string) *Template {
	return &Template{template.New(name), nil, Delims{"{{", "}}"}}
}
func (t *Template) NewTemplate(name string) *Template {
	return &Template{Template: t.Template.New(name), fsys: t.fsys, delims: Delims{"{{", "}}"}}
}

func (t Template) Parse(text string) (ret Template, err error) {
	ret.Template, err = t.Template.Parse(templ.Translate(text, t.delims.Left, t.delims.Right))
	return
}

func (t *Template) SetDelims(delims Delims) {
	t.delims = delims
	t.Delims(t.delims.Left, t.delims.Right)
}

func ParseFSFile(f fs.FS, file string) (t Template, err error) {
	b, err := fs.ReadFile(f, file)
	if err != nil {
		return
	}
	name := filepath.Base(file)
	return NewTemplate(name).Parse(string(b))
}

func ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(nil, readFileOS, filenames...)
}

func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(t, readFileOS, filenames...)
}

func parseFiles(t *Template, readFile func(string) (string, []byte, error), filenames ...string) (*Template, error) {

	if len(filenames) == 0 {
		return nil, fmt.Errorf("yap/template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		name, b, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		var tmpl *Template
		if t == nil {
			t = NewTemplate(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.NewTemplate(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func ParseGlob(pattern string) (*Template, error) {
	return parseGlob(nil, pattern)
}

func (t *Template) ParseGlob(pattern string) (*Template, error) {
	return parseGlob(t, pattern)
}

func parseGlob(t *Template, pattern string) (*Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: pattern matches no files: %#q", pattern)
	}
	return parseFiles(t, readFileOS, filenames...)
}

// IsTrue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value. This is the definition of
// truth used by if and other such actions.
func IsTrue(val any) (truth, ok bool) {
	return template.IsTrue(val)
}

// ParseFS is like ParseFiles or ParseGlob but reads from the file system fs
// instead of the host operating system's file system.
// It accepts a list of glob patterns.
// (Note that most file names serve as glob patterns matching only themselves.)
func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	return parseFS(nil, fs, patterns)
}

// ParseFS is like ParseFiles or ParseGlob but reads from the file system fs
// instead of the host operating system's file system.
// It accepts a list of glob patterns.
// (Note that most file names serve as glob patterns matching only themselves.)
func (t *Template) ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	return parseFS(t, fs, patterns)
}

func parseFS(t *Template, fsys fs.FS, patterns []string) (*Template, error) {
	var filenames []string
	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		filenames = append(filenames, list...)
	}
	return parseFiles(t, readFileFS(fsys), filenames...)
}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

func readFileFS(fsys fs.FS) func(string) (string, []byte, error) {
	return func(file string) (name string, b []byte, err error) {
		name = path.Base(file)
		name = strings.TrimSuffix(name, "_yap.html")
		b, err = fs.ReadFile(fsys, file)
		return
	}
}
