package graphBetaWinGetApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	utils "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/utilities"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a WinGetApp resource using data from the Terraform model.
// It fetches additional details from the Microsoft Store using FetchStoreAppDetails.
func constructResource(ctx context.Context, data *WinGetAppResourceModel) (models.WinGetAppable, error) {
	construct.DebugPrintStruct(ctx, "Constructing WinGet App resource from model", data)

	requestBody := models.NewWinGetApp()

	// Set the packageIdentifier
	packageIdentifier := data.PackageIdentifier.ValueString()
	upperPackageIdentifier := utils.ToUpperCase(packageIdentifier)
	requestBody.SetPackageIdentifier(&upperPackageIdentifier)

	// Fetch metadata from the Microsoft Store using the packageIdentifier
	title, imageURL, description, publisher, err := FetchStoreAppDetails(packageIdentifier)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to fetch store details for packageIdentifier '%s': %v", packageIdentifier, err))
		return nil, fmt.Errorf("failed to fetch store details for packageIdentifier '%s': %v", packageIdentifier, err)
	}

	// Check if any required value (title, description, publisher) is empty
	if title == "" || description == "" || publisher == "" {
		errMsg := fmt.Sprintf("Incomplete store details for packageIdentifier '%s'. Missing required fields: Title='%s', Description='%s', Publisher='%s'", packageIdentifier, title, description, publisher)
		tflog.Error(ctx, errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Set the fetched title, description, and publisher
	requestBody.SetDisplayName(&title)
	requestBody.SetDescription(&description)
	requestBody.SetPublisher(&publisher)

	// Attempt to download and set the large icon from the msft app store
	if imageURL != "" {
		iconBytes, err := utils.DownloadImage(imageURL)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to download icon image from URL '%s': %v. Continuing without setting the icon.", imageURL, err))
		} else {
			largeIcon := models.NewMimeContent()
			iconType := "image/png"
			largeIcon.SetTypeEscaped(&iconType)
			largeIcon.SetValue(iconBytes)
			requestBody.SetLargeIcon(largeIcon)
			tflog.Debug(ctx, fmt.Sprintf("Icon set from store URL. Data length: %d bytes", len(iconBytes)))
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("No icon image URL found for packageIdentifier '%s'. The large icon will not be set.", packageIdentifier))
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

	// Set role scope tag IDs if provided
	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, len(data.RoleScopeTagIds))
		for i, id := range data.RoleScopeTagIds {
			roleScopeTagIds[i] = id.ValueString()
		}
		requestBody.SetRoleScopeTagIds(roleScopeTagIds)
	}

	// Set additional data
	additionalData := map[string]interface{}{
		"repositoryType": "microsoftStore",
	}
	requestBody.SetAdditionalData(additionalData)

	// Set the install experience
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
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown runAsAccount value '%s'. Defaulting to 'user'.", runAsAccount))
			defaultRunAs := models.USER_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&defaultRunAs)
		}
		requestBody.SetInstallExperience(installExperience)
	}

	return requestBody, nil
}
