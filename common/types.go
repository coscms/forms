// This package provides basic constants used by go-form-it packages.
package formcommon

import (
	"os"
	"path"
	"reflect"

	"github.com/coscms/tagfast"
)

var TmplDir string = "templates"

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
	if _, err := os.Stat(widget); os.IsNotExist(err) {
		return path.Join(os.Getenv("GOPATH"), "src", PACKAGE_NAME, widget)
	}
	return widget
}

func Tag(t reflect.Type, fieldNo int, tagName string) string {
	return tagfast.Tag1(t, fieldNo, tagName)
}
