/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

//Package common This package provides basic constants used by forms packages.
package common

import (
	"html/template"
	"path/filepath"

	"github.com/webx-top/echo/middleware/tplfunc"
)

func TplFuncs() template.FuncMap {
	return tplfunc.TplFuncMap
}

func ParseFiles(files ...string) (*template.Template,error) {
	name := filepath.Base(files[0])
	b, err := FileReader(files[0])
	if err != nil {
		return nil, err
	}
	tmpl := template.New(name)
	tmpl.Funcs(TplFuncs())
	tmpl = template.Must(tmpl.Parse(string(b)))
	if len(files) > 1 {
		tmpl = template.Must(tmpl.ParseFiles(files[1:]...))
	}
	return tmpl, nil
}
