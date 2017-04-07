package main

import (
	"fmt"
	"log"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
)

func main() {
	o := gopenshift.New()
	pods, err := o.GetPods()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, pod := range pods {
		fmt.Println(pod["name"])
	}
}
