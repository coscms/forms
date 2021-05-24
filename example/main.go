package main

import (
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
}
