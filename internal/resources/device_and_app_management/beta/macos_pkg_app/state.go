package graphbetamacospkgapp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.Publisher = types.StringValue(state.StringPtrToString(remoteResource.GetPublisher()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.IsFeatured = state.BoolPtrToTypeBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetPrivacyInformationUrl()))
	data.InformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetInformationUrl()))
	data.Owner = types.StringValue(state.StringPtrToString(remoteResource.GetOwner()))
	data.Developer = types.StringValue(state.StringPtrToString(remoteResource.GetDeveloper()))
	data.Notes = types.StringValue(state.StringPtrToString(remoteResource.GetNotes()))
	data.UploadState = state.Int32PtrToTypeInt64(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = sharedmodels.MimeContentResourceModel{
			Type:  types.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
			Value: types.StringValue(state.ByteToString(largeIcon.GetValue())),
		}
	}

	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	data.DependentAppCount = state.Int32PtrToTypeInt64(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersededAppCount())
	data.CommittedContentVersion = types.StringValue(state.StringPtrToString(remoteResource.GetCommittedContentVersion()))
	data.FileName = types.StringValue(state.StringPtrToString(remoteResource.GetFileName()))
	data.Size = state.Int64PtrToTypeInt64(remoteResource.GetSize())
	data.PrimaryBundleId = types.StringValue(state.StringPtrToString(remoteResource.GetPrimaryBundleId()))
	data.PrimaryBundleVersion = types.StringValue(state.StringPtrToString(remoteResource.GetPrimaryBundleVersion()))
	data.IgnoreVersionDetection = state.BoolPtrToTypeBool(remoteResource.GetIgnoreVersionDetection())

	// Handle IncludedApps
	includedApps := remoteResource.GetIncludedApps()
	if len(includedApps) == 0 {
		data.IncludedApps = []MacOSIncludedAppResourceModel{}
	} else {
		data.IncludedApps = make([]MacOSIncludedAppResourceModel, len(includedApps))
		for i, app := range includedApps {
			data.IncludedApps[i] = MacOSIncludedAppResourceModel{
				BundleId:      types.StringValue(state.StringPtrToString(app.GetBundleId())),
				BundleVersion: types.StringValue(state.StringPtrToString(app.GetBundleVersion())),
			}
		}
	}

	// Handle MinimumSupportedOperatingSystem
	minOS := remoteResource.GetMinimumSupportedOperatingSystem()
	data.MinimumSupportedOperatingSystem = MacOSMinimumOperatingSystemResourceModel{
		V10_7:  state.BoolPtrToTypeBool(minOS.GetV107()),
		V10_8:  state.BoolPtrToTypeBool(minOS.GetV108()),
		V10_9:  state.BoolPtrToTypeBool(minOS.GetV109()),
		V10_10: state.BoolPtrToTypeBool(minOS.GetV1010()),
		V10_11: state.BoolPtrToTypeBool(minOS.GetV1011()),
		V10_12: state.BoolPtrToTypeBool(minOS.GetV1012()),
		V10_13: state.BoolPtrToTypeBool(minOS.GetV1013()),
		V10_14: state.BoolPtrToTypeBool(minOS.GetV1014()),
		V10_15: state.BoolPtrToTypeBool(minOS.GetV1015()),
		V11_0:  state.BoolPtrToTypeBool(minOS.GetV110()),
		V12_0:  state.BoolPtrToTypeBool(minOS.GetV120()),
		V13_0:  state.BoolPtrToTypeBool(minOS.GetV130()),
		V14_0:  state.BoolPtrToTypeBool(minOS.GetV140()),
	}

	// Handle PreInstallScript / PostInstallScript
	data.PreInstallScript = MacOSAppScriptResourceModel{
		ScriptContent: types.StringValue(state.StringPtrToString(remoteResource.GetPreInstallScript().GetScriptContent())),
	}

	data.PostInstallScript = MacOSAppScriptResourceModel{
		ScriptContent: types.StringValue(state.StringPtrToString(remoteResource.GetPostInstallScript().GetScriptContent())),
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
