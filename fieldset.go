/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package forms

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
)

// FieldSetType is a collection of fields grouped within a form.
type FieldSetType struct {
	OrigName     string                 `json:"origName" xml:"origName"`
	CurrName     string                 `json:"currName" xml:"currName"`
	Label        string                 `json:"label" xml:"label"`
	LabelCols    int                    `json:"labelCols" xml:"labelCols"`
	FieldCols    int                    `json:"fieldCols" xml:"fieldCols"`
	Classes      common.HTMLAttrValues  `json:"classes" xml:"classes"`
	Tags         common.HTMLAttrValues  `json:"tags" xml:"tags"`
	Helptext     string                 `json:"helpText" xml:"helpText"`
	FieldList    []config.FormElement   `json:"fieldList" xml:"fieldList"`
	AppendData   map[string]interface{} `json:"appendData,omitempty" xml:"appendData,omitempty"`
	FormTheme    string                 `json:"formTheme" xml:"formTheme"`
	Language     string                 `json:"language,omitempty" xml:"language,omitempty"`
	Template     string                 `json:"template" xml:"template"`
	fieldMap     map[string]int
	containerMap map[string]string
	data         map[string]interface{}
}

func (f *FieldSetType) Cols() int {
	return config.GetCols(f.LabelCols, f.FieldCols)
}

func (f *FieldSetType) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

// SetHelptext saves the field helptext.
func (f *FieldSetType) SetHelptext(text string) *FieldSetType {
	f.Helptext = text
	return f
}

func (f *FieldSetType) SetLabelCols(cols int) {
	f.LabelCols = cols
}

func (f *FieldSetType) SetFieldCols(cols int) {
	f.FieldCols = cols
}

func (f *FieldSetType) SetName(name string) {
	f.CurrName = name
}

func (f *FieldSetType) OriginalName() string {
	return f.OrigName
}

func (f *FieldSetType) Data() map[string]interface{} {
	if len(f.data) > 0 {
		return f.data
	}
	f.data = map[string]interface{}{
		"container": "fieldset",
		"name":      f.CurrName,
		"label":     f.Label,
		"labelCols": f.LabelCols,
		"fieldCols": f.FieldCols,
		"fields":    f.FieldList,
		"groups":    config.SplitGroup(f.FieldList),
		"classes":   f.Classes,
		"tags":      f.Tags,
		"helptext":  f.Helptext,
	}
	for k, v := range f.AppendData {
		f.data[k] = v
	}
	return f.data
}

func (f *FieldSetType) render() string {
	buf := bytes.NewBuffer(nil)
	tpf := common.TmplDir(f.FormTheme) + "/" + f.FormTheme + "/" + f.Template + ".html"
	tpl, err := common.GetOrSetCachedTemplate(tpf, func() (*template.Template, error) {
		return common.ParseFiles(common.LookupPath(tpf))
	})
	if err != nil {
		return err.Error()
	}
	err = tpl.Execute(buf, f.Data())
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// Render translates a FieldSetType into HTML code and returns it as a template.HTML object.
func (f *FieldSetType) Render() template.HTML {
	return template.HTML(f.render())
}

func (f *FieldSetType) String() string {
	return f.render()
}

func (f *FieldSetType) SetLang(lang string) {
	f.Language = lang
}

func (f *FieldSetType) Lang() string {
	return f.Language
}

func (f *FieldSetType) Clone() config.FormElement {
	fc := *f
	return &fc
}

func (f *FieldSetType) SetTemplate(tmpl string) *FieldSetType {
	f.Template = tmpl
	return f
}

// FieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func FieldSet(name string, label string, theme string, elems ...config.FormElement) *FieldSetType {
	ret := &FieldSetType{
		Template:     "fieldset",
		CurrName:     name,
		OrigName:     name,
		Label:        label,
		Classes:      common.HTMLAttrValues{},
		Tags:         common.HTMLAttrValues{},
		FieldList:    elems,
		containerMap: make(map[string]string),
		fieldMap:     map[string]int{},
		AppendData:   map[string]interface{}{},
		FormTheme:    theme,
	}
	for i, elem := range elems {
		ret.fieldMap[elem.OriginalName()] = i
	}
	return ret
}

// SortAll("field1,field2") or SortAll("field1","field2")
func (f *FieldSetType) SortAll(sortList ...string) *FieldSetType {
	elem := f.FieldList
	size := len(elem)
	f.FieldList = make([]config.FormElement, size)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	for k, fieldName := range sortSlice {
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			f.FieldList[k] = elem[oldIndex]
			f.fieldMap[fieldName] = k
		}
	}
	return f
}

