// This package provides form creation and rendering functionalities, as well as FieldSet definition.
// Two kind of forms can be created: base forms and Bootstrap3 compatible forms; even though the latters are automatically provided
// the required classes to make them render correctly in a Bootstrap environment, every form can be given custom parameters such as
// classes, id, generic parameters (in key-value form) and stylesheet options.
package forms

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/fields"
)

// Form methods: POST or GET.
const (
	POST = "POST"
	GET  = "GET"
)

// Form structure.
type Form struct {
	fields       []FormElement
	fieldMap     map[string]int
	containerMap map[string]string
	style        string
	template     *template.Template
	class        []string
	id           string
	params       map[string]string
	css          map[string]string
	method       string
	action       template.HTML
}

func NewForm(style string, args ...string) *Form {
	if style == "" {
		style = formcommon.BASE
	}
	var method, action string
	var tmplFile string = formcommon.TmplDir + "/baseform.html"
	switch len(args) {
	case 0:
		tmplFile = formcommon.TmplDir + "/allfields.html"
	case 1:
		method = args[0]
	case 2:
		method = args[0]
		action = args[1]
	case 3:
		method = args[0]
		action = args[1]
		tmplFile = args[2]
	}
	tmpl, err := template.ParseFiles(formcommon.CreateUrl(tmplFile))
	if err != nil {
		panic(err)
	}
	return &Form{
		make([]FormElement, 0),
		make(map[string]int),
		make(map[string]string),
		style,
		tmpl,
		[]string{},
		"",
		map[string]string{},
		map[string]string{},
		method,
		template.HTML(action),
	}
}

// NewFormFromModel returns a base form inferring fields, data types and contents from the provided instance.
// A Submit button is automatically added as a last field; the form is editable and fields can be added, modified or removed as needed.
// Tags can be used to drive automatic creation: change default widgets for each field, skip fields or provide additional parameters.
// Basic field -> widget mapping is as follows: string -> textField, bool -> checkbox, time.Time -> datetimeField, int -> numberField;
// nested structs are also converted and added to the form.
func NewFormFromModel(m interface{}, style string, args ...string) *Form {
	form := NewForm(style, args...)
	flist, fsort := unWindStructure(m, "")
	for _, v := range flist {
		form.Elements(v.(FormElement))
	}
	form.Elements(FieldSet(
		"_button_group",
		fields.SubmitButton("submit", formcommon.LabelFn("Submit")),
		fields.ResetButton("reset", formcommon.LabelFn("Reset")),
	).SetTmpl("fieldset_buttons"))
	if fsort != "" {
		form.Sort(fsort)
	}
	return form
}

