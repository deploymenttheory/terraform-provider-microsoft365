package graphBetaApplications

import (
	"context"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func initializeContentIfNeeded(ctx context.Context, r *ApplicationsResource, object *ApplicationsResourceModel, resp *resource.CreateResponse) error {
	var fileSourcePath *types.String

	switch object.ApplicationType.ValueString() {
	case "MacOSPkgApp":
		if object.MacOSPkgApp != nil {
			fileSourcePath = &object.MacOSPkgApp.PackageInstallerFileSource
		}
	}

	if fileSourcePath != nil && !fileSourcePath.IsNull() {
		contentVersion, err := r.initializeContentUpload(ctx, object)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return err
		}

		_, err = r.contentFileUpload(ctx, object, contentVersion)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return err
		}
	}

	return nil
}

// initializeContentUpload handles content initialization of application upload
func (r *ApplicationsResource) initializeContentUpload(ctx context.Context, object *ApplicationsResourceModel) (graphmodels.MobileAppContentable, error) {
	appType := object.ApplicationType.ValueString()

	contentVersionsBuilder, err := r.NewContentVersionsBuilder(appType, object.ID.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to get content versions builder: %w", err)
	}

	contentVersion, err := contentVersionsBuilder.(interface {
		Post(ctx context.Context, body graphmodels.MobileAppContentable, config *ApplicationsResourceModel) (graphmodels.MobileAppContentable, error)
	}).Post(ctx, graphmodels.NewMobileAppContent(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize content version for resource %s: %s", object.ID.ValueString(), err.Error())
	}

	return contentVersion, nil
}

// contentFileUpload handles the file upload for the application content
func (r *ApplicationsResource) contentFileUpload(ctx context.Context, object *ApplicationsResourceModel, contentVersion graphmodels.MobileAppContentable) (string, error) {
	appType := object.ApplicationType.ValueString()

	filesBuilder, err := r.NewContentVersionsBuilder(appType, object.ID.ValueString())
	if err != nil {
		return "", fmt.Errorf("failed to get files builder: %w", err)
	}

	uploadUrl, err := filesBuilder.(interface {
		ByMobileAppContentId(contentId string) interface {
			Files() interface {
				Post(ctx context.Context, body graphmodels.MobileAppContentFileable, config *ApplicationsResourceModel) (string, error)
			}
		}
	}).ByMobileAppContentId(*contentVersion.GetId()).Files().Post(ctx, graphmodels.NewMobileAppContentFile(), nil)

	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}

// NewContentVersionsBuilder returns a new content versions builder based on the application type
func (r *ApplicationsResource) NewContentVersionsBuilder(appType, appID string) (interface{}, error) {
	switch appType {
	case "AndroidLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphAndroidLobApp().ContentVersions(), nil
	case "IosLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphIosLobApp().ContentVersions(), nil
	case "MacOSDmgApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSDmgApp().ContentVersions(), nil
	case "MacOSLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSLobApp().ContentVersions(), nil
	case "MacOSPkgApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSPkgApp().ContentVersions(), nil
	case "ManagedAndroidLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedAndroidLobApp().ContentVersions(), nil
	case "ManagedIOSLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedIOSLobApp().ContentVersions(), nil
	case "ManagedMobileLobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedMobileLobApp().ContentVersions(), nil
	case "Win32LobApp":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphWin32LobApp().ContentVersions(), nil
	case "WindowsAppX":
		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphWindowsAppX().ContentVersions(), nil
	default:
		return nil, fmt.Errorf("unsupported application type: %s", appType)
	}
}

func getFileBytes(object *ApplicationsResourceModel) ([]byte, error) {
	var filePath string

	switch object.ApplicationType.ValueString() {
	case "MacOSPkgApp":
		if object.MacOSPkgApp == nil || object.MacOSPkgApp.PackageInstallerFileSource.IsNull() {
			return nil, fmt.Errorf("package_installer_file_source is required for MacOSPkgApp")
		}
		filePath = object.MacOSPkgApp.PackageInstallerFileSource.ValueString()
	// TODO: Add other app types here
	default:
		return nil, fmt.Errorf("unsupported application type for file upload: %s", object.ApplicationType.ValueString())
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", filePath, err.Error())
	}
	return fileBytes, nil
}
