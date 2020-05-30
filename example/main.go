package main

import (
	"fmt"
	"time"

	. "github.com/coscms/forms"
)

type Test struct {
	User     string
	Birthday string
}

func main() {
	//1.===================================
	startTime := time.Now()
	config := Unmarshal(`forms.json`)
	t := Test{
		User:     `webx`,
		Birthday: `1985`,
	}
	form := NewWithModelConfig(t, config)
	fmt.Println(form.Render())
	//return
	fmt.Println(`1.________________________________________CostTime:`, time.Now().Sub(startTime).Seconds(), `s`)
	fmt.Println(``)

	//2.===================================
	startTime = time.Now()
	form = New()
	fmt.Println(form.Init(config, t).ValidFromConfig().ParseFromConfig(true))

	fmt.Println(`2.________________________________________CostTime:`, time.Now().Sub(startTime).Seconds(), `s`)
	fmt.Println(``)

	//3.===================================
	startTime = time.Now()
	form = New()
	fmt.Println(form.Generate(t, `forms.json`))
	fmt.Println(`3.________________________________________CostTime:`, time.Now().Sub(startTime).Seconds(), `s`)
}
