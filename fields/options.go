package fields

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/coscms/forms/common"
)

// Id - Value pair used to define an option for select and redio input fields.
type InputChoice struct {
	Id, Val string
	Checked bool
}

type ChoiceIndex struct {
	Group string
	Index int
}

// // Radio button type.
// type RadioType struct {
// 	Field
// }

// // Select field type.
// type SelectType struct {
// 	Field
// }

// // Checkbox field type.
// type CheckBoxType struct {
// 	Field
// }

// =============== RADIO

// RadioField creates a default radio button input field with the provided name and list of choices.
func RadioField(name string, choices []InputChoice) *Field {
	ret := FieldWithType(name, formcommon.RADIO)
	ret.choices = []InputChoice{}
	ret.SetChoices(choices)
	return ret
}

// RadioFieldFromInstance creates and initializes a radio field based on its name, the reference object instance and field number.
// This method looks for "form_choices" and "form_value" tags to add additional parameters to the field. "form_choices" tag is a list
// of <id>|<value> options, joined by "|" character; ex: "A|Option A|B|Option B" translates into 2 options: <A, Option A> and <B, Option B>.
// It also uses i object's [fieldNo]-th field content (if any) to override the "form_value" option and fill the HTML field.
func RadioFieldFromInstance(i interface{}, fieldNo int, name string) *Field {
	t := reflect.TypeOf(i)
	choices := strings.Split(formcommon.Tag(t, fieldNo, "form_choices"), "|")
	chArr := make([]InputChoice, 0)
	ret := RadioField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-1; i += 2 {
		ret.choiceKeys[choices[i]] = ChoiceIndex{Group: "", Index: len(chArr)}
		chArr = append(chArr, InputChoice{choices[i], formcommon.LabelFn(choices[i+1]), false})
		chMap[choices[i]] = choices[i+1]
	}
	ret.SetChoices(chArr)
	var v string = formcommon.Tag(t, fieldNo, "form_value")
	if v == "" {
		v = fmt.Sprintf("%s", reflect.ValueOf(i).Field(fieldNo).String())
	}
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// ================ SELECT

// SelectField creates a default select input field with the provided name and map of choices. Choices for SelectField are grouped
// by name (if <optgroup> is needed); "" group is the default one and does not trigger a <optgroup></optgroup> rendering.
func SelectField(name string, choices map[string][]InputChoice) *Field {
	ret := FieldWithType(name, formcommon.SELECT)
	ret.choices = map[string][]InputChoice{}
	ret.SetChoices(choices)
	return ret
}

// SelectFieldFromInstance creates and initializes a select field based on its name, the reference object instance and field number.
// This method looks for "form_choices" and "form_value" tags to add additional parameters to the field. "form_choices" tag is a list
// of <group<|<id>|<value> options, joined by "|" character; ex: "G1|A|Option A|G1|B|Option B" translates into 2 options in the same group G1:
// <A, Option A> and <B, Option B>. "" group is the default one.
// It also uses i object's [fieldNo]-th field content (if any) to override the "form_value" option and fill the HTML field.
func SelectFieldFromInstance(i interface{}, fieldNo int, name string, options map[string]struct{}) *Field {
	t := reflect.TypeOf(i)
	choices := strings.Split(formcommon.Tag(t, fieldNo, "form_choices"), "|")
	chArr := make(map[string][]InputChoice)
	ret := SelectField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-2; i += 3 {
		optgroupLabel := formcommon.LabelFn(choices[i])
		if _, ok := chArr[optgroupLabel]; !ok {
			chArr[optgroupLabel] = make([]InputChoice, 0)
		}
		id := choices[i+1]
		ret.choiceKeys[id] = ChoiceIndex{Group: optgroupLabel, Index: len(chArr[optgroupLabel])}
		chArr[optgroupLabel] = append(chArr[optgroupLabel], InputChoice{id, formcommon.LabelFn(choices[i+2]), false})
		chMap[id] = choices[i+2]
	}
	ret.SetChoices(chArr)

	if _, ok := options["multiple"]; ok {
		ret.MultipleChoice()
	}

	var v string = fmt.Sprintf("%s", reflect.ValueOf(i).Field(fieldNo).String())
	if v == "" {
		// TODO: multiple value parsing
		v = formcommon.Tag(t, fieldNo, "form_value")
	}
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// ================== CHECKBOX

func CheckboxField(name string, choices []InputChoice) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)
	ret.choices = []InputChoice{}
	ret.SetChoices(choices)
	if len(ret.choices.([]InputChoice)) > 1 {
		ret.MultipleChoice()
	}
	return ret
}

func CheckboxFieldFromInstance(i interface{}, fieldNo int, name string) *Field {
	t := reflect.TypeOf(i)
	choices := strings.Split(formcommon.Tag(t, fieldNo, "form_choices"), "|")
	chArr := make([]InputChoice, 0)
	ret := CheckboxField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-1; i += 2 {
		ret.choiceKeys[choices[i]] = ChoiceIndex{Group: "", Index: len(chArr)}
		chArr = append(chArr, InputChoice{choices[i], formcommon.LabelFn(choices[i+1]), false})
		chMap[choices[i]] = choices[i+1]
	}
	ret.SetChoices(choices)
	if len(ret.choices.([]InputChoice)) > 1 {
		ret.MultipleChoice()
	}

	var v string = formcommon.Tag(t, fieldNo, "form_value")
	if v == "" {
		v = fmt.Sprintf("%s", reflect.ValueOf(i).Field(fieldNo).String())
	}
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// Checkbox creates a default checkbox field with the provided name. It also makes it checked by default based
// on the checked parameter.
func Checkbox(name string, checked bool) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)
	if checked {
		ret.AddTag("checked")
	}
	return ret
}

// CheckboxFromInstance creates and initializes a checkbox field based on its name, the reference object instance, field number and field options.
// It uses i object's [fieldNo]-th field content (if any) to override the "checked" option in the options map and check the field.
func CheckboxFromInstance(i interface{}, fieldNo int, name string, options map[string]struct{}) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)

	if _, ok := options["checked"]; ok {
		ret.AddTag("checked")
	} else {
		val := reflect.ValueOf(i).Field(fieldNo).Bool()
		if val {
			ret.AddTag("checked")
		}
	}
	return ret
}
