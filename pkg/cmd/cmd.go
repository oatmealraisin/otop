package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	"github.com/oatmealraisin/otop/pkg/otop"
	gc "github.com/rthornton128/goncurses"
	"github.com/spf13/cobra"
)

const (
	// TODO: Finish these
	cliExplain = `TODO`
	cliLong    = `TODO`
	cliShort   = `OpenShift Resource Monitor`
	cliUse     = `otop`
)

type OtopCmd struct {
	FrontEnd *otop.Otop
	LogsFile string
}

func NewCmdOtop() *cobra.Command {
	ocmd := &OtopCmd{
		FrontEnd: &otop.Otop{},
	}

	command := &cobra.Command{
		Use:   cliUse,
		Short: cliShort,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ocmd.Run()
			if err != nil && strings.Contains(err.Error(), "connection refused") {
				fmt.Println("Could not connect to OpenShift server.")
				return nil
			}
			return err
		},
	}

	command.Flags().StringVar(&ocmd.LogsFile, "logs", "", "File to output logs to.")

	ocmd.FrontEnd.OpenShift = gopenshift.New()

	return command
}

func (cmd OtopCmd) Run() error {
	if err := cmd.CheckInput(); err != nil {
		return err
	}

	var logFile io.Writer
	// TODO: Make this better
	if cmd.LogsFile != "" {
		logFile, err := os.Create(cmd.LogsFile)
		if err != nil {
			return err
		}
		defer logFile.Close()
	}

	// For ease of reading/writing
	front := cmd.FrontEnd

	// Initialize goncurses. It's essential End() is called to ensure the
	// terminal isn't altered after the program ends
	mainWindow, err := gc.Init()
	if err != nil {
		return err
	}

	front.Window = mainWindow

	// We need to catch ctrl+c so that gc doesn't mess up the terminal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		gc.End()
		os.Exit(1)
	}()

	defer gc.End()

	return front.Run(logFile)
}

func (o OtopCmd) CheckInput() error {
	// TODO: Check if o.LogFile is a valid file, no windows support
	return nil
}
