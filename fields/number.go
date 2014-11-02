package fields

import (
	"fmt"
	"reflect"

	"github.com/coscms/forms/common"
)

// RangeField creates a default range field with the provided name. Min, max and step parameters define the expected behavior
// of the HTML field.
func RangeField(name string, min, max, step int) *Field {
	ret := FieldWithType(name, formcommon.RANGE)
	ret.SetParam("min", fmt.Sprintf("%d", min))
	ret.SetParam("max", fmt.Sprintf("%d", max))
	ret.SetParam("step", fmt.Sprintf("%d", step))
	return ret
}

// NumberField craetes a default number field with the provided name.
func NumberField(name string) *Field {
	ret := FieldWithType(name, formcommon.NUMBER)
	return ret
}

// NumberFieldFromInstance creates and initializes a number field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max" and "form_value" tags to add additional parameters to the field.
func NumberFieldFromInstance(val reflect.Value,t reflect.Type, fieldNo int, name string) *Field {
	ret := NumberField(name)
	// check tags
	if v := formcommon.Tag(t, fieldNo, "form_min"); v != "" {
		ret.SetParam("min", v)
	}
	if v := formcommon.Tag(t, fieldNo, "form_max"); v != "" {
		ret.SetParam("max", v)
	}
	if v := formcommon.Tag(t, fieldNo, "form_value"); v != "" {
		ret.SetValue(v)
	} else {
		ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
	}
	return ret
}

// RangeFieldFromInstance creates and initializes a range field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max", "form_step" and "form_value" tags to add additional parameters to the field.
func RangeFieldFromInstance(val reflect.Value,t reflect.Type, fieldNo int, name string) *Field {
	ret := RangeField(name, 0, 10, 1)
	// check tags
	if v := formcommon.Tag(t, fieldNo, "form_min"); v != "" {
		ret.SetParam("min", v)
	}
	if v := formcommon.Tag(t, fieldNo, "form_max"); v != "" {
		ret.SetParam("max", v)
	}
	if v := formcommon.Tag(t, fieldNo, "form_step"); v != "" {
		ret.SetParam("step", v)
	}
	if v := formcommon.Tag(t, fieldNo, "form_value"); v != "" {
		ret.SetValue(v)
	} else {
		ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
	}
	return ret
}
