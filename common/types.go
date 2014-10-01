// This package provides basic constants used by go-form-it packages.
package formcommon

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/coscms/tagfast"
)

var (
	TmplDir string              = "templates"
	LabelFn func(string) string = func(s string) string {
		return s
	}

	//private
	cachedTemplate map[string]*template.Template = make(map[string]*template.Template)
	lock           *sync.RWMutex                 = new(sync.RWMutex)
)

const (
	PACKAGE_NAME = "github.com/coscms/forms"
)

// Input field types
const (
	BUTTON         = "button"
	CHECKBOX       = "checkbox"
	COLOR          = "color" // Not yet implemented
	DATE           = "date"
	DATETIME       = "datetime"
	DATETIME_LOCAL = "datetime-local"
	EMAIL          = "email" // Not yet implemented
	FILE           = "file"  // Not yet implemented
	HIDDEN         = "hidden"
	IMAGE          = "image" // Not yet implemented
	MONTH          = "month" // Not yet implemented
	NUMBER         = "number"
	PASSWORD       = "password"
	RADIO          = "radio"
	RANGE          = "range"
	RESET          = "reset"
	SEARCH         = "search" // Not yet implemented
	SUBMIT         = "submit"
	TEL            = "tel" // Not yet implemented
	TEXT           = "text"
	TIME           = "time"
	URL            = "url"  // Not yet implemented
	WEEK           = "week" // Not yet implemented
	TEXTAREA       = "textarea"
	SELECT         = "select"
	STATIC         = "static"
)

// Available form styles
const (
	BASE      = "base"
	BOOTSTRAP = "bootstrap3"
)

// CreateUrl creates the complete url of the desired widget template
func CreateUrl(widget string) string {
	//println(widget)
	if _, err := os.Stat(widget); os.IsNotExist(err) {
		return path.Join(os.Getenv("GOPATH"), "src", PACKAGE_NAME, widget)
	}
	return widget
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

func ParseTmpl(data interface{}, fn_tpl template.FuncMap, fn_fixTpl func(tpls ...string) ([]string, error), tpls ...string) string {
	var s string
	buf := bytes.NewBufferString(s)
	tpf := fmt.Sprintf("%v", tpls)
	tpl, ok := CachedTemplate(tpf)
	if !ok {
		c := template.New(path.Base(tpls[0]))
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

func Tag(t reflect.Type, fieldNo int, tagName string) string {
	return tagfast.Tag1(t, fieldNo, tagName)
}
