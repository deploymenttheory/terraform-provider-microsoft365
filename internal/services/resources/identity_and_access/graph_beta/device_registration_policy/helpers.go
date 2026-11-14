package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"errors"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	errRemotePolicyRequired        = errors.New("the current device registration policy is required to build a full-replace request")
	errReadPlanValue               = errors.New("failed to read planned value")
	errUnsupportedMembershipType   = errors.New("unsupported membership type")
	errMembershipListNotAllowed    = errors.New(`can only be set when type is "enumerated"`)
	errRegistrationNotConfigurable = errors.New("azure_ad_registration.allowed_to_register cannot be changed")
)

// membershipTypeFromODataType converts the membership @odata.type to the schema type value.
func membershipTypeFromODataType(odataType *string) string {
	if odataType == nil {
		return ""
	}
	switch *odataType {
	case odataTypeAllMembership:
		return membershipTypeAll
	case odataTypeNoMembership:
		return membershipTypeNone
	case odataTypeEnumeratedMembership:
		return membershipTypeEnumerated
	default:
		return strings.TrimPrefix(*odataType, "#microsoft.graph.")
	}
}

// priorAttribute returns a nested attribute of a known object, or nil when it is not available.
func priorAttribute(value attr.Value, name string) attr.Value {
	object, ok := value.(basetypes.ObjectValue)
	if !ok || object.IsNull() || object.IsUnknown() {
		return nil
	}
	return object.Attributes()[name]
}

// isKnown reports whether a plan value was configured or carried from state, as opposed to unknown or null.
func isKnown(value attr.Value) bool {
	return value != nil && !value.IsNull() && !value.IsUnknown()
}

// membershipsEqual compares two membership objects by type and by users and groups, ignoring null versus empty.
func membershipsEqual(a, b types.Object) bool {
	if a.IsNull() || a.IsUnknown() || b.IsNull() || b.IsUnknown() {
		return a.Equal(b)
	}
	aAttributes, bAttributes := a.Attributes(), b.Attributes()
	if !aAttributes["type"].Equal(bAttributes["type"]) {
		return false
	}
	for _, name := range []string{"users", "groups"} {
		if !slices.Equal(sortedSetStrings(aAttributes[name]), sortedSetStrings(bAttributes[name])) {
			return false
		}
	}
	return true
}

// sortedSetStrings returns the sorted string elements of a set value.
func sortedSetStrings(value attr.Value) []string {
	set, ok := value.(basetypes.SetValue)
	if !ok || set.IsNull() || set.IsUnknown() {
		return []string{}
	}
	values := make([]string, 0, len(set.Elements()))
	for _, element := range set.Elements() {
		if stringValue, ok := element.(basetypes.StringValue); ok {
			values = append(values, stringValue.ValueString())
		}
	}
	slices.Sort(values)
	return values
}
