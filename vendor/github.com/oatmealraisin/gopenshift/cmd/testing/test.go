package main

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
)

func main() {
	oc := gopenshift.New()

	testCases := []struct {
		// The elements of a test
		TestFn func(*gopenshift.OpenShift) error
		Name   string
	}{
		{
			Name:   "TestGet",
			TestFn: TestGet,
		},
		{
			Name:   "TestWhoAmI",
			TestFn: TestWhoAmI,
		},
		{
			Name:   "TestGetPods",
			TestFn: TestGetPods,
		},
	}

	for _, test := range testCases {
		// How we run the tests
		fmt.Printf("\nRunning %s\n", test.Name)
		err := test.TestFn(oc)
		if err != nil {
			fmt.Printf("FAILURE: Test %s failed with\n%s\n", test.Name, err.Error())
			return
		}
		fmt.Printf("SUCCESS - %s\n", test.Name)
	}
}

func TestGet(oc *gopenshift.OpenShift) error {
	for _, resource := range []string{
		"pods",
		"po",
		"services",
		"routes",
		"builds",
	} {
		fmt.Println(resource)
		_, err := oc.Get(resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestGetPods(oc *gopenshift.OpenShift) error {
	pods, err := oc.GetPods()
	if err != nil {
		return err
	}

	if len(pods) == 0 {
		return fmt.Errorf("No pods found!")
	}

	for _, pod := range pods {
		fmt.Println(pod.Name)
	}

	return nil
}

func TestWhoAmI(oc *gopenshift.OpenShift) error {
	return nil
}
