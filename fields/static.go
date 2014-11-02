package fields

import (
	"fmt"
	"github.com/coscms/forms/common"
	"reflect"
)

// StaticField returns a static field with the provided name and content
func StaticField(name, content string) *Field {
	ret := FieldWithType(name, formcommon.STATIC)
	ret.SetText(content)
	return ret
}

// RadioFieldFromInstance creates and initializes a radio field based on its name, the reference object instance and field number.
func StaticFieldFromInstance(val reflect.Value,t reflect.Type, fieldNo int, name string) *Field {
	ret := StaticField(name, fmt.Sprintf("%s", val.Field(fieldNo).Interface()))
	return ret
}
