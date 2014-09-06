// This package contains the base logic for the creation and rendering of field widgets. Base widgets are defined for most input fields,
// both in classic and Bootstrap3 style; custom widgets can be defined and associated to a field, provided that they implement the
// WidgetInterface interface.
package widgets

import (
	"bytes"
	"github.com/coscms/forms/common"
	"html/template"
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
func BaseWidget(style, inputType string) *Widget {
	var fpath string = formcommon.TmplDir + "/" + style + "/"
	var urls []string = []string{formcommon.CreateUrl(fpath + "generic.tmpl")}
	switch inputType {
	case formcommon.BUTTON:
		urls = append(urls, formcommon.CreateUrl(fpath+"button.html"))
	case formcommon.CHECKBOX:
		urls = append(urls, formcommon.CreateUrl(fpath+"options/checkbox.html"))
	case formcommon.TEXTAREA:
		urls = append(urls, formcommon.CreateUrl(fpath+"text/textareainput.html"))
	case formcommon.SELECT:
		urls = append(urls, formcommon.CreateUrl(fpath+"options/select.html"))
	case formcommon.PASSWORD:
		urls = append(urls, formcommon.CreateUrl(fpath+"text/passwordinput.html"))
	case formcommon.RADIO:
		urls = append(urls, formcommon.CreateUrl(fpath+"options/radiobutton.html"))
	case formcommon.TEXT:
		urls = append(urls, formcommon.CreateUrl(fpath+"text/textinput.html"))
	case formcommon.RANGE:
		urls = append(urls, formcommon.CreateUrl(fpath+"number/range.html"))
	case formcommon.NUMBER:
		urls = append(urls, formcommon.CreateUrl(fpath+"number/number.html"))
	case formcommon.RESET:
		urls = append(urls, formcommon.CreateUrl(fpath+"button.html"))
	case formcommon.SUBMIT:
		urls = append(urls, formcommon.CreateUrl(fpath+"button.html"))
	case formcommon.DATE:
		urls = append(urls, formcommon.CreateUrl(fpath+"datetime/date.html"))
	case formcommon.DATETIME:
		urls = append(urls, formcommon.CreateUrl(fpath+"datetime/datetime.html"))
	case formcommon.TIME:
		urls = append(urls, formcommon.CreateUrl(fpath+"datetime/time.html"))
	case formcommon.DATETIME_LOCAL:
		urls = append(urls, formcommon.CreateUrl(fpath+"datetime/datetime.html"))
	case formcommon.STATIC:
		urls = append(urls, formcommon.CreateUrl(fpath+"static.html"))
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
		urls = append(urls, formcommon.CreateUrl(fpath+"input.html"))
	default:
		urls = append(urls, formcommon.CreateUrl(fpath+"input.html"))
	}
	templ := template.Must(template.ParseFiles(urls...))
	return &Widget{templ}
}
