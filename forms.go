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
	"github.com/coscms/xweb/validation"
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
	AppendData   map[string]interface{}
	valid        *validation.Validation
	model        interface{}
}

func (f *Form) Valid(args ...string) (valid *validation.Validation, passed bool) {
	if f.valid == nil {
		f.valid = &validation.Validation{}
	}
	valid = f.valid
	if f.model == nil {
		return
	}
	var err error
	passed, err = valid.Valid(f.model, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !passed { // validation does not pass
		for field, err := range valid.ErrorsMap {
			f.Field(field).AddError(formcommon.LabelFn(err.Message))
		}
	}
	return
}

func (f *Form) SetModel(m interface{}) *Form {
	f.model = m
	return f
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
	tmpl, ok := formcommon.CachedTemplate(tmplFile)
	if !ok {
		tmpl = template.Must(template.ParseFiles(formcommon.CreateUrl(tmplFile)))
		formcommon.SetCachedTemplate(tmplFile, tmpl)
	}
	return &Form{
		fields:       make([]FormElement, 0),
		fieldMap:     make(map[string]int),
		containerMap: make(map[string]string),
		style:        style,
		template:     tmpl,
		class:        []string{},
		id:           "",
		params:       map[string]string{},
		css:          map[string]string{},
		method:       method,
		action:       template.HTML(action),
		AppendData:   map[string]interface{}{},
	}
}

// NewFormFromModel returns a base form inferring fields, data types and contents from the provided instance.
// A Submit button is automatically added as a last field; the form is editable and fields can be added, modified or removed as needed.
// Tags can be used to drive automatic creation: change default widgets for each field, skip fields or provide additional parameters.
// Basic field -> widget mapping is as follows: string -> textField, bool -> checkbox, time.Time -> datetimeField, int -> numberField;
// nested structs are also converted and added to the form.
func NewFormFromModel(m interface{}, style string, args ...string) *Form {
	form := NewForm(style, args...)
	form.SetModel(m)
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
				f = fields.TextFieldFromInstance(v, t, i, fName)
			case "hidden":
				f = fields.HiddenFieldFromInstance(v, t, i, fName)
			case "textarea":
				f = fields.TextAreaFieldFromInstance(v, t, i, fName)
			case "password":
				f = fields.PasswordFieldFromInstance(v, t, i, fName)
			case "select":
				f = fields.SelectFieldFromInstance(v, t, i, fName, options)
			case "date":
				f = fields.DateFieldFromInstance(v, t, i, fName)
			case "datetime":
				f = fields.DatetimeFieldFromInstance(v, t, i, fName)
			case "time":
				f = fields.TimeFieldFromInstance(v, t, i, fName)
			case "number":
				f = fields.NumberFieldFromInstance(v, t, i, fName)
			case "range":
				f = fields.RangeFieldFromInstance(v, t, i, fName)
			case "radio":
				f = fields.RadioFieldFromInstance(v, t, i, fName)
			case "checkbox":
				f = fields.CheckboxFieldFromInstance(v, t, i, fName)
			case "static":
				f = fields.StaticFieldFromInstance(v, t, i, fName)
			default:
				switch t.Field(i).Type.String() {
				case "string":
					f = fields.TextFieldFromInstance(v, t, i, fName)
				case "bool":
					f = fields.CheckboxFromInstance(v, t, i, fName, options)
				case "time.Time":
					f = fields.DatetimeFieldFromInstance(v, t, i, fName)
				case "int", "int64":
					f = fields.NumberFieldFromInstance(v, t, i, fName)
				case "float", "float64":
					f = fields.NumberFieldFromInstance(v, t, i, fName)
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
						f = fields.TextFieldFromInstance(v, t, i, fName)
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
				valid := formcommon.Tag(t, i, "valid")
				if valid != "" {
					ValidTagFn(valid, f)
				}
				fieldset := formcommon.Tag(t, i, "form_fieldset")
				fieldsort := formcommon.Tag(t, i, "form_sort")
				if fieldset != "" {
					fieldset = formcommon.LabelFn(fieldset)
					f.SetData("container", "fieldset")
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

var ValidTagFn func(string, fields.FieldInterface) = Html5Validate

func ValidationEngine(valid string, f fields.FieldInterface) {
	//for jQuery-Validation-Engine
	validFuncs := strings.Split(valid, ";")
	var validClass string
	for _, v := range validFuncs {
		pos := strings.Index(v, "(")
		var fn string
		if pos > -1 {
			fn = v[0:pos]
		} else {
			fn = v
		}
		switch fn {
		case "Required":
			validClass += "," + strings.ToLower(fn)
		case "Min", "Max":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += "," + strings.ToLower(fn) + "[" + val + "]"
		case "Range":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			rangeVals := strings.SplitN(val, ",", 2)
			validClass += ",min[" + strings.TrimSpace(rangeVals[0]) + "],max[" + strings.TrimSpace(rangeVals[1]) + "]"
		case "MinSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += ",minSize[" + val + "]"
		case "MaxSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += ",maxSize[" + val + "]"
		case "Numeric":
			validClass += ",number"
		case "AlphaNumeric":
			validClass += ",custom[onlyLetterNumber]"
		/*
			case "Length":
				validClass += ",length"
			case "Match":
				val := v[pos+1:]
				val = strings.TrimSuffix(val, ")")
				val = strings.Trim(val, "/")
				validClass += ",match[]"
		*/
		case "AlphaDash":
			validClass += ",custom[onlyLetterNumber]"
		case "IP":
			validClass += ",custom[ipv4]"
		case "Alpha", "Email", "Base64", "Mobile", "Tel", "Phone":
			validClass += ",custom[" + strings.ToLower(fn) + "]"
		case "ZipCode":
			validClass += ",custom[zip]"
		}
	}
	if validClass != "" {
		validClass = strings.TrimPrefix(validClass, ",")
		validClass = "validate[" + validClass + "]"
		f.AddClass(validClass)
	}
}

func Html5Validate(valid string, f fields.FieldInterface) {
	validFuncs := strings.Split(valid, ";")
	for _, v := range validFuncs {
		pos := strings.Index(v, "(")
		var fn string
		if pos > -1 {
			fn = v[0:pos]
		} else {
			fn = v
		}
		switch fn {
		case "Required":
			f.AddTag(strings.ToLower(fn))
		case "Min", "Max":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam(strings.ToLower(fn), val)
		case "Range":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			rangeVals := strings.SplitN(val, ",", 2)
			f.SetParam("min", strings.TrimSpace(rangeVals[0]))
			f.SetParam("max", strings.TrimSpace(rangeVals[1]))
		case "MinSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("data-min", val)
		case "MaxSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("maxlength", val)
			f.SetParam("data-max", val)
		case "Numeric":
			f.SetParam("pattern", template.HTML("^\\-?\\d+(\\.\\d+)?$"))
		case "AlphaNumeric":
			f.SetParam("pattern", template.HTML("^[\\w\\d]+$"))
		case "Length":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("pattern", ".{"+val+"}")
		case "Match":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			val = strings.Trim(val, "/")
			f.SetParam("pattern", template.HTML(val))

		case "AlphaDash":
			f.SetParam("pattern", template.HTML("^[\\d\\w-]+$"))
		case "IP":
			f.SetParam("pattern", template.HTML("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$"))
		case "Alpha":
			f.SetParam("pattern", template.HTML("^[a-zA-Z]+$"))
		case "Email":
			f.SetParam("pattern", template.HTML("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?"))
		case "Base64":
			f.SetParam("pattern", template.HTML("^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$"))
		case "Mobile":
			f.SetParam("pattern", template.HTML("^((\\+86)|(86))?(1(([35][0-9])|(47)|[8][\\d]))\\d{8}$"))
		case "Tel":
			f.SetParam("pattern", template.HTML("^(0\\d{2,3}(\\-)?)?\\d{7,8}$"))
		case "Phone":
			f.SetParam("pattern", template.HTML("^(((\\+86)|(86))?(1(([35][0-9])|(47)|[8][012356789]))\\d{8}|(0\\d{2,3}(\\-)?)?\\d{7,8})$"))
		case "ZipCode":
			f.SetParam("pattern", template.HTML("^[1-9]\\d{5}$"))
		}
	}
}
