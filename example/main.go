package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/coscms/forms"
	_ "github.com/coscms/forms/defaults"
)

type Test struct {
	User     string
	Birthday string
	Show     string
}

var expected = []string{
	"user",
	"user1",
	"user2",
	"birthday",
	"Language[en][title]",
	"Language[zh-CN][title]",
	"Language[en][content]",
	"Language[zh-CN][content]",
	"Language[en][metaTitle]",
	"Language[zh-CN][metaTitle]",
	"Language[en][metaKeywords]",
	"Language[zh-CN][metaKeywords]",
}

func main() {
	//1.===================================
	startTime := time.Now()
	config, err := UnmarshalFile(`forms.json`)
	if err != nil {
		log.Println(err)
	}
	//com.Dump(config.GetMultilingualText())
	//return
	t := Test{
		User:     `webx`,
		Birthday: `1985`,
		Show:     `0`,
	}
	b, _ := json.MarshalIndent(config.GetNames(), ``, " ")
	println(string(b))
	for _, name := range expected {
		if !config.HasName(name) {
			panic(`not found "` + name + `"`)
		}
	}
	form := NewWithModelConfig(t, config)
	htmlResult := form.Render()
	fmt.Println(htmlResult)
	os.Mkdir(`_test`, os.ModePerm)
	os.WriteFile(`_test/test.html`, []byte(`<html><head>
	<link rel="stylesheet" href="https://www.coscms.com/public/assets/backend/js/bootstrap/dist/css/bootstrap.min.css?t=0" />
	</head><body><div class="container">`+htmlResult+`</div><script src="https://www.coscms.com/public/assets/backend/js/jquery3.6.min.js?t=0"></script>
	<script type="text/javascript" src="https://www.coscms.com/public/assets/backend/js/bootstrap/dist/js/bootstrap.min.js?t=0"></script></body></html>`), os.ModePerm)
	//return
	fmt.Println(`1.________________________________________CostTime:`, time.Since(startTime).Seconds(), `s`)
	fmt.Println(``)

	//2.===================================
	startTime = time.Now()
	form = New()
	fmt.Println(form.Init(config, t).ValidFromConfig().ParseFromConfig(true))

	fmt.Println(`2.________________________________________CostTime:`, time.Since(startTime).Seconds(), `s`)
	fmt.Println(``)

	//3.===================================
	startTime = time.Now()
	form = New()
	form.Generate(t, `forms.json`)
	fmt.Println(form)
	fmt.Println(`3.________________________________________CostTime:`, time.Since(startTime).Seconds(), `s`)

	//b, _ = json.MarshalIndent(form, ``, " ")
	//println(string(b))
}
