// This package provides all the input fields logic and customization methods.
package fields

import (
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/widgets"
	"html/template"
	"strings"
)

// Field is a generic type containing all data associated to an input field.
type Field struct {
	fieldType      string
	Widget         widgets.WidgetInterface // Public Widget field for widget customization
	name           string
	class          []string
	id             string
	params         map[string]string
	css            map[string]string
	label          string
	labelClass     []string
	tag            map[string]struct{}
	value          string
	helptext       string
	errors         []string
	additionalData map[string]interface{}
	choices        interface{}
	choiceKeys     map[string]ChoiceIndex
}

// FieldInterface defines the interface an object must implement to be used in a form. Every method returns a FieldInterface object
// to allow methods chaining.
type FieldInterface interface {
	Name() string
	Render(appendData ...map[string]interface{}) template.HTML
	AddClass(class string) FieldInterface
	RemoveClass(class string) FieldInterface
	AddTag(class string) FieldInterface
	RemoveTag(class string) FieldInterface
	SetId(id string) FieldInterface
	SetParam(key, value string) FieldInterface
	DeleteParam(key string) FieldInterface
	AddCss(key, value string) FieldInterface
	RemoveCss(key string) FieldInterface
	SetStyle(style string) FieldInterface
	SetLabel(label string) FieldInterface
	AddLabelClass(class string) FieldInterface
	RemoveLabelClass(class string) FieldInterface
	SetValue(value string) FieldInterface
	Disabled() FieldInterface
	Enabled() FieldInterface
	SetHelptext(text string) FieldInterface
	AddError(err string) FieldInterface
	MultipleChoice() FieldInterface
	SingleChoice() FieldInterface
	AddSelected(opt ...string) FieldInterface
	RemoveSelected(opt string) FieldInterface
	SetChoices(choices interface{}) FieldInterface
	SetText(text string) FieldInterface
}

// FieldWithType creates an empty field of the given type and identified by name.
func FieldWithType(name, t string) *Field {
	return &Field{
		fieldType:      t,
		Widget:         nil,
		name:           name,
		class:          []string{},
		id:             "",
		params:         map[string]string{},
		css:            map[string]string{},
		label:          "",
		labelClass:     []string{},
		tag:            map[string]struct{}{},
		value:          "",
		helptext:       "",
		errors:         []string{},
		additionalData: map[string]interface{}{},
		choices:        nil,
		choiceKeys:     map[string]ChoiceIndex{},
	}
}

// SetStyle sets the style (e.g.: BASE, BOOTSTRAP) of the field, correctly populating the Widget field.
func (f *Field) SetStyle(style string) FieldInterface {
	f.Widget = widgets.BaseWidget(style, f.fieldType)
	return f
}

// Name returns the name of the field.
func (f *Field) Name() string {
	return strings.TrimSuffix(f.name, "[]")
}

func (f *Field) dataForRender(appendData ...map[string]interface{}) map[string]interface{} {
	safeParams := make(map[template.HTMLAttr]string)
	for k, v := range f.params {
		safeParams[template.HTMLAttr(k)] = v
	}
	data := map[string]interface{}{
		"classes":      f.class,
		"id":           f.id,
		"name":         f.name,
		"params":       safeParams,
		"css":          f.css,
		"type":         f.fieldType,
		"label":        f.label,
		"labelClasses": f.labelClass,
		"tags":         f.tag,
		"value":        f.value,
		"helptext":     f.helptext,
		"errors":       f.errors,
		"container":    "form",
		"choices":      f.choices,
	}
	for k, v := range f.additionalData {
		data[k] = v
	}
	for _, val := range appendData {
		for k, v := range val {
			data[k] = v
		}
	}
	return data
}

// Render packs all data and executes widget render method.
func (f *Field) Render(appendData ...map[string]interface{}) template.HTML {
	if f.Widget != nil {
		data := f.dataForRender(appendData...)
		return template.HTML(f.Widget.Render(data))
	}
	return template.HTML("")
}

// AddClass adds a class to the field.
func (f *Field) AddClass(class string) FieldInterface {
	f.class = append(f.class, class)
	return f
}

