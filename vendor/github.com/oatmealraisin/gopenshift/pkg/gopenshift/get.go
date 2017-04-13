package gopenshift

import (
	kapi "k8s.io/kubernetes/pkg/api"
	ekapi "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"

	_ "github.com/openshift/origin/pkg/api/install"
	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	routeapi "github.com/openshift/origin/pkg/route/api"
)

// TODO: Since we'll be watching this stuff we should implement the watcher logic
//       from the client
// TODO: We are missing a lot of stuff for getting "all"
// TODO: Find a way to decide when to print all namespaces
func (o *OpenShift) Get(reqResource string) ([]runtime.Object, error) {
	var result []runtime.Object

	cmdNamespace, _, err := o.Factory.DefaultNamespace()
	if err != nil {
		return nil, err
	}

	// TODO: Comment on what this does

	r := resource.NewBuilder(o.Mapper, o.Typer, resource.ClientMapperFunc(o.Factory.ClientForMapping), runtime.UnstructuredJSONScheme).
		NamespaceParam(cmdNamespace).DefaultNamespace().AllNamespaces(false).
		SelectorParam("").
		ExportParam(false).
		ResourceTypeOrNameArgs(true, reqResource).
		ContinueOnError().
		Latest().
		Flatten().
		RequireObject(true).
		SingleResourceType().
		Do()

	err = r.Err()
	if err != nil {
		return nil, err
	}

	infos, err := r.Infos()
	if err != nil {
		return nil, err
	}

	//if len(infos) == 0 {
	//	return nil, fmt.Errorf("No items found")
	//}

	for _, info := range infos {
		result = append(result, info.Object)
	}

	return result, nil
}

func (o *OpenShift) GetPods() ([]*kapi.Pod, error) {
	pods := []*kapi.Pod{}

	objects, err := o.Get("po")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if pod := object.(*kapi.Pod); pod != nil {
			pods = append(pods, pod)
		}
	}

	return pods, nil
}

func (o *OpenShift) GetServices() ([]*kapi.Service, error) {
	services := []*kapi.Service{}

	objects, err := o.Get("svc")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if service := object.(*kapi.Service); service != nil {
			services = append(services, service)
		}
	}

	return services, nil
}

func (o *OpenShift) GetDeployments() ([]*ekapi.Deployment, error) {
	deployments := []*ekapi.Deployment{}

	objects, err := o.Get("deploy")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if deployment := object.(*ekapi.Deployment); deployment != nil {
			deployments = append(deployments, deployment)
		}
	}

	return deployments, nil
}

func (o *OpenShift) GetBuilds() ([]*buildapi.Build, error) {
	builds := []*buildapi.Build{}

	objects, err := o.Get("builds")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if build := object.(*buildapi.Build); build != nil {
			builds = append(builds, build)
		}
	}

	return builds, nil
}

func (o *OpenShift) GetRoutes() ([]*routeapi.Route, error) {
	routes := []*routeapi.Route{}

	objects, err := o.Get("routes")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if route := object.(*routeapi.Route); route != nil {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

func (o *OpenShift) GetDeploymentConfigs() ([]*deployapi.DeploymentConfig, error) {
	deploys := []*deployapi.DeploymentConfig{}

	objects, err := o.Get("dc")
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if deploy := object.(*deployapi.DeploymentConfig); deploy != nil {
			deploys = append(deploys, deploy)
		}
	}

	return deploys, nil
}
