package graphBetaWinGetApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *WinGetAppResourceModel) (models.WinGetAppable, error) {
	construct.DebugPrintStruct(ctx, "Constructing WinGet App resource from model", data)

	requestBody := models.NewWinGetApp()

	// Set required fields
	requestBody.SetDisplayName(data.DisplayName.ValueStringPointer())
	requestBody.SetPackageIdentifier(data.PackageIdentifier.ValueStringPointer())

	// Set optional fields
	if !data.Description.IsNull() {
		requestBody.SetDescription(data.Description.ValueStringPointer())
	}
	if !data.Publisher.IsNull() {
		requestBody.SetPublisher(data.Publisher.ValueStringPointer())
	}
	if !data.IsFeatured.IsNull() {
		requestBody.SetIsFeatured(data.IsFeatured.ValueBoolPointer())
	}
	if !data.PrivacyInformationUrl.IsNull() {
		requestBody.SetPrivacyInformationUrl(data.PrivacyInformationUrl.ValueStringPointer())
	}
	if !data.InformationUrl.IsNull() {
		requestBody.SetInformationUrl(data.InformationUrl.ValueStringPointer())
	}
	if !data.Owner.IsNull() {
		requestBody.SetOwner(data.Owner.ValueStringPointer())
	}
	if !data.Developer.IsNull() {
		requestBody.SetDeveloper(data.Developer.ValueStringPointer())
	}
	if !data.Notes.IsNull() {
		requestBody.SetNotes(data.Notes.ValueStringPointer())
	}
	if !data.ManifestHash.IsNull() {
		requestBody.SetManifestHash(data.ManifestHash.ValueStringPointer())
	}

	// Set role scope tag IDs
	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, len(data.RoleScopeTagIds))
		for i, id := range data.RoleScopeTagIds {
			roleScopeTagIds[i] = id.ValueString()
		}
		requestBody.SetRoleScopeTagIds(roleScopeTagIds)
	}

	// Set large icon
	if data.LargeIcon != nil && !data.LargeIcon.Type.IsNull() && !data.LargeIcon.Value.IsNull() {
		largeIcon := models.NewMimeContent()
		largeIcon.SetTypeEscaped(data.LargeIcon.Type.ValueStringPointer())
		largeIcon.SetValue([]byte(data.LargeIcon.Value.ValueString()))
		requestBody.SetLargeIcon(largeIcon)
	}

	// Set install experience
	if data.InstallExperience != nil && !data.InstallExperience.RunAsAccount.IsNull() {
		installExperience := models.NewWinGetAppInstallExperience()
		runAsAccount := data.InstallExperience.RunAsAccount.ValueString()
		switch runAsAccount {
		case "system":
			systemAccount := models.SYSTEM_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&systemAccount)
		case "user":
			userAccount := models.USER_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&userAccount)
		}
		requestBody.SetInstallExperience(installExperience)
	}

	return requestBody, nil
}
