package graphBetaTermsAndConditions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote terms and conditions to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data TermsAndConditionsResourceModel, termsAndConditions graphmodels.TermsAndConditionsable) TermsAndConditionsResourceModel {
	if termsAndConditions == nil {
		tflog.Debug(ctx, "Remote terms and conditions is nil")
		return data
	}

	data.ID = state.StringPointerValue(termsAndConditions.GetId())
	data.DisplayName = state.StringPointerValue(termsAndConditions.GetDisplayName())
	data.Description = state.StringPointerValue(termsAndConditions.GetDescription())
	data.Title = state.StringPointerValue(termsAndConditions.GetTitle())
	data.BodyText = state.StringPointerValue(termsAndConditions.GetBodyText())
	data.AcceptanceStatement = state.StringPointerValue(termsAndConditions.GetAcceptanceStatement())
	data.Version = state.Int32PtrToTypeInt32(termsAndConditions.GetVersion())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, termsAndConditions.GetRoleScopeTagIds())
	data.CreatedDateTime = state.TimeToString(termsAndConditions.GetCreatedDateTime())
	data.ModifiedDateTime = state.TimeToString(termsAndConditions.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
