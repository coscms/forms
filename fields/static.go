package fields

import (
	"fmt"
	"github.com/coscms/forms/common"
	"reflect"
)

// // Static field type
// type StaticType struct {
// 	Field
// }

// StaticField returns a static field with the provided name and content
func StaticField(name, content string) *Field {
	ret := FieldWithType(name, formcommon.STATIC)
	ret.SetText(content)
	return ret
}

// RadioFieldFromInstance creates and initializes a radio field based on its name, the reference object instance and field number.
// It uses i object's [fieldNo]-th field content (if any) to set the field content.
func StaticFieldFromInstance(i interface{}, fieldNo int, name string) *Field {
	ret := StaticField(name, fmt.Sprintf("%s", reflect.ValueOf(i).Field(fieldNo).Interface()))
	return ret
}
