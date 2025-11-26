package config

import (
	"strings"

	"github.com/webx-top/com"
)

type Element struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	Label        string                 `json:"label"`
	LabelCols    int                    `json:"labelCols,omitempty"`
	FieldCols    int                    `json:"fieldCols,omitempty"`
	LabelClasses []string               `json:"labelClasses,omitempty"`
	Value        string                 `json:"value"`
	HelpText     string                 `json:"helpText"`
	Template     string                 `json:"template"`
	Valid        string                 `json:"valid"`
	Attributes   [][]string             `json:"attributes"`
	Choices      []*Choice              `json:"choices"`
	Elements     []*Element             `json:"elements"`
	Format       string                 `json:"format"`
	Languages    []*Language            `json:"languages,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

func (c *Element) Merge(source *Element) *Element {
	if len(c.ID) == 0 && len(source.ID) > 0 {
		c.ID = source.ID
	}
	if len(c.Type) == 0 && len(source.Type) > 0 {
		c.Type = source.Type
	}
	if len(c.Template) == 0 && len(source.Template) > 0 {
		c.Template = source.Template
	}
	if len(c.Label) == 0 && len(source.Label) > 0 {
		c.Label = source.Label
	}
	if len(c.Name) == 0 && len(source.Name) > 0 {
		c.Name = source.Name
	}
	if c.LabelCols == 0 && source.LabelCols > 0 {
		c.LabelCols = source.LabelCols
	}
	if c.FieldCols == 0 && source.FieldCols > 0 {
		c.FieldCols = source.FieldCols
	}
	if len(c.Value) == 0 && len(source.Value) > 0 {
		c.Value = source.Value
	}
	if len(c.HelpText) == 0 && len(source.HelpText) > 0 {
		c.HelpText = source.HelpText
	}
	if len(c.Valid) == 0 && len(source.Valid) > 0 {
		c.Valid = source.Valid
	}
	if len(c.Format) == 0 && len(source.Format) > 0 {
		c.Format = source.Format
	}
	var found bool
	for _, v := range source.Attributes {
		if len(v) == 0 {
			continue
		}
		for _, v2 := range c.Attributes {
			if len(v2) == 0 {
				continue
			}
			if v2[0] == v[0] {
				found = true
				break
			}
		}
		if !found {
			c.Attributes = append(c.Attributes, v)
		} else {
			found = false
		}
	}
	for _, v := range source.LabelClasses {
		for _, v2 := range c.LabelClasses {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			c.LabelClasses = append(c.LabelClasses, v)
		} else {
			found = false
		}
	}
	for _, v := range source.Choices {
		for _, v2 := range c.Choices {
			if v.Group == v2.Group {
				found = true
				v2.Merge(v)
				break
			}
		}
		if !found {
			c.Choices = append(c.Choices, v)
		} else {
			found = false
		}
	}
	for _, v := range source.Elements {
		if len(v.Name) > 0 {
			for _, v2 := range c.Elements {
				if v.Name == v2.Name {
					found = true
					v2.Merge(v)
					break
				}
			}
		}
		if !found {
			c.Elements = append(c.Elements, v)
		} else {
			found = false
		}
	}
	for _, v := range source.Languages {
		if len(v.ID) > 0 {
			for _, v2 := range c.Languages {
				if v.ID == v2.ID {
					found = true
					break
				}
			}
		}
		if !found {
			c.Languages = append(c.Languages, v)
		} else {
			found = false
		}
	}
	if source.Data != nil {
		if c.Data == nil {
			c.Data = map[string]interface{}{}
		}
		for k, v := range source.Data {
			c.Data[k] = v
		}
	}
	return c
}

func (e *Element) Cols() int {
	return GetCols(e.LabelCols, e.FieldCols)
}

func (e *Element) Clone() *Element {
	elements := make([]*Element, len(e.Elements))
	languages := make([]*Language, len(e.Languages))
	choices := make([]*Choice, len(e.Choices))
	for index, elem := range e.Elements {
		elements[index] = elem.Clone()
	}
	for index, value := range e.Languages {
		languages[index] = value.Clone()
	}
	for index, value := range e.Choices {
		choices[index] = value.Clone()
	}
	r := &Element{
		ID:           e.ID,
		Type:         e.Type,
		Name:         e.Name,
		Label:        e.Label,
		LabelCols:    e.LabelCols,
		FieldCols:    e.FieldCols,
		LabelClasses: make([]string, len(e.LabelClasses)),
		Value:        e.Value,
		HelpText:     e.HelpText,
		Template:     e.Template,
		Valid:        e.Valid,
		Attributes:   make([][]string, len(e.Attributes)),
		Choices:      choices,
		Elements:     elements,
		Format:       e.Format,
		Languages:    languages,
		Data:         map[string]interface{}{},
	}
	for k, v := range e.Data {
		r.Data[k] = v
	}
	if len(e.LabelClasses) > 0 {
		copy(r.LabelClasses, e.LabelClasses)
	}
	for k, v := range e.Attributes {
		cv := make([]string, len(v))
		copy(cv, v)
		r.Attributes[k] = cv
	}
	return r
}

func (e *Element) HasAttr(attrs ...string) bool {
	mk := map[string]struct{}{}
	for _, attr := range attrs {
		mk[strings.ToLower(attr)] = struct{}{}
	}
	for _, v := range e.Attributes {
		if len(v) == 0 || len(v[0]) == 0 {
			continue
		}
		v[0] = strings.ToLower(v[0])
		if _, ok := mk[v[0]]; ok {
			return true
		}
	}
	return false
}

func (e *Element) AddElement(elements ...*Element) *Element {
	e.Elements = append(e.Elements, elements...)
	return e
}

func (e *Element) AddLanguage(languages ...*Language) *Element {
	e.Languages = append(e.Languages, languages...)
	return e
}

func (e *Element) AddAttribute(attributes ...string) *Element {
	e.Attributes = append(e.Attributes, attributes)
	return e
}

func (e *Element) AddChoice(choices ...*Choice) *Element {
	e.Choices = append(e.Choices, choices...)
	return e
}

func (e *Element) AddLabelClass(labelClasses ...string) *Element {
	e.LabelClasses = append(e.LabelClasses, labelClasses...)
	return e
}

func (e *Element) Set(name string, value interface{}) *Element {
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	e.Data[name] = value
	return e
}

func (e *Element) GetMultilingualText(recv *map[string]struct{}) {
	if len(e.Label) > 0 {
		(*recv)[e.Label] = struct{}{}
	}
	if len(e.HelpText) > 0 {
		(*recv)[e.HelpText] = struct{}{}
	}
	for _, v := range e.Attributes {
		if len(v) == 2 && len(v[0]) > 0 && len(v[1]) > 0 {
			switch v[0] {
			case `title`:
				fallthrough
			case `placeholder`:
				(*recv)[v[1]] = struct{}{}
			}
		}
	}
	for _, v := range e.Choices {
		if len(v.Group) > 0 && !IsExistsKey(recv, v.Group) && !com.StrIsNumeric(v.Group) {
			(*recv)[v.Group] = struct{}{}
		}
		if len(v.Option) == 2 && len(v.Option[1]) > 0 && !IsExistsKey(recv, v.Option[1]) && !com.StrIsNumeric(v.Option[1]) {
			(*recv)[v.Option[1]] = struct{}{}
		}
	}
}
