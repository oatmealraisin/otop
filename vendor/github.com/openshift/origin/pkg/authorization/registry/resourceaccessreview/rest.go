package resourceaccessreview

import (
	"errors"
	"fmt"

	kapi "k8s.io/kubernetes/pkg/api"
	kapierrors "k8s.io/kubernetes/pkg/api/errors"
	kauthorizer "k8s.io/kubernetes/pkg/auth/authorizer"
	"k8s.io/kubernetes/pkg/auth/user"
	"k8s.io/kubernetes/pkg/runtime"

	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	authorizationvalidation "github.com/openshift/origin/pkg/authorization/api/validation"
	"github.com/openshift/origin/pkg/authorization/authorizer"
)

// REST implements the RESTStorage interface in terms of an Registry.
type REST struct {
	authorizer     kauthorizer.Authorizer
	subjectLocator authorizer.SubjectLocator
}

// NewREST creates a new REST for policies.
func NewREST(authorizer kauthorizer.Authorizer, subjectLocator authorizer.SubjectLocator) *REST {
	return &REST{authorizer, subjectLocator}
}

// New creates a new ResourceAccessReview object
func (r *REST) New() runtime.Object {
	return &authorizationapi.ResourceAccessReview{}
}

// Create registers a given new ResourceAccessReview instance to r.registry.
func (r *REST) Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error) {
	resourceAccessReview, ok := obj.(*authorizationapi.ResourceAccessReview)
	if !ok {
		return nil, kapierrors.NewBadRequest(fmt.Sprintf("not a resourceAccessReview: %#v", obj))
	}
	if errs := authorizationvalidation.ValidateResourceAccessReview(resourceAccessReview); len(errs) > 0 {
		return nil, kapierrors.NewInvalid(authorizationapi.Kind(resourceAccessReview.Kind), "", errs)
	}

	user, ok := kapi.UserFrom(ctx)
	if !ok {
		return nil, kapierrors.NewInternalError(errors.New("missing user on request"))
	}

	// if a namespace is present on the request, then the namespace on the on the RAR is overwritten.
	// This is to support backwards compatibility.  To have gotten here in this state, it means that
	// the authorizer decided that a user could run an RAR against this namespace
	if namespace := kapi.NamespaceValue(ctx); len(namespace) > 0 {
		resourceAccessReview.Action.Namespace = namespace

	} else if err := r.isAllowed(user, resourceAccessReview); err != nil {
		// this check is mutually exclusive to the condition above.  localSAR and localRAR both clear the namespace before delegating their calls
		// We only need to check if the RAR is allowed **again** if the authorizer didn't already approve the request for a legacy call.
		return nil, err
	}

	attributes := authorizer.ToDefaultAuthorizationAttributes(nil, resourceAccessReview.Action.Namespace, resourceAccessReview.Action)
	users, groups, err := r.subjectLocator.GetAllowedSubjects(attributes)

	response := &authorizationapi.ResourceAccessReviewResponse{
		Namespace: resourceAccessReview.Action.Namespace,
		Users:     users,
		Groups:    groups,
	}
	if err != nil {
		response.EvaluationError = err.Error()
	}

	return response, nil
}

// isAllowed checks to see if the current user has rights to issue a LocalSubjectAccessReview on the namespace they're attempting to access
func (r *REST) isAllowed(user user.Info, rar *authorizationapi.ResourceAccessReview) error {
	localRARAttributes := kauthorizer.AttributesRecord{
		User:            user,
		Verb:            "create",
		Namespace:       rar.Action.Namespace,
		Resource:        "localresourceaccessreviews",
		ResourceRequest: true,
	}
	allowed, reason, err := r.authorizer.Authorize(localRARAttributes)

	if err != nil {
		return kapierrors.NewForbidden(authorizationapi.Resource(localRARAttributes.GetResource()), localRARAttributes.GetName(), err)
	}
	if !allowed {
		forbiddenError := kapierrors.NewForbidden(authorizationapi.Resource(localRARAttributes.GetResource()), localRARAttributes.GetName(), errors.New("") /*discarded*/)
		forbiddenError.ErrStatus.Message = reason
		return forbiddenError
	}

	return nil
}
