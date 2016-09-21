package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	"github.com/oatmealraisin/otop/pkg/otop"
	gc "github.com/rthornton128/goncurses"
	"github.com/spf13/cobra"
)

const (
	// TODO: Finish these
	cliExplain = `TODO`
	cliLong    = `TODO`
	cliShort   = `OpenShift cluster viewer`
	cliUse     = `otop`
)

type OtopCmd struct {
	// TODO: Rename this something sensible
	Thing *otop.Otop
}

func NewCmdOtop() *cobra.Command {
	ocmd := &OtopCmd{
		Thing: &otop.Otop{},
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
			return ocmd.Run()
		},
	}

	ocmd.Thing.OpenShift = gopenshift.New()

	return command
}

func (cmd OtopCmd) Run() error {
	if err := cmd.CheckInput(); err != nil {
		return err
	}

	// For ease of reading/writing
	thing := cmd.Thing

	// Initialize goncurses. It's essential End() is called to ensure the
	// terminal isn't altered after the program ends
	mw, err := gc.Init()
	if err != nil {
		return err
	}

	thing.Window = mw

	// We need to catch ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		gc.End()
		os.Exit(1)
	}()

	defer gc.End()

	return thing.Run()
}

func (o OtopCmd) CheckInput() error {
	return nil
}
