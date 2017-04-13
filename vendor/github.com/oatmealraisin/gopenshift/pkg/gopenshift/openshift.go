package gopenshift

import (
	"fmt"

	"github.com/openshift/origin/pkg/client"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api/meta"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/runtime"
)

type OpenShift struct {
	Context       string
	Token         string
	Factory       *clientcmd.Factory
	UserInterface client.UserInterface
	Mapper        meta.RESTMapper
	Typer         runtime.ObjectTyper
	Command       *cobra.Command
}

func New() *OpenShift {
	command := &cobra.Command{}
	cmdutil.AddPrinterFlags(command)
	factory := clientcmd.New(command.Flags())
	mapper, typer := factory.Object()
	return &OpenShift{
		Command: command,
		Factory: factory,
		Mapper:  mapper,
		Typer:   typer,
	}
}

// TODO: Stub
func (o *OpenShift) Describe() string {
	return ""
}

// TODO: Stub
func (o *OpenShift) Edit(resource, name string) string {
	return ""
}

// TODO: Stub
func (o *OpenShift) ListVolumes(resource string) string {
	return ""
}

// TODO: Stub
func (o *OpenShift) Logs(resource string) string {
	return ""
}

func (o *OpenShift) Project() (string, error) {
	openshiftConfig, err := o.Factory.OpenShiftClientConfig().RawConfig()
	if err != nil {
		return "", err
	}

	return openshiftConfig.Contexts[openshiftConfig.CurrentContext].Namespace, nil
}

// TODO: Stub
func (o *OpenShift) Projects() string {
	return ""
}

// TODO: Stub
func (o *OpenShift) Status() string {
	return ""
}

func (o *OpenShift) WhoAmI() (string, error) {
	if o.Token != "" {
		cfg, err := o.Factory.OpenShiftClientConfig().ClientConfig()
		if err != nil {
			return "", err
		}
		if len(cfg.BearerToken) == 0 {
			return "", fmt.Errorf("no token is currently in use for this session")
		}
		return cfg.BearerToken, nil
	}
	if o.Context != "" {
		cfg, err := o.Factory.OpenShiftClientConfig().RawConfig()
		if err != nil {
			return "", err
		}
		if len(cfg.CurrentContext) == 0 {
			return "", fmt.Errorf("no context has been set")
		}
		return cfg.CurrentContext, nil
	}

	client, _, err := o.Factory.Clients()
	if err != nil {
		return "", err
	}

	userInterface := client.Users()

	me, err := userInterface.Get("~")
	if err != nil {
		return "", err
	}

	return me.Name, nil
}
