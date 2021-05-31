package fields

import (
	"testing"

	"github.com/coscms/forms/common"
	_ "github.com/coscms/forms/defaults"
	"github.com/stretchr/testify/assert"
)

func TestTextField(t *testing.T) {
	f := FieldWithType(`title`, common.TEXT)
	f.SetStyle(`base`)
	f.AddClass(`form-control`).AddClass(`row`)

	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`, `row`}), f.Classes)
	assert.Equal(t, `<input type="text" name="title" class="form-control row">`, f.String())
	f.data = nil

	f.RemoveClass(`row`)
	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`}), f.Classes)
	assert.Equal(t, `<input type="text" name="title" class="form-control">`, f.String())
	f.data = nil

	f.RemoveClass(`form-control`)
	assert.Equal(t, common.HTMLAttrValues([]string{}), f.Classes)

	assert.Equal(t, `<input type="text" name="title">`, f.String())
}

func TestCheckboxField(t *testing.T) {
	f := FieldWithType(`title`, common.CHECKBOX)
	f.SetStyle(`base`)
	f.AddChoice(`value1`, `text1`)
	f.AddChoice(`value2`, `text2`, true)
	f.AddChoice(`value3`, `text3`)
	f.AddChoice(`value4`, `text4`)

	assert.Equal(t, "<label for=\"value1\">text1</label>\n<input type=\"checkbox\" name=\"title\" value=\"value1\">", f.String())
}

func TestSelectField(t *testing.T) {
	f := FieldWithType(`title`, common.SELECT)
	f.SetStyle(`base`)
	f.AddChoice(`value1`, `text1`)
	f.AddChoice(`value2`, `text2`, true)
	f.AddChoice(`value3`, `text3`)
	f.AddChoice(`value4`, `text4`)

	assert.Equal(t, "<select name=\"title\"><option value=\"value1\">text1</option><option value=\"value2\" selected=\"selected\">text2</option><option value=\"value3\">text3</option><option value=\"value4\">text4</option></select>", f.String())
}
