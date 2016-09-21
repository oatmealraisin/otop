package main

import "C"

import (
	"log"
	"os"

	"github.com/oatmealraisin/otop/pkg/cmd"
)

func main() {

	otopCmd := cmd.NewCmdOtop()
	if err := otopCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
