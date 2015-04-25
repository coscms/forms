package forms

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"

	"github.com/coscms/forms/fields"
)

// FormElement interface defines a form object (usually a Field or a FieldSet) that can be rendered as a template.HTML object.
type FormElement interface {
	Render() template.HTML
	Name() string
	String() string
	SetData(key string, value interface{})
	Data() map[string]interface{}
}

func (f *Form) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

func (f *Form) Data() map[string]interface{} {
	data := map[string]interface{}{
		"container": "",
		"fields":    f.fields,
		"classes":   f.class,
		"id":        f.id,
		"params":    f.params,
		"css":       f.css,
		"method":    f.method,
		"action":    f.action,
	}
	for k, v := range f.AppendData {
		data[k] = v
	}
	return data
}

func (f *Form) dataForRender() string {
	buf := bytes.NewBufferString("")
	err := f.template.Execute(buf, f.Data())
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Render executes the internal template and renders the form, returning the result as a template.HTML object embeddable
// in any other template.
func (f *Form) Render() template.HTML {
	return template.HTML(f.dataForRender())
}

func (f *Form) Html(value interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("%v", value))
}

func (f *Form) String() string {
	return f.dataForRender()
}

// Elements adds the provided elements to the form.
func (f *Form) Elements(elems ...FormElement) *Form {
	for _, e := range elems {
		t := reflect.TypeOf(e)
		switch {
		case t.Implements(reflect.TypeOf((*fields.FieldInterface)(nil)).Elem()):
			f.addField(e.(fields.FieldInterface))
		case reflect.ValueOf(e).Type().String() == "*forms.FieldSetType":
			f.addFieldSet(e.(*FieldSetType))
		}
	}
	return f
}

func (f *Form) addField(field fields.FieldInterface) *Form {
	field.SetStyle(f.style)
	f.fields = append(f.fields, field)
	f.fieldMap[field.Name()] = len(f.fields) - 1
	return f
}

func (f *Form) addFieldSet(fs *FieldSetType) *Form {
	for _, v := range fs.fields {
		v.SetStyle(f.style)
		v.SetData("container", "fieldset")
		f.containerMap[v.Name()] = fs.name
	}
	f.fields = append(f.fields, fs)
	f.fieldMap[fs.Name()] = len(f.fields) - 1
	return f
}

// RemoveElement removes an element (identified by name) from the Form.
func (f *Form) RemoveElement(name string) *Form {
	ind, ok := f.fieldMap[name]
	if !ok {
		return f
	}
	delete(f.fieldMap, name)
	f.fields = append(f.fields[:ind], f.fields[ind+1:]...)
	return f
}

// AddClass associates the provided class to the Form.
func (f *Form) AddClass(class string) *Form {
	f.class = append(f.class, class)
	return f
}

// RemoveClass removes the given class (if present) from the Form.
func (f *Form) RemoveClass(class string) *Form {
	ind := -1
	for i, v := range f.class {
		if v == class {
			ind = i
			break
		}
	}

	if ind != -1 {
		f.class = append(f.class[:ind], f.class[ind+1:]...)
	}
	return f
}

// SetId set the given id to the form.
func (f *Form) SetId(id string) *Form {
	f.id = id
	return f
}

// SetParam adds the given key-value pair to form parameters list.
func (f *Form) SetParam(key, value string) *Form {
	f.params[key] = value
	return f
}

// DeleteParm removes the parameter identified by key from form parameters list.
func (f *Form) DeleteParam(key string) *Form {
	delete(f.params, key)
	return f
}

// AddCss add a CSS value (in the form of option-value - e.g.: border - auto) to the form.
func (f *Form) AddCss(key, value string) *Form {
	f.css[key] = value
	return f
}

// RemoveCss removes CSS style from the form.
func (f *Form) RemoveCss(key string) *Form {
	delete(f.css, key)
	return f
}

// Field returns the field identified by name. It returns an empty field if it is missing.
func (f *Form) Field(name string) fields.FieldInterface {
	ind, ok := f.fieldMap[name]
	if !ok || !reflect.TypeOf(f.fields[ind]).Implements(reflect.TypeOf((*fields.FieldInterface)(nil)).Elem()) {
		if v, ok2 := f.containerMap[name]; ok2 {
			return f.FieldSet(v).Field(name)
		}
		return &fields.Field{}
	}
	return f.fields[ind].(fields.FieldInterface)
}

// Fields returns all field
func (f *Form) Fields() []FormElement {
	return f.fields
}

// Field returns the field identified by name. It returns an empty field if it is missing.
func (f *Form) FieldSet(name string) *FieldSetType {
	ind, ok := f.fieldMap[name]
	if !ok || reflect.ValueOf(f.fields[ind]).Type().String() != "*forms.FieldSetType" {
		return &FieldSetType{}
	}
	return f.fields[ind].(*FieldSetType)
}

// FieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func (f *Form) NewFieldSet(name string, elems ...fields.FieldInterface) *FieldSetType {
	return FieldSet(name, elems...)
}

//SortAll("field1,field2") or SortAll("field1","field2")
func (f *Form) SortAll(sortList ...string) *Form {
	elem := f.fields
	size := len(elem)
	f.fields = make([]FormElement, size)
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

//Sort("field1:1,field2:2") or Sort("field1:1","field2:2")
func (f *Form) Sort(sortList ...string) *Form {
	size := len(f.fields)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	var index int
	var endIdx int = size - 1

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
					index = size + idx
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

func (f *Form) Sort2Last(fieldsName ...string) *Form {
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

func (f *Form) sortFields(index, oldIndex, endIdx, size int) {

	var newFields []FormElement = make([]FormElement, 0)
	var oldFields []FormElement = make([]FormElement, size)
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