// Elements adds the provided elements to the fieldset.
func (f *FieldSetType) Elements(elems ...config.FormElement) {
	for _, e := range elems {
		switch v := e.(type) {
		case fields.FieldInterface:
			f.addField(v)
		case *FieldSetType:
			f.addFieldSet(v)
		case *LangSetType:
			f.addLangSet(v)
		}
	}
}

func (f *FieldSetType) addFieldSet(fs *FieldSetType) *FieldSetType {
	for _, v := range fs.FieldList {
		v.SetData("container", "fieldset")
		f.containerMap[v.OriginalName()] = fs.OriginalName()
	}
	f.FieldList = append(f.FieldList, fs)
	f.fieldMap[fs.OriginalName()] = len(f.FieldList) - 1
	return f
}

func (f *FieldSetType) addLangSet(fs *LangSetType) *FieldSetType {
	for _, v := range fs.fieldMap {
		v.SetData("container", "langset")
		f.containerMap[v.OriginalName()] = fs.OriginalName()
	}
	f.FieldList = append(f.FieldList, fs)
	f.fieldMap[fs.OriginalName()] = len(f.FieldList) - 1
	return f
}

func (f *FieldSetType) addField(field fields.FieldInterface) *FieldSetType {
	field.SetTheme(f.FormTheme)
	field.SetData(`container`, `fieldset`)
	f.FieldList = append(f.FieldList, field)
	f.fieldMap[field.OriginalName()] = len(f.FieldList) - 1
	return f
}

// Sort("field1:1,field2:2") or Sort("field1:1","field2:2")
func (f *FieldSetType) Sort(sortList ...string) *FieldSetType {
	size := len(f.FieldList)
	endIdx := size - 1
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
	size := len(f.FieldList)
	endIdx := size - 1
	index := endIdx
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
	switch v := f.FieldList[ind].(type) {
	case fields.FieldInterface:
		return v
	case *FieldSetType:
		if v, ok := f.containerMap[name]; ok {
			return f.FieldSet(v).Field(name)
		}
	case *LangSetType:
		if v, ok := f.containerMap[name]; ok {
			return f.LangSet(v).Field(name)
		}
	}
	return &fields.Field{}
}

// FieldSet returns the fieldset identified by name. It returns an empty field if it is missing.
func (f *FieldSetType) FieldSet(name string) *FieldSetType {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &FieldSetType{}
	}
	switch v := f.FieldList[ind].(type) {
	case *FieldSetType:
		return v
	default:
		return &FieldSetType{}
	}
}

// LangSet returns the fieldset identified by name. It returns an empty field if it is missing.
func (f *FieldSetType) LangSet(name string) *LangSetType {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &LangSetType{}
	}
	switch v := f.FieldList[ind].(type) {
	case *LangSetType:
		return v
	default:
		return &LangSetType{}
	}
}

// Name returns the name of the fieldset.
func (f *FieldSetType) Name() string {
	return f.CurrName
}

// AddClass saves the provided class for the fieldset.
func (f *FieldSetType) AddClass(class string) *FieldSetType {
	f.Classes.Add(class)
	return f
}

// RemoveClass removes the provided class from the fieldset, if it was present. Nothing is done if it was not originally present.
func (f *FieldSetType) RemoveClass(class string) *FieldSetType {
	f.Classes.Remove(class)
	return f
}

// AddTag adds a no-value parameter (e.g.: "disabled", "checked") to the fieldset.
func (f *FieldSetType) AddTag(tag string) *FieldSetType {
	f.Tags.Add(tag)
	return f
}

// RemoveTag removes a tag from the fieldset, if it was present.
func (f *FieldSetType) RemoveTag(tag string) *FieldSetType {
	f.Tags.Remove(tag)
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

	var newFields []config.FormElement
	oldFields := make([]config.FormElement, size)
	copy(oldFields, f.FieldList)
	var min, max int
	if index > oldIndex {
		//[ ][I][ ][ ][ ][ ] I:oldIndex=1
		//[ ][ ][ ][ ][I][ ] I:index=4
		if oldIndex > 0 {
			newFields = oldFields[0:oldIndex]
		}
		newFields = append(newFields, oldFields[oldIndex+1:index+1]...)
		newFields = append(newFields, f.FieldList[oldIndex])
		if index+1 <= endIdx {
			newFields = append(newFields, f.FieldList[index+1:]...)
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
		newFields = append(newFields, f.FieldList[index:oldIndex]...)
		if oldIndex+1 <= endIdx {
			newFields = append(newFields, f.FieldList[oldIndex+1:]...)
		}
		min = index
		max = oldIndex
	}
	for i := min; i <= max; i++ {
		f.fieldMap[newFields[i].OriginalName()] = i
	}
	f.FieldList = newFields
}
