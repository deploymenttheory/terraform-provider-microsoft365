package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// registrationNotConfigurableDetail explains why azureADRegistration cannot be changed.
const registrationNotConfigurableDetail = "The tenant's device registration policy reports azureADRegistration.isAdminConfigurable = false, " +
	"which means Microsoft Intune manages who can register devices and Microsoft Graph ignores changes to it. " +
	"Remove azure_ad_registration from the configuration, or set allowed_to_register to the current tenant value."

// ValidateConfig checks that users and groups are only configured for enumerated memberships.
func (r *DeviceRegistrationPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config DeviceRegistrationPolicyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberships := []struct {
		path  path.Path
		value attr.Value
	}{
		{path.Root("azure_ad_registration").AtName("allowed_to_register"), priorAttribute(config.AzureADRegistration, "allowed_to_register")},
		{path.Root("azure_ad_join").AtName("allowed_to_join"), priorAttribute(config.AzureADJoin, "allowed_to_join")},
		{path.Root("azure_ad_join").AtName("local_admins").AtName("registering_users"), priorAttribute(priorAttribute(config.AzureADJoin, "local_admins"), "registering_users")},
	}

	for _, membership := range memberships {
		if err := validateMembership(membership.value); err != nil {
			resp.Diagnostics.AddAttributeError(
				membership.path,
				"Invalid device registration membership",
				fmt.Sprintf("%s: %s", membership.path.String(), err.Error()),
			)
		}
	}
}

// validateMembership checks a configured membership object.
func validateMembership(value attr.Value) error {
	object, ok := value.(basetypes.ObjectValue)
	if !ok || object.IsNull() || object.IsUnknown() {
		return nil
	}

	attributes := object.Attributes()
	membershipType, ok := attributes["type"].(basetypes.StringValue)
	if !ok || membershipType.IsNull() || membershipType.IsUnknown() || membershipType.ValueString() == membershipTypeEnumerated {
		return nil
	}

	for _, name := range []string{"users", "groups"} {
		if set, ok := attributes[name].(basetypes.SetValue); ok && !set.IsNull() {
			return fmt.Errorf("%s %w, but type is %q", name, errMembershipListNotAllowed, membershipType.ValueString())
		}
	}
	return nil
}

// validateRequest validates the planned policy against the current tenant policy before the PUT.
func validateRequest(ctx context.Context, data *DeviceRegistrationPolicyResourceModel, remoteResource graphmodels.DeviceRegistrationPolicyable) error {
	tflog.Debug(ctx, "Starting device registration policy request validation")

	if err := validateMembership(priorAttribute(data.AzureADRegistration, "allowed_to_register")); err != nil {
		return fmt.Errorf("azure_ad_registration.allowed_to_register: %w", err)
	}
	if err := validateMembership(priorAttribute(data.AzureADJoin, "allowed_to_join")); err != nil {
		return fmt.Errorf("azure_ad_join.allowed_to_join: %w", err)
	}
	if err := validateMembership(priorAttribute(priorAttribute(data.AzureADJoin, "local_admins"), "registering_users")); err != nil {
		return fmt.Errorf("azure_ad_join.local_admins.registering_users: %w", err)
	}

	registration := remoteResource.GetAzureADRegistration()
	if registration == nil || registration.GetIsAdminConfigurable() == nil || *registration.GetIsAdminConfigurable() {
		return nil
	}

	planned, ok := priorAttribute(data.AzureADRegistration, "allowed_to_register").(basetypes.ObjectValue)
	if !ok || planned.IsNull() || planned.IsUnknown() {
		return nil
	}

	current := mapMembership(ctx, planned, registration.GetAllowedToRegister())
	if !membershipsEqual(planned, current) {
		return fmt.Errorf("%w: %s", errRegistrationNotConfigurable, registrationNotConfigurableDetail)
	}

	tflog.Debug(ctx, "Finished device registration policy request validation")
	return nil
}
