package forms

import (
	"fmt"
	"testing"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/webx-top/com"
)

func TestForms(t *testing.T) {
	type Data struct {
		Test string
	}
	mp := map[string]interface{}{
		`name`: `test`,
		`age`:  20,
		`items`: map[string]string{
			`itemK1`: `itemV1`,
		},
		`data`: &Data{
			Test: `test-data`,
		},
		`list`: []string{
			`1`, `2`,
		},
		`listData`: []*Data{
			&Data{
				Test: `test-listdata-1`,
			}, &Data{
				Test: `test-listdata-2`,
			},
		},
	}
	cfg := NewConfig()
	cfg.AddElement(&config.Element{
		ID:    `input-name`,
		Type:  `text`,
		Name:  `name`,
		Label: `名称`,
	}, &config.Element{
		ID:    `input-items-k1`,
		Type:  `text`,
		Name:  `items.itemK1`,
		Label: `Item K1`,
	}, &config.Element{
		ID:    `input-data-test`,
		Type:  `text`,
		Name:  `data.test`,
		Label: `Data`,
	}, &config.Element{
		ID:    `input-list-0`,
		Type:  `text`,
		Name:  `list.0`,
		Label: `List 0`,
	}, &config.Element{
		ID:    `input-list-2`,
		Type:  `text`,
		Name:  `list.2`,
		Label: `List 2`,
	}, &config.Element{
		ID:    `input-listdata-0`,
		Type:  `text`,
		Name:  `listData.0.test`,
		Label: `ListData 0`,
	})
	form := NewWithModelConfig(mp, cfg)
	com.Dump(form.Data())
	result := form.String()
	fmt.Println(result)
}

func TestParseConfig(t *testing.T) {
	cfg := config.Config{
		Elements: []*config.Element{
			{
				ID:   ``,
				Type: `text`,
				Name: `test`,
			},
		},
	}
	f := NewForms(New())
	f.Theme = common.BOOTSTRAP
	f.Init(&cfg)
	f.ParseFromConfig(true)
	com.Dump(cfg.Clone())
}
