package graphBetaMacosPkgApp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *MacOSPkgAppResourceModel, remoteResource graphmodels.MacOSPkgAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = types.StringPointerValue(remoteResource.GetPublisher())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.IsFeatured = types.BoolPointerValue(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringPointerValue(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = types.StringPointerValue(remoteResource.GetInformationUrl())
	data.Owner = types.StringPointerValue(remoteResource.GetOwner())
	data.Developer = types.StringPointerValue(remoteResource.GetDeveloper())
	data.Notes = types.StringPointerValue(remoteResource.GetNotes())
	data.UploadState = state.Int32PtrToTypeInt64(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = types.BoolPointerValue(remoteResource.GetIsAssigned())

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = sharedmodels.MimeContentResourceModel{
			Type:  types.StringPointerValue(largeIcon.GetTypeEscaped()),
			Value: types.StringValue(state.ByteToString(largeIcon.GetValue())),
		}
	}

	var roleScopeTagIds []attr.Value
	for _, v := range state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds()) {
		roleScopeTagIds = append(roleScopeTagIds, v)
	}

	data.RoleScopeTagIds = types.ListValueMust(
		types.StringType,
		roleScopeTagIds,
	)

	data.DependentAppCount = state.Int32PtrToTypeInt64(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersededAppCount())
	data.CommittedContentVersion = types.StringPointerValue(remoteResource.GetCommittedContentVersion())
	data.FileName = types.StringPointerValue(remoteResource.GetFileName())
	data.Size = state.Int64PtrToTypeInt64(remoteResource.GetSize())
	data.PrimaryBundleId = types.StringPointerValue(remoteResource.GetPrimaryBundleId())
	data.PrimaryBundleVersion = types.StringPointerValue(remoteResource.GetPrimaryBundleVersion())
	data.IgnoreVersionDetection = types.BoolPointerValue(remoteResource.GetIgnoreVersionDetection())

	// Handle IncludedApps
	includedApps := remoteResource.GetIncludedApps()
	if len(includedApps) == 0 {
		data.IncludedApps = []MacOSIncludedAppResourceModel{}
	} else {
		data.IncludedApps = make([]MacOSIncludedAppResourceModel, len(includedApps))
		for i, app := range includedApps {
			data.IncludedApps[i] = MacOSIncludedAppResourceModel{
				BundleId:      types.StringPointerValue(app.GetBundleId()),
				BundleVersion: types.StringPointerValue(app.GetBundleVersion()),
			}
		}
	}

	// Handle MinimumSupportedOperatingSystem
	minOS := remoteResource.GetMinimumSupportedOperatingSystem()
	data.MinimumSupportedOperatingSystem = MacOSMinimumOperatingSystemResourceModel{
		V10_7:  types.BoolPointerValue(minOS.GetV107()),
		V10_8:  types.BoolPointerValue(minOS.GetV108()),
		V10_9:  types.BoolPointerValue(minOS.GetV109()),
		V10_10: types.BoolPointerValue(minOS.GetV1010()),
		V10_11: types.BoolPointerValue(minOS.GetV1011()),
		V10_12: types.BoolPointerValue(minOS.GetV1012()),
		V10_13: types.BoolPointerValue(minOS.GetV1013()),
		V10_14: types.BoolPointerValue(minOS.GetV1014()),
		V10_15: types.BoolPointerValue(minOS.GetV1015()),
		V11_0:  types.BoolPointerValue(minOS.GetV110()),
		V12_0:  types.BoolPointerValue(minOS.GetV120()),
		V13_0:  types.BoolPointerValue(minOS.GetV130()),
		V14_0:  types.BoolPointerValue(minOS.GetV140()),
	}

	// Handle PreInstallScript / PostInstallScript
	data.PreInstallScript = MacOSAppScriptResourceModel{
		ScriptContent: types.StringPointerValue(remoteResource.GetPreInstallScript().GetScriptContent()),
	}

	data.PostInstallScript = MacOSAppScriptResourceModel{
		ScriptContent: types.StringPointerValue(remoteResource.GetPostInstallScript().GetScriptContent()),
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
