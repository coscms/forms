package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	. "github.com/coscms/forms"
	_ "github.com/coscms/forms/defaults"
)

type Test struct {
	User     string
	Birthday string
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
	"Language[en][meta_title]",
	"Language[zh-CN][meta_title]",
	"Language[en][meta_kewwords]",
	"Language[zh-CN][meta_kewwords]",
}

func main() {
	//1.===================================
	startTime := time.Now()
	config, err := UnmarshalFile(`forms.json`)
	if err != nil {
		log.Println(err)
	}
	t := Test{
		User:     `webx`,
		Birthday: `1985`,
	}
	b, _ := json.MarshalIndent(config.GetNames(), ``, " ")
	println(string(b))
	for _, name := range expected {
		if !config.HasName(name) {
			panic(`not found "` + name + `"`)
		}
	}
	form := NewWithModelConfig(t, config)
	fmt.Println(form.Render())
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

	b, _ = json.MarshalIndent(form, ``, " ")
	println(string(b))
}
