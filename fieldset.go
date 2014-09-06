package forms

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/fields"
	"html/template"
)

// FieldSetType is a collection of fields grouped within a form.
type FieldSetType struct {
	tmpl     string
	name     string
	class    map[string]struct{}
	tags     map[string]struct{}
	fields   []fields.FieldInterface
	fieldMap map[string]int
}

// Render translates a FieldSetType into HTML code and returns it as a template.HTML object.
func (f *FieldSetType) Render(appendData ...map[string]interface{}) template.HTML {
	var s string
	buf := bytes.NewBufferString(s)
	data := map[string]interface{}{
		"container": "fieldset",
		"name":      f.name,
		"fields":    f.fields,
		"classes":   f.class,
		"tags":      f.tags,
	}
	for _, val := range appendData {
		for k, v := range val {
			data[k] = v
		}
	}
	data["append"] = map[string]interface{}{"container": "fieldset"}
	err := template.Must(template.ParseFiles(formcommon.CreateUrl(formcommon.TmplDir+"/"+f.tmpl+".html"))).Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(buf.String())
}

func (f *FieldSetType) SetTmpl(tmpl string) *FieldSetType {
	f.tmpl = tmpl
	return f
}

// FieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func FieldSet(name string, elems ...fields.FieldInterface) *FieldSetType {
	ret := &FieldSetType{
		"fieldset",
		name,
		map[string]struct{}{},
		map[string]struct{}{},
		elems,
		map[string]int{},
	}
	for i, elem := range elems {
		ret.fieldMap[elem.Name()] = i
	}
	return ret
}

//SortAll("field1,field2") or SortAll("field1","field2")
func (f *FieldSetType) SortAll(sortList ...string) *FieldSetType {
	elem := f.fields
	size := len(elem)
	f.fields = make([]fields.FieldInterface, size)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	for k, fieldName := range sortSlice {
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			f.fields[k] = elem[oldIndex]
			f.fieldMap[fieldName] = k
		}
	}
	return f
}

// Elements adds the provided elements to the fieldset.
func (f *FieldSetType) Elements(elems ...FormElement) *FieldSetType {
	for _, e := range elems {
		t := reflect.TypeOf(e)
		switch {
		case t.Implements(reflect.TypeOf((*fields.FieldInterface)(nil)).Elem()):
			f.addField(e.(fields.FieldInterface))
		}
	}
	return f
}

func (f *FieldSetType) addField(field fields.FieldInterface) *FieldSetType {
	f.fields = append(f.fields, field)
	f.fieldMap[field.Name()] = len(f.fields) - 1
	return f
}

//Sort("field1:1,field2:2") or Sort("field1:1","field2:2")
func (f *FieldSetType) Sort(sortList ...string) *FieldSetType {
	size := len(f.fields)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	var index int
	for _, nameIndex := range sortSlice {
		ni := strings.Split(nameIndex, ":")
		fieldName := ni[0]
		if len(ni) > 1 {
			if idx, err := strconv.Atoi(ni[1]); err != nil {
				continue
			} else {
				index = idx
			}
		}
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			if oldIndex != index && size > index {
				f.fields[oldIndex], f.fields[index] = f.fields[index], f.fields[oldIndex]
				f.fieldMap[f.fields[index].Name()] = index
				f.fieldMap[f.fields[oldIndex].Name()] = oldIndex
			}
		}
		index++
	}
	return f
}

// Field returns the field identified by name. It returns an empty field if it is missing.
func (f *FieldSetType) Field(name string) fields.FieldInterface {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &fields.Field{}
	}
	return f.fields[ind].(fields.FieldInterface)
}

// Name returns the name of the fieldset.
func (f *FieldSetType) Name() string {
	return f.name
}

// AddClass saves the provided class for the fieldset.
func (f *FieldSetType) AddClass(class string) *FieldSetType {
	f.class[class] = struct{}{}
	return f
}

// RemoveClass removes the provided class from the fieldset, if it was present. Nothing is done if it was not originally present.
func (f *FieldSetType) RemoveClass(class string) *FieldSetType {
	delete(f.class, class)
	return f
}

// AddTag adds a no-value parameter (e.g.: "disabled", "checked") to the fieldset.
func (f *FieldSetType) AddTag(tag string) *FieldSetType {
	f.tags[tag] = struct{}{}
	return f
}

// RemoveTag removes a tag from the fieldset, if it was present.
func (f *FieldSetType) RemoveTag(tag string) *FieldSetType {
	delete(f.tags, tag)
	return f
}

// Disable adds tag "disabled" to the fieldset, making it unresponsive in some environment (e.g.: Bootstrap).
func (f *FieldSetType) Disable() *FieldSetType {
	f.AddTag("disabled")
	return f
}

// Enable removes tag "disabled" from the fieldset, making it responsive.
func (f *FieldSetType) Enable() *FieldSetType {
	f.RemoveTag("disabled")
	return f
}
