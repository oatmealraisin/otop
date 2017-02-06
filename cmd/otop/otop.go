package main

import "C"
import (
	"fmt"
	"os"

	"github.com/oatmealraisin/otop/pkg/cmd"
)

func main() {
	otopCmd := cmd.NewCmdOtop()
	if err := otopCmd.Execute(); err != nil {
		//log.Fatal(err)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
