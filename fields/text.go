package fields

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/coscms/forms/common"
)

func ColorField(name string) *Field {
	return FieldWithType(name, formcommon.COLOR)
}

func EmailField(name string) *Field {
	return FieldWithType(name, formcommon.EMAIL)
}

func FileField(name string) *Field {
	return FieldWithType(name, formcommon.FILE)
}

func ImageField(name string) *Field {
	return FieldWithType(name, formcommon.IMAGE)
}

func MonthField(name string) *Field {
	return FieldWithType(name, formcommon.MONTH)
}

func SearchField(name string) *Field {
	return FieldWithType(name, formcommon.SEARCH)
}

func TelField(name string) *Field {
	return FieldWithType(name, formcommon.TEL)
}

func UrlField(name string) *Field {
	return FieldWithType(name, formcommon.URL)
}

func WeekField(name string) *Field {
	return FieldWithType(name, formcommon.WEEK)
}

// TextField creates a default text input field based on the provided name.
func TextField(name string, typ ...string) *Field {
	var t = formcommon.TEXT
	if len(typ) > 0 {
		t = typ[0]
	}
	return FieldWithType(name, t)
}

// PasswordField creates a default password text input field based on the provided name.
func PasswordField(name string) *Field {
	return FieldWithType(name, formcommon.PASSWORD)
}

// =========== TEXT AREA

// TextAreaField creates a default textarea input field based on the provided name and dimensions.
func TextAreaField(name string, rows, cols int) *Field {
	ret := FieldWithType(name, formcommon.TEXTAREA)
	ret.SetParam("rows", fmt.Sprintf("%d", rows))
	ret.SetParam("cols", fmt.Sprintf("%d", cols))
	return ret
}

// ========================

// HiddenField creates a default hidden input field based on the provided name.
func HiddenField(name string) *Field {
	return FieldWithType(name, formcommon.HIDDEN)
}

// TextFieldFromInstance creates and initializes a text field based on its name, the reference object instance and field number.
func TextFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, typ ...string) *Field {
	ret := TextField(name, typ...)
	ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
	return ret
}

// PasswordFieldFromInstance creates and initializes a password field based on its name, the reference object instance and field number.
func PasswordFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	ret := PasswordField(name)
	ret.SetValue(fmt.Sprintf("%s", val.Field(fieldNo).String()))
	return ret
}

// TextFieldFromInstance creates and initializes a text field based on its name, the reference object instance and field number.
// This method looks for "form_rows" and "form_cols" tags to add additional parameters to the field.
func TextAreaFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	var rows, cols int = 20, 50
	var err error
	if v := formcommon.Tag(t, fieldNo, "form_rows"); v != "" {
		rows, err = strconv.Atoi(v)
		if err != nil {
			return nil
		}
	}
	if v := formcommon.Tag(t, fieldNo, "form_cols"); v != "" {
		cols, err = strconv.Atoi(v)
		if err != nil {
			return nil
		}
	}
	ret := TextAreaField(name, rows, cols)
	ret.SetText(fmt.Sprintf("%s", val.Field(fieldNo).String()))
	return ret
}

// HiddenFieldFromInstance creates and initializes a hidden field based on its name, the reference object instance and field number.
func HiddenFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	ret := HiddenField(name)
	ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
	return ret
}
