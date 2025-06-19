package graphBetaTermsAndConditions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote terms and conditions to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data TermsAndConditionsResourceModel, termsAndConditions graphmodels.TermsAndConditionsable) TermsAndConditionsResourceModel {
	if termsAndConditions == nil {
		tflog.Debug(ctx, "Remote terms and conditions is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(termsAndConditions.GetId())
	data.DisplayName = convert.GraphToFrameworkString(termsAndConditions.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(termsAndConditions.GetDescription())
	data.Title = convert.GraphToFrameworkString(termsAndConditions.GetTitle())
	data.BodyText = convert.GraphToFrameworkString(termsAndConditions.GetBodyText())
	data.AcceptanceStatement = convert.GraphToFrameworkString(termsAndConditions.GetAcceptanceStatement())
	data.Version = convert.GraphToFrameworkInt32(termsAndConditions.GetVersion())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, termsAndConditions.GetRoleScopeTagIds())
	data.CreatedDateTime = convert.GraphToFrameworkTime(termsAndConditions.GetCreatedDateTime())
	data.ModifiedDateTime = convert.GraphToFrameworkTime(termsAndConditions.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
