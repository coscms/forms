package config_test

import (
	"testing"

	"github.com/coscms/forms"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestSplitGroup(t *testing.T) {
	langset := forms.LangSet(`test`, ``, &config.Language{ID: `zh-CN`}, &config.Language{ID: `en`})
	langset.Elements(&fields.Field{
		OrigName:  `1`,
		LabelCols: 0,
		FieldCols: 0,
	})
	r := config.SplitGroup([]config.FormElement{
		&fields.Field{
			OrigName:  `1`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `2`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `3`,
			LabelCols: 0,
			FieldCols: 8,
		},
		&fields.Field{
			OrigName:  `4`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `5`,
			LabelCols: 0,
			FieldCols: 4,
			Errors:    []string{`Test`},
		},
	})
	assert.Equal(t, 3, len(r))

	r2 := config.SplitGroup([]config.FormElement{
		&fields.Field{
			OrigName:  `1`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `2`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `3`,
			LabelCols: 0,
			FieldCols: 8,
		},
		langset,
		&fields.Field{
			OrigName:  `4`,
			LabelCols: 0,
			FieldCols: 4,
		},
		&fields.Field{
			OrigName:  `5`,
			LabelCols: 0,
			FieldCols: 4,
			Errors:    []string{`Test`},
		},
	})
	assert.Equal(t, 4, len(r2))
	com.Dump(r)
}
