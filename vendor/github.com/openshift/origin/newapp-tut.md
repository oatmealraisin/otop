= pkg/generate/app/cli/newapp.go Tutorial

1. Grab env and parameters from c.Validate()
   Parses Env, Build-Env, and Parameters to combine duplicates and check for
	errs. We only catch env and param here.
2. c.ensureDockerSearch()
   Sets these DockerSearcher is not already set.
	TODO: Is this normally set before here?
3. `Resolve()`
	From the function description:
		Resolve transforms unstructured inputs (component names, templates,
		images) into a set of resolved components, or returns an error.
	Orders the inputs in the AppConfig to resemble something we can work futher
	with in order to figure out what we're doing.
	TODO: Finish
4. Check to make sure we have components/repos
5. `validateBuilders()`
	If using Source strategy, check to make sure you're using a s2i builder image
6. If name manually set, check it
7. Validate exposed ports
8. `validateOutputImageReferences()`
9. `installComponents()`
10. `c.BuildPipelines()`
TODO: Finish
