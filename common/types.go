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

package common

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/coscms/forms/config"
	"github.com/webx-top/tagfast"
)

// Available form styles
const (
	BASE      = "base"
	BOOTSTRAP = "bootstrap3"
)

var (
	tmplDirs = map[string]string{BASE: "templates", BOOTSTRAP: "templates"}
	LabelFn  = func(s string) string {
		return s
	}
	FileReader = ioutil.ReadFile

	//private
	cachedTemplate = make(map[string]*template.Template)
	cachedConfig   = make(map[string]*config.Config)
	lock           = new(sync.RWMutex)
)

const (
	PACKAGE_NAME = "github.com/coscms/forms"
)

// Input field types
const (
	BUTTON         = "button"
	CHECKBOX       = "checkbox"
	COLOR          = "color"
	DATE           = "date"
	DATETIME       = "datetime"
	DATETIME_LOCAL = "datetime-local"
	EMAIL          = "email"
	FILE           = "file"
	HIDDEN         = "hidden"
	IMAGE          = "image"
	MONTH          = "month"
	NUMBER         = "number"
	PASSWORD       = "password"
	RADIO          = "radio"
	RANGE          = "range"
	RESET          = "reset"
	SEARCH         = "search"
	SUBMIT         = "submit"
	TEL            = "tel"
	TEXT           = "text"
	TIME           = "time"
	URL            = "url"
	WEEK           = "week"
	TEXTAREA       = "textarea"
	SELECT         = "select"
	STATIC         = "static"
)

func SetTmplDir(style, tmplDir string) {
	tmplDirs[style] = tmplDir
}

func TmplDir(style string) (tmplDir string) {
	tmplDir, _ = tmplDirs[style]
	return
}

// CreateUrl creates the complete url of the desired widget template
func CreateUrl(widget string) string {
	if !TmplExists(widget) {
		return filepath.Join(os.Getenv("GOPATH"), "src", PACKAGE_NAME, widget)
	}
	return widget
}

func TmplExists(tmpl string) bool {
	_, err := os.Stat(tmpl)
	return !os.IsNotExist(err)
}

func CachedTemplate(cachedKey string) (r *template.Template, ok bool) {
	lock.RLock()
	defer lock.RUnlock()

	r, ok = cachedTemplate[cachedKey]
	return
}

func SetCachedTemplate(cachedKey string, tmpl *template.Template) bool {
	lock.Lock()
	defer lock.Unlock()

	cachedTemplate[cachedKey] = tmpl
	return true
}

func ClearCachedTemplate() {
	cachedTemplate = make(map[string]*template.Template)
}

func DelCachedTemplate(key string) bool {
	if _, ok := cachedTemplate[key]; ok {
		delete(cachedTemplate, key)
		return true
	}
	return false
}

func CachedConfig(cachedKey string) (r *config.Config, ok bool) {
	lock.RLock()
	defer lock.RUnlock()

	r, ok = cachedConfig[cachedKey]
	return
}

func SetCachedConfig(cachedKey string, c *config.Config) bool {
	lock.Lock()
	defer lock.Unlock()

	cachedConfig[cachedKey] = c
	return true
}

func ClearCachedConfig() {
	cachedConfig = make(map[string]*config.Config)
}

func DelCachedConfig(key string) bool {
	if _, ok := cachedConfig[key]; ok {
		delete(cachedConfig, key)
		return true
	}
	return false
}

func ParseTmpl(data interface{}, fn_tpl template.FuncMap, fn_fixTpl func(tpls ...string) ([]string, error), tpls ...string) string {
	var s string
	buf := bytes.NewBufferString(s)
	tpf := fmt.Sprintf("%v", tpls)
	tpl, ok := CachedTemplate(tpf)
	if !ok {
		c := template.New(filepath.Base(tpls[0]))
		if fn_tpl != nil {
			c.Funcs(fn_tpl)
		}
		if fn_fixTpl != nil {
			var err error
			tpls, err = fn_fixTpl(tpls...)
			if err != nil {
				return fmt.Sprintf(`%v`, err)
			}
		}
		tpl = template.Must(c.ParseFiles(tpls...))
		SetCachedTemplate(tpf, tpl)
	}
	err := tpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func TagVal(t reflect.Type, fieldNo int, tagName string) string {
	return tagfast.Value(t, t.Field(fieldNo), tagName)
}

func Tag(t reflect.Type, f reflect.StructField, tagName string) (value string, tf tagfast.Faster) {
	return tagfast.Tag(t, f, tagName)
}
