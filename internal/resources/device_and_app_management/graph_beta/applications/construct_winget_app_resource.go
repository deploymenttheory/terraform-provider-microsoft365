package graphBetaApplications

import (
	"context"
	"fmt"
	"strings"

	utils "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/utilities"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructWinGetAppResource constructs a WinGetApp resource using data from the Terraform model.
// It fetches additional details from the Microsoft Store using FetchStoreAppDetails.
func constructWinGetAppResource(ctx context.Context, data *WinGetAppResourceModel, baseApp graphmodels.WinGetAppable) (graphmodels.WinGetAppable, error) {
	packageIdentifier := data.PackageIdentifier.ValueString()
	upperPackageIdentifier := strings.ToUpper(packageIdentifier)
	baseApp.SetPackageIdentifier(&upperPackageIdentifier)

	// Fetch metadata if auto-generate is enabled
	if data.AutomaticallyGenerateMetadata.ValueBool() {
		title, imageURL, description, publisher, err := FetchStoreAppDetails(ctx, packageIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch store details: %v", err)
		}

		if title == "" || description == "" || publisher == "" {
			return nil, fmt.Errorf("incomplete store details for packageIdentifier '%s'", packageIdentifier)
		}

		baseApp.SetDisplayName(&title)
		baseApp.SetDescription(&description)
		baseApp.SetPublisher(&publisher)

		if imageURL != "" {
			if iconBytes, err := utils.DownloadImage(imageURL); err == nil {
				largeIcon := graphmodels.NewMimeContent()
				iconType := "image/png"
				largeIcon.SetTypeEscaped(&iconType)
				largeIcon.SetValue(iconBytes)
				baseApp.SetLargeIcon(largeIcon)
			}
		}
	}

	baseApp.SetAdditionalData(map[string]interface{}{
		"repositoryType": "microsoftStore",
	})

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
			defaultRunAs := graphmodels.USER_RUNASACCOUNTTYPE
			installExperience.SetRunAsAccount(&defaultRunAs)
		}
		baseApp.SetInstallExperience(installExperience)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s WinGetApp resource values", ResourceName))

	return baseApp, nil
}
