package fields

import (
	"testing"

	"github.com/coscms/forms/common"
	_ "github.com/coscms/forms/defaults"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	f := FieldWithType(`title`, common.TEXT)
	f.SetStyle(`base`)
	f.AddClass(`form-control`).AddClass(`row`)

	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`, `row`}), f.Class)
	assert.Equal(t, `<input type="text" name="title" class="form-control row">`, f.String())
	f.data = nil

	f.RemoveClass(`row`)
	assert.Equal(t, common.HTMLAttrValues([]string{`form-control`}), f.Class)
	assert.Equal(t, `<input type="text" name="title" class="form-control">`, f.String())
	f.data = nil

	f.RemoveClass(`form-control`)
	assert.Equal(t, common.HTMLAttrValues([]string{}), f.Class)

	assert.Equal(t, `<input type="text" name="title">`, f.String())
}
