package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// newPolicyBody returns an empty request body of the SDK type that matches the policy type
func newPolicyBody(policyType string) graphmodels.MobilityManagementPolicyable {
	if policyType == policyTypeApp {
		return graphmodels.NewMobileAppManagementPolicy()
	}
	return graphmodels.NewMobileDeviceManagementPolicy()
}

// nullableValue is any Terraform value that can be null or unknown
type nullableValue interface {
	IsNull() bool
	IsUnknown() bool
}

// isConfigured reports whether a planned value was set in configuration and is known
func isConfigured(value nullableValue) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// stringChanged reports whether a configured string differs from state, or there is no state yet (create)
func stringChanged(plan types.String, state *MobilityManagementPolicyResourceModel, stateValue func(*MobilityManagementPolicyResourceModel) types.String) bool {
	if !isConfigured(plan) {
		return false
	}
	return state == nil || !plan.Equal(stateValue(state))
}

// constructSettingsRequest builds the PATCH body for every configured setting other than appliesTo.
// appliesTo cannot be PATCHed together with other properties, so it is sent by constructAppliesToRequest.
// On update only settings that differ from state are sent. Returns nil when there is nothing to write.
func constructSettingsRequest(
	ctx context.Context,
	plan *MobilityManagementPolicyResourceModel,
	state *MobilityManagementPolicyResourceModel,
) graphmodels.MobilityManagementPolicyable {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s settings request", ResourceName))

	requestBody := newPolicyBody(plan.PolicyType.ValueString())
	hasChanges := false

	if stringChanged(plan.ComplianceUrl, state, func(m *MobilityManagementPolicyResourceModel) types.String { return m.ComplianceUrl }) {
		convert.FrameworkToGraphString(plan.ComplianceUrl, requestBody.SetComplianceUrl)
		hasChanges = true
	}
	if stringChanged(plan.DiscoveryUrl, state, func(m *MobilityManagementPolicyResourceModel) types.String { return m.DiscoveryUrl }) {
		convert.FrameworkToGraphString(plan.DiscoveryUrl, requestBody.SetDiscoveryUrl)
		hasChanges = true
	}
	if stringChanged(plan.TermsOfUseUrl, state, func(m *MobilityManagementPolicyResourceModel) types.String { return m.TermsOfUseUrl }) {
		convert.FrameworkToGraphString(plan.TermsOfUseUrl, requestBody.SetTermsOfUseUrl)
		hasChanges = true
	}

	if deviceBody, ok := requestBody.(graphmodels.MobileDeviceManagementPolicyable); ok && isConfigured(plan.IsMdmEnrollmentDuringRegistrationDisabled) {
		if state == nil || !plan.IsMdmEnrollmentDuringRegistrationDisabled.Equal(state.IsMdmEnrollmentDuringRegistrationDisabled) {
			convert.FrameworkToGraphBool(plan.IsMdmEnrollmentDuringRegistrationDisabled, deviceBody.SetIsMdmEnrollmentDuringRegistrationDisabled)
			hasChanges = true
		}
	}

	if !hasChanges {
		tflog.Debug(ctx, fmt.Sprintf("No %s settings to write", ResourceName))
		return nil
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s settings", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody
}

// constructAppliesToRequest builds the PATCH body that sets appliesTo on its own.
// Only none and all are written this way; selected is applied by adding included group references.
func constructAppliesToRequest(ctx context.Context, policyType, appliesTo string) (graphmodels.MobilityManagementPolicyable, error) {
	requestBody := newPolicyBody(policyType)

	if err := convert.FrameworkToGraphEnum(types.StringValue(appliesTo), graphmodels.ParsePolicyScope, requestBody.SetAppliesTo); err != nil {
		return nil, fmt.Errorf("failed to set applies_to: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s appliesTo", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructGroupReference builds the $ref body that adds a group to includedGroups
func constructGroupReference(groupID string) graphmodels.ReferenceCreateable {
	requestBody := graphmodels.NewReferenceCreate()
	odataID := groupReferenceURL(groupID)
	requestBody.SetOdataId(&odataID)
	return requestBody
}
