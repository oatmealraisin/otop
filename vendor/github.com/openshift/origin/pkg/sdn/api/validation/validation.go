package validation

import (
	"fmt"
	"net"

	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/api/validation/path"
	"k8s.io/kubernetes/pkg/util/validation/field"

	sdnapi "github.com/openshift/origin/pkg/sdn/api"
	"github.com/openshift/origin/pkg/util/netutils"
)

// ValidateClusterNetwork tests if required fields in the ClusterNetwork are set.
func ValidateClusterNetwork(clusterNet *sdnapi.ClusterNetwork) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&clusterNet.ObjectMeta, false, path.ValidatePathSegmentName, field.NewPath("metadata"))

	clusterIPNet, err := netutils.ParseCIDRMask(clusterNet.Network)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("network"), clusterNet.Network, err.Error()))
	} else {
		ones, bitSize := clusterIPNet.Mask.Size()
		if uint32(bitSize-ones) <= clusterNet.HostSubnetLength {
			allErrs = append(allErrs, field.Invalid(field.NewPath("hostSubnetLength"), clusterNet.HostSubnetLength, "subnet length is greater than cluster Mask"))
		}
	}

	serviceIPNet, err := netutils.ParseCIDRMask(clusterNet.ServiceNetwork)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("serviceNetwork"), clusterNet.ServiceNetwork, err.Error()))
	}

	if (clusterIPNet != nil) && (serviceIPNet != nil) && clusterIPNet.Contains(serviceIPNet.IP) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("serviceNetwork"), clusterNet.ServiceNetwork, "service network overlaps with cluster network"))
	}
	if (serviceIPNet != nil) && (clusterIPNet != nil) && serviceIPNet.Contains(clusterIPNet.IP) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("network"), clusterNet.Network, "cluster network overlaps with service network"))
	}

	return allErrs
}

func validateNewNetwork(obj *sdnapi.ClusterNetwork, old *sdnapi.ClusterNetwork) *field.Error {
	oldNet, err := netutils.ParseCIDRMask(old.Network)
	if err != nil {
		// Shouldn't happen, but if the existing value is invalid, then any change should be an improvement...
		return nil
	}
	oldSize, _ := oldNet.Mask.Size()
	newNet, err := netutils.ParseCIDRMask(obj.Network)
	if err != nil {
		return field.Invalid(field.NewPath("network"), obj.Network, err.Error())
	}
	newSize, _ := newNet.Mask.Size()
	// oldSize/newSize is, eg the "16" in "10.1.0.0/16", so "newSize < oldSize" means
	// the new network is larger
	if newSize < oldSize && newNet.Contains(oldNet.IP) {
		return nil
	} else {
		return field.Invalid(field.NewPath("network"), obj.Network, "cannot change the cluster's network CIDR to a value that does not include the existing network.")
	}
}

func ValidateClusterNetworkUpdate(obj *sdnapi.ClusterNetwork, old *sdnapi.ClusterNetwork) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&obj.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateClusterNetwork(obj)...)

	if obj.Network != old.Network {
		err := validateNewNetwork(obj, old)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}
	if obj.HostSubnetLength != old.HostSubnetLength {
		allErrs = append(allErrs, field.Invalid(field.NewPath("hostSubnetLength"), obj.HostSubnetLength, "cannot change the cluster's hostSubnetLength midflight."))
	}
	if obj.ServiceNetwork != old.ServiceNetwork && old.ServiceNetwork != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("serviceNetwork"), obj.ServiceNetwork, "cannot change the cluster's serviceNetwork CIDR midflight."))
	}

	return allErrs
}

// ValidateHostSubnet tests fields for the host subnet, the host should be a network resolvable string,
//  and subnet should be a valid CIDR
func ValidateHostSubnet(hs *sdnapi.HostSubnet) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&hs.ObjectMeta, false, path.ValidatePathSegmentName, field.NewPath("metadata"))

	if hs.Host != hs.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("host"), hs.Host, fmt.Sprintf("must be the same as metadata.name: %q", hs.Name)))
	}

	if hs.Subnet == "" {
		// check if annotation exists, then let the Subnet field be empty
		if _, ok := hs.Annotations[sdnapi.AssignHostSubnetAnnotation]; !ok {
			allErrs = append(allErrs, field.Invalid(field.NewPath("subnet"), hs.Subnet, "field cannot be empty"))
		}
	} else {
		_, err := netutils.ParseCIDRMask(hs.Subnet)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("subnet"), hs.Subnet, err.Error()))
		}
	}
	if net.ParseIP(hs.HostIP) == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("hostIP"), hs.HostIP, "invalid IP address"))
	}
	return allErrs
}

func ValidateHostSubnetUpdate(obj *sdnapi.HostSubnet, old *sdnapi.HostSubnet) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&obj.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateHostSubnet(obj)...)

	if obj.Subnet != old.Subnet {
		allErrs = append(allErrs, field.Invalid(field.NewPath("subnet"), obj.Subnet, "cannot change the subnet lease midflight."))
	}

	return allErrs
}

// ValidateNetNamespace tests fields for a greater-than-zero NetID
func ValidateNetNamespace(netnamespace *sdnapi.NetNamespace) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&netnamespace.ObjectMeta, false, path.ValidatePathSegmentName, field.NewPath("metadata"))

	if netnamespace.NetName != netnamespace.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("netname"), netnamespace.NetName, fmt.Sprintf("must be the same as metadata.name: %q", netnamespace.Name)))
	}

	if err := sdnapi.ValidVNID(netnamespace.NetID); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("netid"), netnamespace.NetID, err.Error()))
	}
	return allErrs
}

func ValidateNetNamespaceUpdate(obj *sdnapi.NetNamespace, old *sdnapi.NetNamespace) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&obj.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateNetNamespace(obj)...)
	return allErrs
}

// ValidateEgressNetworkPolicy tests if required fields in the EgressNetworkPolicy are set.
func ValidateEgressNetworkPolicy(policy *sdnapi.EgressNetworkPolicy) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&policy.ObjectMeta, true, path.ValidatePathSegmentName, field.NewPath("metadata"))

	for i, rule := range policy.Spec.Egress {
		if rule.Type != sdnapi.EgressNetworkPolicyRuleAllow && rule.Type != sdnapi.EgressNetworkPolicyRuleDeny {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("egress").Index(i).Child("type"), rule.Type, "invalid policy type"))
		}

		_, err := netutils.ParseCIDRMask(rule.To.CIDRSelector)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("egress").Index(i).Child("to"), rule.To.CIDRSelector, err.Error()))
		}
	}

	if len(policy.Spec.Egress) > sdnapi.EgressNetworkPolicyMaxRules {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("egress"), "", ("too many egress rules (max 50)")))
	}

	return allErrs
}

func ValidateEgressNetworkPolicyUpdate(obj *sdnapi.EgressNetworkPolicy, old *sdnapi.EgressNetworkPolicy) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&obj.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateEgressNetworkPolicy(obj)...)
	return allErrs
}
