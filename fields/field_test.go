package fields

import (
	"strings"
	"testing"

	"github.com/coscms/forms/common"
	_ "github.com/coscms/forms/defaults"
	"github.com/stretchr/testify/assert"
)

func TestTextField(t *testing.T) {
	f := FieldWithType(`title`, common.TEXT)
	f.SetTheme(`base`)
	f.AddClass(`form-control`).AddClass(`row`)

	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`, `row`}), f.Classes)
	assert.Equal(t, `<input type="text" name="title" class="form-control row">`, strings.TrimSpace(f.String()))
	f.data = nil

	f.RemoveClass(`row`)
	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`}), f.Classes)
	assert.Equal(t, `<input type="text" name="title" class="form-control">`, strings.TrimSpace(f.String()))
	f.data = nil

	f.RemoveClass(`form-control`)
	assert.Equal(t, common.HTMLAttrValues([]string{}), f.Classes)

	assert.Equal(t, `<input type="text" name="title">`, strings.TrimSpace(f.String()))
}

func TestCheckboxField(t *testing.T) {
	f := FieldWithType(`title`, common.CHECKBOX)
	f.SetTheme(`base`)
	f.AddChoice(`value1`, `text1`)
	f.AddChoice(`value2`, `text2`, true)
	f.AddChoice(`value3`, `text3`)
	f.AddChoice(`value4`, `text4`)

	assert.Equal(t, "\n<label for=\"value1\">text1</label>\n<input type=\"checkbox\" name=\"title\" value=\"value1\">\n<label for=\"value2\">text2</label>\n<input type=\"checkbox\" name=\"title\" value=\"value2\" checked=\"checked\">\n<label for=\"value3\">text3</label>\n<input type=\"checkbox\" name=\"title\" value=\"value3\">\n<label for=\"value4\">text4</label>\n<input type=\"checkbox\" name=\"title\" value=\"value4\">", f.String())
}

func TestSelectField(t *testing.T) {
	f := FieldWithType(`title`, common.SELECT)
	f.SetTheme(`base`)
	f.AddChoice(`value1`, `text1`)
	f.AddChoice(`value2`, `text2`, true)
	f.AddChoice(`value3`, `text3`)
	f.AddChoice(`value4`, `text4`)

	assert.Equal(t, "\n<select name=\"title\">\n        <option value=\"value1\">text1</option>\n        <option value=\"value2\" selected=\"selected\">text2</option>\n        <option value=\"value3\">text3</option>\n        <option value=\"value4\">text4</option>\n</select>", f.String())
}
