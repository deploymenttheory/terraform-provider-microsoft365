package graphBetaApplications

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *ApplicationsResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var baseApp graphmodels.MobileAppable
	var err error

	switch data.ApplicationType.ValueString() {
	case "WindowsStoreApp":
		if data.WinGetApp == nil {
			return nil, fmt.Errorf("winget_app configuration block required for WinGetApp type")
		}
		baseApp = graphmodels.NewWinGetApp()
		if err = constructBaseProperties(ctx, baseApp, data); err != nil {
			return nil, err
		}
		requestBody, err := constructWinGetAppResource(ctx, data.WinGetApp, baseApp.(graphmodels.WinGetAppable))
		if err != nil {
			return nil, err
		}

		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return requestBody, nil

	case "MacOSPkgApp":
		if data.MacOSPkgApp == nil {
			return nil, fmt.Errorf("macos_pkg_app configuration block required for MacOSPkgApp type")
		}
		baseApp = graphmodels.NewMacOSPkgApp()
		if err = constructBaseProperties(ctx, baseApp, data); err != nil {
			return nil, err
		}
		requestBody, err := constructMacOSPkgAppResource(ctx, data.MacOSPkgApp, baseApp.(graphmodels.MacOSPkgAppable))
		if err != nil {
			return nil, err
		}

		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return requestBody, nil

	default:
		return nil, fmt.Errorf("unsupported intune application type: %s", data.ApplicationType.ValueString())
	}
}

// constructBaseProperties sets the base properties of the MobileAppable object
func constructBaseProperties(ctx context.Context, baseApp graphmodels.MobileAppable, data *ApplicationsResourceModel) error {
	constructors.SetStringProperty(data.Description, baseApp.SetDescription)
	constructors.SetStringProperty(data.Publisher, baseApp.SetPublisher)
	constructors.SetStringProperty(data.DisplayName, baseApp.SetDisplayName)
	constructors.SetStringProperty(data.InformationUrl, baseApp.SetInformationUrl)
	constructors.SetBoolProperty(data.IsFeatured, baseApp.SetIsFeatured)
	constructors.SetStringProperty(data.Owner, baseApp.SetOwner)
	constructors.SetStringProperty(data.Developer, baseApp.SetDeveloper)
	constructors.SetStringProperty(data.Notes, baseApp.SetNotes)

	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if len(data.Categories) > 0 {
		categories := make([]graphmodels.MobileAppCategoryable, 0, len(data.Categories))
		for _, category := range data.Categories {
			mobileAppCategory := graphmodels.NewMobileAppCategory()
			constructors.SetStringProperty(category.ID, mobileAppCategory.SetId)
			constructors.SetStringProperty(category.DisplayName, mobileAppCategory.SetDisplayName)
			categories = append(categories, mobileAppCategory)
		}
		baseApp.SetCategories(categories)
	}

	if !data.LargeIcon.IsNull() {
		largeIcon := graphmodels.NewMimeContent()
		var iconData map[string]attr.Value
		data.LargeIcon.As(ctx, &iconData, basetypes.ObjectAsOptions{})

		iconType := "image/png"
		largeIcon.SetTypeEscaped(&iconType)

		if valueVal, ok := iconData["value"].(types.String); ok {
			iconBytes, err := base64.StdEncoding.DecodeString(valueVal.ValueString())
			if err != nil {
				return fmt.Errorf("failed to decode icon base64: %v", err)
			}
			largeIcon.SetValue(iconBytes)
		}
		baseApp.SetLargeIcon(largeIcon)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s base resource values ", ResourceName))
	return nil
}
