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
	tmpl       string
	name       string
	class      map[string]struct{}
	tags       map[string]struct{}
	fields     []fields.FieldInterface
	fieldMap   map[string]int
	AppendData map[string]interface{}
}

func (f *FieldSetType) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

func (f *FieldSetType) Data() map[string]interface{} {
	data := map[string]interface{}{
		"container": "fieldset",
		"name":      f.name,
		"fields":    f.fields,
		"classes":   f.class,
		"tags":      f.tags,
	}
	for k, v := range f.AppendData {
		data[k] = v
	}
	return data
}

func (f *FieldSetType) dataForRender() string {
	buf := bytes.NewBufferString("")
	tpf := formcommon.TmplDir + "/" + f.tmpl + ".html"
	tpl, ok := formcommon.CachedTemplate(tpf)
	if !ok {
		tpl = template.Must(template.ParseFiles(formcommon.CreateUrl(tpf)))
		formcommon.SetCachedTemplate(tpf, tpl)
	}
	err := tpl.Execute(buf, f.Data())
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Render translates a FieldSetType into HTML code and returns it as a template.HTML object.
func (f *FieldSetType) Render() template.HTML {
	return template.HTML(f.dataForRender())
}

func (f *FieldSetType) String() string {
	return f.dataForRender()
}

func (f *FieldSetType) SetTmpl(tmpl string) *FieldSetType {
	f.tmpl = tmpl
	return f
}

// FieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func FieldSet(name string, elems ...fields.FieldInterface) *FieldSetType {
	ret := &FieldSetType{
		tmpl:       "fieldset",
		name:       name,
		class:      map[string]struct{}{},
		tags:       map[string]struct{}{},
		fields:     elems,
		fieldMap:   map[string]int{},
		AppendData: map[string]interface{}{},
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
	var endIdx int = size - 1
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
			if ni[1] == "last" {
				index = endIdx
			} else if idx, err := strconv.Atoi(ni[1]); err != nil {
				continue
			} else {
				if idx >= 0 {
					index = idx
				} else {
					index = endIdx + idx
				}

			}
		}
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			if oldIndex != index && size > index {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index++
	}
	return f
}

func (f *FieldSetType) Sort2Last(fieldsName ...string) *FieldSetType {
	size := len(f.fields)
	var endIdx int = size - 1
	var index int = endIdx
	for n := len(fieldsName) - 1; n >= 0; n-- {
		fieldName := fieldsName[n]
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			if oldIndex != index && index >= 0 {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index--
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

func (f *FieldSetType) sortFields(index, oldIndex, endIdx, size int) {

	var newFields []fields.FieldInterface = make([]fields.FieldInterface, 0)
	var oldFields []fields.FieldInterface = make([]fields.FieldInterface, size)
	copy(oldFields, f.fields)
	var min, max int
	if index > oldIndex {
		//[ ][I][ ][ ][ ][ ] I:oldIndex=1
		//[ ][ ][ ][ ][I][ ] I:index=4
		if oldIndex > 0 {
			newFields = oldFields[0:oldIndex]
		}
		newFields = append(newFields, oldFields[oldIndex+1:index+1]...)
		newFields = append(newFields, f.fields[oldIndex])
		if index+1 <= endIdx {
			newFields = append(newFields, f.fields[index+1:]...)
		}
		min = oldIndex
		max = index
	} else {
		//[ ][ ][ ][ ][I][ ] I:oldIndex=4
		//[ ][I][ ][ ][ ][ ] I:index=1
		if index > 0 {
			newFields = oldFields[0:index]
		}
		newFields = append(newFields, oldFields[oldIndex])
		newFields = append(newFields, f.fields[index:oldIndex]...)
		if oldIndex+1 <= endIdx {
			newFields = append(newFields, f.fields[oldIndex+1:]...)
		}
		min = index
		max = oldIndex
	}
	for i := min; i <= max; i++ {
		f.fieldMap[newFields[i].Name()] = i
	}
	f.fields = newFields
}
