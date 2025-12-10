package graphBetaAgentUser

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote user resource to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentUserResourceModel, remoteResource graphmodels.Userable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Required fields
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.UserPrincipalName = convert.GraphToFrameworkString(remoteResource.GetUserPrincipalName())
	data.MailNickname = convert.GraphToFrameworkString(remoteResource.GetMailNickname())

	// Agent user specific field - SDK has first-class getter
	data.AgentIdentityId = convert.GraphToFrameworkString(remoteResource.GetIdentityParentId())

	// Computed fields (read-only)
	data.Mail = convert.GraphToFrameworkString(remoteResource.GetMail())
	data.UserType = convert.GraphToFrameworkString(remoteResource.GetUserType())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.CreationType = convert.GraphToFrameworkString(remoteResource.GetCreationType())

	// Optional fields
	data.GivenName = convert.GraphToFrameworkString(remoteResource.GetGivenName())
	data.Surname = convert.GraphToFrameworkString(remoteResource.GetSurname())
	data.JobTitle = convert.GraphToFrameworkString(remoteResource.GetJobTitle())
	data.Department = convert.GraphToFrameworkString(remoteResource.GetDepartment())
	data.CompanyName = convert.GraphToFrameworkString(remoteResource.GetCompanyName())
	data.OfficeLocation = convert.GraphToFrameworkString(remoteResource.GetOfficeLocation())
	data.City = convert.GraphToFrameworkString(remoteResource.GetCity())
	data.State = convert.GraphToFrameworkString(remoteResource.GetState())
	data.Country = convert.GraphToFrameworkString(remoteResource.GetCountry())
	data.PostalCode = convert.GraphToFrameworkString(remoteResource.GetPostalCode())
	data.StreetAddress = convert.GraphToFrameworkString(remoteResource.GetStreetAddress())
	data.UsageLocation = convert.GraphToFrameworkString(remoteResource.GetUsageLocation())
	data.PreferredLanguage = convert.GraphToFrameworkString(remoteResource.GetPreferredLanguage())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
