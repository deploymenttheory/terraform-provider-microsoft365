package graphBetaWinGetApp

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	utils "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/utilities"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a WinGetApp resource using data from the Terraform model.
// It fetches additional details from the Microsoft Store using FetchStoreAppDetails.
func constructResource(ctx context.Context, data *WinGetAppResourceModel) (graphmodels.WinGetAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWinGetApp()

	packageIdentifier := data.PackageIdentifier.ValueString()
	upperPackageIdentifier := strings.ToUpper(packageIdentifier)
	requestBody.SetPackageIdentifier(&upperPackageIdentifier)

	// Fetch metadata from the Microsoft Store using the packageIdentifier if AutomaticallyGenerateMetadata is true
	if data.AutomaticallyGenerateMetadata.ValueBool() {
		title, imageURL, description, publisher, err := FetchStoreAppDetails(ctx, packageIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch store details for packageIdentifier '%s': %v", packageIdentifier, err)
		}

		if title == "" || description == "" || publisher == "" {
			return nil, fmt.Errorf("incomplete store details for packageIdentifier '%s'. Missing required fields: Title='%s', Description='%s', Publisher='%s'", packageIdentifier, title, description, publisher)
		}

		requestBody.SetDisplayName(&title)
		requestBody.SetDescription(&description)
		requestBody.SetPublisher(&publisher)

		if imageURL != "" {
			iconBytes, err := utils.DownloadImage(imageURL)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to download icon image from URL '%s': %v", imageURL, err))
			} else {
				largeIcon := graphmodels.NewMimeContent()
				iconType := "image/png"
				largeIcon.SetTypeEscaped(&iconType)
				largeIcon.SetValue(iconBytes)
				requestBody.SetLargeIcon(largeIcon)
			}
		}
	} else {
		// Use the provided values from the model
		constructors.SetStringProperty(data.Description, requestBody.SetDescription)
		constructors.SetStringProperty(data.Publisher, requestBody.SetPublisher)
		constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)

		if !data.LargeIcon.IsNull() {
			largeIcon := graphmodels.NewMimeContent()
			var iconData map[string]attr.Value
			data.LargeIcon.As(context.Background(), &iconData, basetypes.ObjectAsOptions{})

			iconType := "image/png"
			largeIcon.SetTypeEscaped(&iconType)

			if valueVal, ok := iconData["value"].(types.String); ok {
				iconBytes, err := base64.StdEncoding.DecodeString(valueVal.ValueString())
				if err != nil {
					return nil, fmt.Errorf("failed to decode icon base64: %v", err)
				}
				largeIcon.SetValue(iconBytes)
			}
			requestBody.SetLargeIcon(largeIcon)
		}
	}

	constructors.SetBoolProperty(data.IsFeatured, requestBody.SetIsFeatured)
	constructors.SetStringProperty(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	constructors.SetStringProperty(data.InformationUrl, requestBody.SetInformationUrl)
	constructors.SetStringProperty(data.Owner, requestBody.SetOwner)
	constructors.SetStringProperty(data.Developer, requestBody.SetDeveloper)
	constructors.SetStringProperty(data.Notes, requestBody.SetNotes)
	constructors.SetStringProperty(data.ManifestHash, requestBody.SetManifestHash)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	additionalData := map[string]interface{}{
		"repositoryType": "microsoftStore",
	}
	requestBody.SetAdditionalData(additionalData)

	if data.InstallExperience != nil && !data.InstallExperience.RunAsAccount.IsNull() {
		installExperience := graphmodels.NewWinGetAppInstallExperience()
		runAsAccount := data.InstallExperience.RunAsAccount.ValueString()
		switch runAsAccount {
		case "system":
			systemAccount := graphmodels.SYSTEM_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&systemAccount)
		case "user":
			userAccount := graphmodels.USER_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&userAccount)
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown runAsAccount value '%s'. Defaulting to 'user'.", runAsAccount))
			defaultRunAs := graphmodels.USER_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&defaultRunAs)
		}
		requestBody.SetInstallExperience(installExperience)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
