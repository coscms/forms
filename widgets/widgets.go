// This package contains the base logic for the creation and rendering of field widgets. Base widgets are defined for most input fields,
// both in classic and Bootstrap3 style; custom widgets can be defined and associated to a field, provided that they implement the
// WidgetInterface interface.
package widgets

import (
	"bytes"
	"html/template"

	"github.com/coscms/forms/common"
)

// Simple widget object that gets executed at render time.
type Widget struct {
	template *template.Template
}

// WidgetInterface defines the requirements for custom widgets.
type WidgetInterface interface {
	Render(data interface{}) string
}

// Render executes the internal template and returns the result as a template.HTML object.
func (w *Widget) Render(data interface{}) string {
	var s string
	buf := bytes.NewBufferString(s)
	w.template.ExecuteTemplate(buf, "main", data)
	return buf.String()
}

// BaseWidget creates a Widget based on style and inpuType parameters, both defined in the common package.
func BaseWidget(style, inputType, tmplName string) *Widget {
	var cachedKey string = style+", "+inputType+", "+tmplName
	templ, ok := formcommon.CachedTemplate(cachedKey)
	if !ok {
		var (
			fpath string = formcommon.TmplDir + "/" + style + "/"
			urls []string = []string{formcommon.CreateUrl(fpath + "generic.tmpl")}
			tpath string = widgetTmpl(inputType, tmplName)
		)
		urls = append(urls, formcommon.CreateUrl(fpath+tpath+".html"))
		templ = template.Must(template.ParseFiles(urls...))
		formcommon.SetCachedTemplate(cachedKey, templ)
	}
	return &Widget{templ}
}

func widgetTmpl(inputType, tmpl string) (tpath string) {
	switch inputType {
	case formcommon.BUTTON:
		tpath = "button"
		if tmpl != "" {
			tpath = tmpl
		}
	case formcommon.TEXTAREA:
		tpath = "text/textareainput"
		if tmpl != "" {
			tpath = "text/" + tmpl
		}
	case formcommon.PASSWORD:
		tpath = "text/passwordinput"
		if tmpl != "" {
			tpath = "text/" + tmpl
		}
	case formcommon.TEXT:
		tpath = "text/textinput"
		if tmpl != "" {
			tpath = "text/" + tmpl
		}
	case formcommon.CHECKBOX:
		tpath = "options/checkbox"
		if tmpl != "" {
			tpath = "options/" + tmpl
		}
	case formcommon.SELECT:
		tpath = "options/select"
		if tmpl != "" {
			tpath = "options/" + tmpl
		}
	case formcommon.RADIO:
		tpath = "options/radiobutton"
		if tmpl != "" {
			tpath = "options/" + tmpl
		}
	case formcommon.RANGE:
		tpath = "number/range"
		if tmpl != "" {
			tpath = "number/" + tmpl
		}
	case formcommon.NUMBER:
		tpath = "number/number"
		if tmpl != "" {
			tpath = "number/" + tmpl
		}
	case formcommon.RESET, formcommon.SUBMIT:
		tpath = "button"
		if tmpl != "" {
			tpath = tmpl
		}
	case formcommon.DATE:
		tpath = "datetime/date"
		if tmpl != "" {
			tpath = "datetime/" + tmpl
		}
	case formcommon.DATETIME:
		tpath = "datetime/datetime"
		if tmpl != "" {
			tpath = "datetime/" + tmpl
		}
	case formcommon.TIME:
		tpath = "datetime/time"
		if tmpl != "" {
			tpath = "datetime/" + tmpl
		}
	case formcommon.DATETIME_LOCAL:
		tpath = "datetime/datetime"
		if tmpl != "" {
			tpath = "datetime/" + tmpl
		}
	case formcommon.STATIC:
		tpath = "static"
		if tmpl != "" {
			tpath = tmpl
		}
	case formcommon.SEARCH,
		formcommon.TEL,
		formcommon.URL,
		formcommon.WEEK,
		formcommon.COLOR,
		formcommon.EMAIL,
		formcommon.FILE,
		formcommon.HIDDEN,
		formcommon.IMAGE,
		formcommon.MONTH:
		fallthrough
	default:
		tpath = "input"
		if tmpl != "" {
			tpath = tmpl
		}
	}
	return
}