func unWindStructure(m interface{}, baseName string) ([]interface{}, string) {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	fieldList := make([]interface{}, 0)
	fieldSort := ""
	fieldSetList := make(map[string]*FieldSetType, 0)
	fieldSetSort := make(map[string]string, 0)
	for i := 0; i < t.NumField(); i++ {
		optionsArr := strings.Split(formcommon.Tag(t, i, "form_options"), ",")
		options := make(map[string]struct{})
		for _, opt := range optionsArr {
			if opt != "" {
				options[opt] = struct{}{}
			}
		}
		if _, ok := options["-"]; !ok {
			widget := formcommon.Tag(t, i, "form_widget")
			var f fields.FieldInterface
			var fName string
			if baseName == "" {
				fName = t.Field(i).Name
			} else {
				fName = strings.Join([]string{baseName, t.Field(i).Name}, ".")
			}
			//fmt.Println(fName, t.Field(i).Type.String(), t.Field(i).Type.Kind())
			switch widget {
			case "text":
				f = fields.TextFieldFromInstance(m, i, fName)
			case "hidden":
				f = fields.HiddenFieldFromInstance(m, i, fName)
			case "textarea":
				f = fields.TextAreaFieldFromInstance(m, i, fName)
			case "password":
				f = fields.PasswordFieldFromInstance(m, i, fName)
			case "select":
				f = fields.SelectFieldFromInstance(m, i, fName, options)
			case "date":
				f = fields.DateFieldFromInstance(m, i, fName)
			case "datetime":
				f = fields.DatetimeFieldFromInstance(m, i, fName)
			case "time":
				f = fields.TimeFieldFromInstance(m, i, fName)
			case "number":
				f = fields.NumberFieldFromInstance(m, i, fName)
			case "range":
				f = fields.RangeFieldFromInstance(m, i, fName)
			case "radio":
				f = fields.RadioFieldFromInstance(m, i, fName)
			case "checkbox":
				f = fields.CheckboxFieldFromInstance(m, i, fName)
			case "static":
				f = fields.StaticFieldFromInstance(m, i, fName)
			default:
				switch t.Field(i).Type.String() {
				case "string":
					f = fields.TextFieldFromInstance(m, i, fName)
				case "bool":
					f = fields.CheckboxFieldFromInstance(m, i, fName, options)
				case "time.Time":
					f = fields.DatetimeFieldFromInstance(m, i, fName)
				case "int", "int64":
					f = fields.NumberFieldFromInstance(m, i, fName)
				case "float", "float64":
					f = fields.NumberFieldFromInstance(m, i, fName)
				case "struct":
					fl, fs := unWindStructure(v.Field(i).Interface(), fName)
					if fs != "" {
						if fieldSort == "" {
							fieldSort = fs
						} else {
							fieldSort += "," + fs
						}
					}
					fieldList = append(fieldList, fl...)
					f = nil
				default:
					if t.Field(i).Type.Kind() == reflect.Struct ||
						(t.Field(i).Type.Kind() == reflect.Ptr && t.Field(i).Type.Elem().Kind() == reflect.Struct) {
						fl, fs := unWindStructure(v.Field(i).Interface(), fName)
						if fs != "" {
							if fieldSort == "" {
								fieldSort = fs
							} else {
								fieldSort += "," + fs
							}
						}
						fieldList = append(fieldList, fl...)
						f = nil
					} else {
						f = fields.TextFieldFromInstance(m, i, fName)
					}
				}
			}
			if f != nil {
				label := formcommon.Tag(t, i, "form_label")
				if label == "" {
					label = strings.Title(t.Field(i).Name)
				}
				label = formcommon.LabelFn(label)
				f.SetLabel(label)

				params := formcommon.Tag(t, i, "form_params")
				if params != "" {
					if paramsMap, err := url.ParseQuery(params); err == nil {
						for k, v := range paramsMap {
							if k == "placeholder" || k == "title" {
								v[0] = formcommon.LabelFn(v[0])
							}
							f.SetParam(k, v[0])
						}
					} else {
						fmt.Println(err)
					}
				}
				fieldset := formcommon.Tag(t, i, "form_fieldset")
				fieldsort := formcommon.Tag(t, i, "form_sort")
				if fieldset != "" {
					fieldset = formcommon.LabelFn(fieldset)
					if _, ok := fieldSetList[fieldset]; !ok {
						fieldSetList[fieldset] = FieldSet(fieldset, f)
					} else {
						fieldSetList[fieldset].Elements(f)
					}
					if fieldsort != "" {
						if _, ok := fieldSetSort[fieldset]; !ok {
							fieldSetSort[fieldset] = fName + ":" + fieldsort
						} else {
							fieldSetSort[fieldset] += "," + fName + ":" + fieldsort
						}
					}
				} else {
					fieldList = append(fieldList, f)
					if fieldsort != "" {
						if fieldSort == "" {
							fieldSort = fName + ":" + fieldsort
						} else {
							fieldSort += "," + fName + ":" + fieldsort
						}
					}
				}
			}
		}
	}
	for _, v := range fieldSetList {
		if s, ok := fieldSetSort[v.Name()]; ok {
			v.Sort(s)
		}
		fieldList = append(fieldList, v)
	}
	return fieldList, fieldSort
}
