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

package config

type Config struct {
	ID             string                 `json:"id"`                       // 表单ID，用于区分不同表单
	Theme          string                 `json:"theme"`                    // 表单主题
	Template       string                 `json:"template"`                 // 表单模板
	Method         string                 `json:"method"`                   // 表单提交方法
	Action         string                 `json:"action"`                   // 表单提交地址
	Attributes     [][]string             `json:"attributes"`               // 表单属性
	WithButtons    bool                   `json:"withButtons"`              // 是否显示表单按钮
	Buttons        []string               `json:"buttons"`                  // 表单按钮
	BtnsTemplate   string                 `json:"btnsTemplate"`             // 表单按钮模板
	Elements       []*Element             `json:"elements"`                 // 表单元素
	Languages      []*Language            `json:"languages"`                // 表单多语言支持
	Data           map[string]interface{} `json:"data,omitempty"`           // 额外数据
	TrimNamePrefix string                 `json:"trimNamePrefix,omitempty"` // 去除字段名前缀
}

func (c *Config) Merge(source *Config) *Config {
	if len(c.ID) == 0 && len(source.ID) > 0 {
		c.ID = source.ID
	}
	if len(c.Theme) == 0 && len(source.Theme) > 0 {
		c.Theme = source.Theme
	}
	if len(c.Template) == 0 && len(source.Template) > 0 {
		c.Template = source.Template
	}
	if len(c.Method) == 0 && len(source.Method) > 0 {
		c.Method = source.Method
	}
	if len(c.Action) == 0 && len(source.Action) > 0 {
		c.Action = source.Action
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
	for _, v := range source.Buttons {
		for _, v2 := range c.Buttons {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			c.Buttons = append(c.Buttons, v)
		} else {
			found = false
		}
	}
	if c.WithButtons != source.WithButtons {
		c.WithButtons = source.WithButtons
	}
	if len(c.BtnsTemplate) == 0 && len(source.BtnsTemplate) > 0 {
		c.BtnsTemplate = source.BtnsTemplate
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
	if len(c.TrimNamePrefix) == 0 && len(source.TrimNamePrefix) > 0 {
		c.TrimNamePrefix = source.TrimNamePrefix
	}
	return c
}

func (c *Config) AddElement(elements ...*Element) *Config {
	c.Elements = append(c.Elements, elements...)
	return c
}

func (c *Config) AddLanguage(languages ...*Language) *Config {
	c.Languages = append(c.Languages, languages...)
	return c
}

func (c *Config) AddButton(buttons ...string) *Config {
	c.Buttons = append(c.Buttons, buttons...)
	return c
}

func (c *Config) AddAttribute(attributes ...string) *Config {
	c.Attributes = append(c.Attributes, attributes)
	return c
}

func (c *Config) Set(name string, value interface{}) *Config {
	if c.Data == nil {
		c.Data = map[string]interface{}{}
	}
	c.Data[name] = value
	return c
}

func (c *Config) Clone() *Config {
	elements := make([]*Element, len(c.Elements))
	for index, elem := range c.Elements {
		elements[index] = elem.Clone()
	}
	languages := make([]*Language, len(c.Languages))
	for index, value := range c.Languages {
		languages[index] = value.Clone()
	}
	r := &Config{
		ID:           c.ID,
		Theme:        c.Theme,
		Template:     c.Template,
		Method:       c.Method,
		Action:       c.Action,
		Attributes:   make([][]string, len(c.Attributes)),
		WithButtons:  c.WithButtons,
		Buttons:      make([]string, len(c.Buttons)),
		BtnsTemplate: c.BtnsTemplate,
		Elements:     elements,
		Languages:    languages,
		Data:         map[string]interface{}{},
	}
	copy(r.Buttons, c.Buttons)
	for k, v := range c.Data {
		r.Data[k] = v
	}
	for k, v := range c.Attributes {
		cv := make([]string, len(v))
		copy(cv, v)
		r.Attributes[k] = cv
	}
	return r
}

func (c *Config) HasName(name string) bool {
	return c.hasName(name, c.Elements, c.Languages)
}

func (c *Config) hasName(name string, elements []*Element, languages []*Language) bool {
	for _, elem := range elements {
		if elem.Name == name {
			return elem.Type != `langset` && elem.Type != `fieldset`
		}
		if elem.Type == `langset` {
			if c.hasName(name, elem.Elements, elem.Languages) {
				return true
			}
			continue
		}
		if elem.Type == `fieldset` {
			if c.hasName(name, elem.Elements, languages) {
				return true
			}
			continue
		}
		if len(languages) == 0 {
			continue
		}
		for _, lang := range languages {
			if lang.HasName(name) || name == lang.Name(elem.Name) {
				return true
			}
		}
	}
	return false
}

func (c *Config) GetNames() []string {
	return getNames(c.Elements, c.Languages)
}

func (c *Config) SetDefaultValue(fieldDefaultValue func(fieldName string) string) {
	if fieldDefaultValue != nil {
		setDefaultValue(c.Elements, c.Languages, fieldDefaultValue)
	}
}

func (c *Config) SetValue(fieldValue func(fieldName string) string) {
	if fieldValue != nil {
		setValue(c.Elements, c.Languages, fieldValue)
	}
}

func (c *Config) GetValue(fieldValue func(fieldName string, fieldValue string) error) error {
	if fieldValue != nil {
		return getValue(c.Elements, c.Languages, fieldValue)
	}
	return nil
}
