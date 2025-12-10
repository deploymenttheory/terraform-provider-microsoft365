package graphBetaAgentUser

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body for creating an agent user.
// Uses the SDK's AgentUser model which auto-sets @odata.type to #microsoft.graph.agentUser
func constructResource(ctx context.Context, data *AgentUserResourceModel) (graphmodels.Userable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAgentUser()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.UserPrincipalName, requestBody.SetUserPrincipalName)
	convert.FrameworkToGraphString(data.MailNickname, requestBody.SetMailNickname)
	convert.FrameworkToGraphBool(data.AccountEnabled, requestBody.SetAccountEnabled)
	convert.FrameworkToGraphString(data.AgentIdentityId, requestBody.SetIdentityParentId)

	// Add sponsors using OData bind during creation (still needs additionalData)
	if !data.SponsorIds.IsNull() && !data.SponsorIds.IsUnknown() {
		var sponsors []string
		diags := data.SponsorIds.ElementsAs(ctx, &sponsors, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract sponsor_ids: %v", diags.Errors())
		}
		if len(sponsors) > 0 {
			sponsorUrls := make([]string, len(sponsors))
			for i, sponsorId := range sponsors {
				sponsorUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/directoryObjects/%s", sponsorId)
			}
			additionalData := make(map[string]any)
			additionalData["sponsors@odata.bind"] = sponsorUrls
			requestBody.SetAdditionalData(additionalData)
			tflog.Debug(ctx, fmt.Sprintf("Adding %d sponsors to agent user", len(sponsors)))
		}
	}

	convert.FrameworkToGraphString(data.GivenName, requestBody.SetGivenName)
	convert.FrameworkToGraphString(data.Surname, requestBody.SetSurname)
	convert.FrameworkToGraphString(data.JobTitle, requestBody.SetJobTitle)
	convert.FrameworkToGraphString(data.Department, requestBody.SetDepartment)
	convert.FrameworkToGraphString(data.CompanyName, requestBody.SetCompanyName)
	convert.FrameworkToGraphString(data.OfficeLocation, requestBody.SetOfficeLocation)
	convert.FrameworkToGraphString(data.City, requestBody.SetCity)
	convert.FrameworkToGraphString(data.State, requestBody.SetState)
	convert.FrameworkToGraphString(data.Country, requestBody.SetCountry)
	convert.FrameworkToGraphString(data.PostalCode, requestBody.SetPostalCode)
	convert.FrameworkToGraphString(data.StreetAddress, requestBody.SetStreetAddress)
	convert.FrameworkToGraphString(data.UsageLocation, requestBody.SetUsageLocation)
	convert.FrameworkToGraphString(data.PreferredLanguage, requestBody.SetPreferredLanguage)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructUpdateResource builds the request body for updating an agent user.
// Note: identityParentId cannot be updated, so it's not included here.
func constructUpdateResource(ctx context.Context, data *AgentUserResourceModel) (graphmodels.Userable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource for update", ResourceName))

	requestBody := graphmodels.NewAgentUser()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphBool(data.AccountEnabled, requestBody.SetAccountEnabled)
	convert.FrameworkToGraphString(data.GivenName, requestBody.SetGivenName)
	convert.FrameworkToGraphString(data.Surname, requestBody.SetSurname)
	convert.FrameworkToGraphString(data.JobTitle, requestBody.SetJobTitle)
	convert.FrameworkToGraphString(data.Department, requestBody.SetDepartment)
	convert.FrameworkToGraphString(data.CompanyName, requestBody.SetCompanyName)
	convert.FrameworkToGraphString(data.OfficeLocation, requestBody.SetOfficeLocation)
	convert.FrameworkToGraphString(data.City, requestBody.SetCity)
	convert.FrameworkToGraphString(data.State, requestBody.SetState)
	convert.FrameworkToGraphString(data.Country, requestBody.SetCountry)
	convert.FrameworkToGraphString(data.PostalCode, requestBody.SetPostalCode)
	convert.FrameworkToGraphString(data.StreetAddress, requestBody.SetStreetAddress)
	convert.FrameworkToGraphString(data.UsageLocation, requestBody.SetUsageLocation)
	convert.FrameworkToGraphString(data.PreferredLanguage, requestBody.SetPreferredLanguage)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource for update", ResourceName))

	return requestBody, nil
}
