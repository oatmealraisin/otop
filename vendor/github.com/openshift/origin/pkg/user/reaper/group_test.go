package reaper

import (
	"reflect"
	"testing"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
	"k8s.io/kubernetes/pkg/client/testing/core"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/davecgh/go-spew/spew"
	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	"github.com/openshift/origin/pkg/client/testclient"

	// install all APIs
	_ "github.com/openshift/origin/pkg/api/install"
)

var (
	groupsResource              = unversioned.GroupVersionResource{Group: "", Version: "", Resource: "groups"}
	clusterRoleBindingsResource = unversioned.GroupVersionResource{Group: "", Version: "", Resource: "clusterrolebindings"}
	roleBindingsResource        = unversioned.GroupVersionResource{Group: "", Version: "", Resource: "rolebindings"}
)

func TestGroupReaper(t *testing.T) {
	tests := []struct {
		name     string
		group    string
		objects  []runtime.Object
		expected []interface{}
	}{
		{
			name:    "no objects",
			group:   "mygroup",
			objects: []runtime.Object{},
			expected: []interface{}{
				core.DeleteActionImpl{ActionImpl: core.ActionImpl{Verb: "delete", Resource: groupsResource}, Name: "mygroup"},
			},
		},
		{
			name:  "cluster bindings",
			group: "mygroup",
			objects: []runtime.Object{
				&authorizationapi.ClusterRoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-no-subjects"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{},
				},
				&authorizationapi.ClusterRoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-one-subject"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{{Name: "mygroup", Kind: "Group"}},
				},
				&authorizationapi.ClusterRoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-mismatched-subject"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{{Name: "mygroup"}, {Name: "mygroup", Kind: "User"}, {Name: "mygroup", Kind: "Other"}},
				},
			},
			expected: []interface{}{
				core.UpdateActionImpl{ActionImpl: core.ActionImpl{Verb: "update", Resource: clusterRoleBindingsResource}, Object: &authorizationapi.ClusterRoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-one-subject"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{},
				}},
				core.DeleteActionImpl{ActionImpl: core.ActionImpl{Verb: "delete", Resource: groupsResource}, Name: "mygroup"},
			},
		},
		{
			name:  "namespaced bindings",
			group: "mygroup",
			objects: []runtime.Object{
				&authorizationapi.RoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-no-subjects", Namespace: "ns1"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{},
				},
				&authorizationapi.RoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-one-subject", Namespace: "ns2"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{{Name: "mygroup", Kind: "Group"}},
				},
				&authorizationapi.RoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-mismatched-subject", Namespace: "ns3"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{{Name: "mygroup"}, {Name: "mygroup", Kind: "User"}, {Name: "mygroup", Kind: "Other"}},
				},
			},
			expected: []interface{}{
				core.UpdateActionImpl{ActionImpl: core.ActionImpl{Verb: "update", Resource: roleBindingsResource, Namespace: "ns2"}, Object: &authorizationapi.RoleBinding{
					ObjectMeta: kapi.ObjectMeta{Name: "binding-one-subject", Namespace: "ns2"},
					RoleRef:    kapi.ObjectReference{Name: "role"},
					Subjects:   []kapi.ObjectReference{},
				}},
				core.DeleteActionImpl{ActionImpl: core.ActionImpl{Verb: "delete", Resource: groupsResource}, Name: "mygroup"},
			},
		},
		{
			name:  "sccs",
			group: "mygroup",
			objects: []runtime.Object{
				&kapi.SecurityContextConstraints{
					ObjectMeta: kapi.ObjectMeta{Name: "scc-no-subjects"},
					Groups:     []string{},
				},
				&kapi.SecurityContextConstraints{
					ObjectMeta: kapi.ObjectMeta{Name: "scc-one-subject"},
					Groups:     []string{"mygroup"},
				},
				&kapi.SecurityContextConstraints{
					ObjectMeta: kapi.ObjectMeta{Name: "scc-mismatched-subjects"},
					Users:      []string{"mygroup"},
					Groups:     []string{"mygroup2"},
				},
			},
			expected: []interface{}{
				core.UpdateActionImpl{ActionImpl: core.ActionImpl{Verb: "update", Resource: unversioned.GroupVersionResource{Resource: "securitycontextconstraints"}}, Object: &kapi.SecurityContextConstraints{
					ObjectMeta: kapi.ObjectMeta{Name: "scc-one-subject"},
					Groups:     []string{},
				}},
				core.DeleteActionImpl{ActionImpl: core.ActionImpl{Verb: "delete", Resource: groupsResource}, Name: "mygroup"},
			},
		},
	}

	for _, test := range tests {
		tc := testclient.NewSimpleFake(test.objects...)
		ktc := fake.NewSimpleClientset(test.objects...)

		actual := []interface{}{}
		oreactor := func(action core.Action) (handled bool, ret runtime.Object, err error) {
			actual = append(actual, action)
			return false, nil, nil
		}
		kreactor := func(action core.Action) (handled bool, ret runtime.Object, err error) {
			actual = append(actual, action)
			return false, nil, nil
		}

		tc.PrependReactor("update", "*", oreactor)
		tc.PrependReactor("delete", "*", oreactor)
		ktc.PrependReactor("update", "*", kreactor)
		ktc.PrependReactor("delete", "*", kreactor)

		reaper := NewGroupReaper(tc, tc, tc, ktc.Core())
		err := reaper.Stop("", test.group, 0, nil)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", test.name, err)
		}

		if !reflect.DeepEqual(test.expected, actual) {
			for i, x := range test.expected {
				t.Logf("Expected %d: %s", i, spew.Sprint(x))
			}
			for i, x := range actual {
				t.Logf("Actual %d:   %s", i, spew.Sprint(x))
			}
			t.Errorf("%s: unexpected actions", test.name)
		}
	}
}