// RemoveClass removes a class from the field, if it was present.
func (f *Field) RemoveClass(class string) FieldInterface {
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

// SetId associates the given id to the field, overwriting any previous id.
func (f *Field) SetId(id string) FieldInterface {
	f.id = id
	return f
}

// SetLabel saves the label to be rendered along with the field.
func (f *Field) SetLabel(label string) FieldInterface {
	f.label = label
	return f
}

// SetLablClass allows to define custom classes for the label.
func (f *Field) AddLabelClass(class string) FieldInterface {
	f.labelClass = append(f.labelClass, class)
	return f
}

// RemoveLabelClass removes the given class from the field label.
func (f *Field) RemoveLabelClass(class string) FieldInterface {
	ind := -1
	for i, v := range f.labelClass {
		if v == class {
			ind = i
			break
		}
	}

	if ind != -1 {
		f.labelClass = append(f.labelClass[:ind], f.labelClass[ind+1:]...)
	}
	return f
}

// SetParam adds a parameter (defined as key-value pair) in the field.
func (f *Field) SetParam(key, value string) FieldInterface {
	f.params[key] = value
	return f
}

// DeleteParam removes a parameter identified by key from the field.
func (f *Field) DeleteParam(key string) FieldInterface {
	delete(f.params, key)
	return f
}

// AddCss adds a custom CSS style the field.
func (f *Field) AddCss(key, value string) FieldInterface {
	f.css[key] = value
	return f
}

// RemoveCss removes CSS options identified by key from the field.
func (f *Field) RemoveCss(key string) FieldInterface {
	delete(f.css, key)
	return f
}

// Disabled add the "disabled" tag to the field, making it unresponsive in some environments (e.g. Bootstrap).
func (f *Field) Disabled() FieldInterface {
	f.AddTag("disabled")
	return f
}

// Enabled removes the "disabled" tag from the field, making it responsive.
func (f *Field) Enabled() FieldInterface {
	f.RemoveTag("disabled")
	return f
}

// AddTag adds a no-value parameter (e.g.: checked, disabled) to the field.
func (f *Field) AddTag(tag string) FieldInterface {
	f.tag[tag] = struct{}{}
	return f
}

// RemoveTag removes a no-value parameter from the field.
func (f *Field) RemoveTag(tag string) FieldInterface {
	delete(f.tag, tag)
	return f
}

// SetValue sets the value parameter for the field.
func (f *Field) SetValue(value string) FieldInterface {
	f.value = value
	return f
}

// SetHelptext saves the field helptext.
func (f *Field) SetHelptext(text string) FieldInterface {
	f.helptext = text
	return f
}

// AddError adds an error string to the field. It's valid only for Bootstrap forms.
func (f *Field) AddError(err string) FieldInterface {
	f.errors = append(f.errors, err)
	return f
}

// MultipleChoice configures the SelectField to accept and display multiple choices.
// It has no effect if type is not SELECT.
func (f *Field) MultipleChoice() FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		f.AddTag("multiple")
		fallthrough
	case formcommon.CHECKBOX:
		// fix name if necessary
		if !strings.HasSuffix(f.name, "[]") {
			f.name = f.name + "[]"
		}
	}
	return f
}

// SingleChoice configures the Field to accept and display only one choice (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) SingleChoice() FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		f.RemoveTag("multiple")
		fallthrough
	case formcommon.CHECKBOX:
		if strings.HasSuffix(f.name, "[]") {
			f.name = strings.TrimSuffix(f.name, "[]")
		}
	}
	return f
}

// If the field is configured as "multiple", AddSelected adds a selected value to the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) AddSelected(opt ...string) FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		for _, v := range opt {
			i := f.choiceKeys[v]
			if vc, ok := f.choices.(map[string][]InputChoice)[i.Group]; ok {
				if len(vc)>i.Index {
					f.choices.(map[string][]InputChoice)[i.Group][i.Index].Checked = true
				}
			}
		}
	case formcommon.CHECKBOX:
		size := len(f.choices.([]InputChoice))
		for _, v := range opt {
			i := f.choiceKeys[v]
			if size > i.Index {
				f.choices.([]InputChoice)[i.Index].Checked = true
			}
		}
	}
	return f
}

// If the field is configured as "multiple", AddSelected removes the selected value from the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) RemoveSelected(opt string) FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		i := f.choiceKeys[opt]
		if vc, ok := f.choices.(map[string][]InputChoice)[i.Group]; ok {
			if len(vc)>i.Index {
				f.choices.(map[string][]InputChoice)[i.Group][i.Index].Checked = false
			}
		}

	case formcommon.CHECKBOX:
		size := len(f.choices.([]InputChoice))
		i := f.choiceKeys[opt]
		if size > i.Index {
			f.choices.([]InputChoice)[i.Index].Checked = false
		}
	}
	return f
}

// SetChoices takes as input a dictionary whose key-value entries are defined as follows: key is the group name (the empty string
// is the default group that is not explicitly rendered) and value is the list of choices belonging to that group.
// Grouping is only useful for Select fields, while groups are ignored in Radio fields.
// It has no effect if type is not SELECT.
func (f *Field) SetChoices(choices interface{}) FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		f.choices, _ = choices.(map[string][]InputChoice)

	case formcommon.RADIO, formcommon.CHECKBOX:
		f.choices, _ = choices.([]InputChoice)
	}
	return f
}

// SetText saves the provided text as content of the field, usually a TextAreaField.
func (f *Field) SetText(text string) FieldInterface {
	if f.fieldType == formcommon.BUTTON ||
		f.fieldType == formcommon.SUBMIT ||
		f.fieldType == formcommon.RESET ||
		f.fieldType == formcommon.STATIC ||
		f.fieldType == formcommon.TEXTAREA {
		f.additionalData["text"] = text
	}
	return f
}
