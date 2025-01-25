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

	var contentVersionsBuilder interface{}

	switch appType {
	case "AndroidLobApp":
		androidLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphAndroidLobApp()
		contentVersionsBuilder = androidLobAppBuilder.ContentVersions()
	case "IosLobApp":
		iosLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphIosLobApp()
		contentVersionsBuilder = iosLobAppBuilder.ContentVersions()
	case "MacOSDmgApp":
		macOSDmgAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSDmgApp()
		contentVersionsBuilder = macOSDmgAppBuilder.ContentVersions()
	case "MacOSLobApp":
		macOSLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSLobApp()
		contentVersionsBuilder = macOSLobAppBuilder.ContentVersions()
	case "MacOSPkgApp":
		macOSPkgAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp()
		contentVersionsBuilder = macOSPkgAppBuilder.ContentVersions()
	case "ManagedAndroidLobApp":
		managedAndroidLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedAndroidLobApp()
		contentVersionsBuilder = managedAndroidLobAppBuilder.ContentVersions()
	case "ManagedIOSLobApp":
		managedIosLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedIOSLobApp()
		contentVersionsBuilder = managedIosLobAppBuilder.ContentVersions()
	case "ManagedMobileLobApp":
		managedMobileLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedMobileLobApp()
		contentVersionsBuilder = managedMobileLobAppBuilder.ContentVersions()
	case "Win32LobApp":
		win32LobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWin32LobApp()
		contentVersionsBuilder = win32LobAppBuilder.ContentVersions()
	case "WindowsAppX":
		windowsAppXBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWindowsAppX()
		contentVersionsBuilder = windowsAppXBuilder.ContentVersions()
	default:
		return nil, fmt.Errorf("unsupported application type for content upload: %s", appType)
	}

	contentVersion, err := contentVersionsBuilder.(interface {
		Post(ctx context.Context, body graphmodels.MobileAppContentable, config *ApplicationsResourceModel) (graphmodels.MobileAppContentable, error)
	}).Post(ctx, graphmodels.NewMobileAppContent(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize content version for resource %s: %s", object.ID.ValueString(), err.Error())
	}

	return contentVersion, nil
}

func (r *ApplicationsResource) contentFileUpload(ctx context.Context, object *ApplicationsResourceModel, contentVersion graphmodels.MobileAppContentable) (string, error) {
	appType := object.ApplicationType.ValueString()
	var filesBuilder interface{}

	switch appType {
	case "AndroidLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphAndroidLobApp().
			ContentVersions()
	case "IosLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphIosLobApp().
			ContentVersions()
	case "MacOSDmgApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSDmgApp().
			ContentVersions()
	case "MacOSLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSLobApp().
			ContentVersions()
	case "MacOSPkgApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp().
			ContentVersions()
	case "ManagedAndroidLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedAndroidLobApp().
			ContentVersions()
	case "ManagedIOSLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedIOSLobApp().
			ContentVersions()
	case "ManagedMobileLobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedMobileLobApp().
			ContentVersions()
	case "Win32LobApp":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWin32LobApp().
			ContentVersions()
	case "WindowsAppX":
		filesBuilder = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWindowsAppX().
			ContentVersions()
	default:
		return "", fmt.Errorf("unsupported application type for content upload: %s", appType)
	}

	uploadUrl, err := filesBuilder.(interface {
		ByMobileAppContentId(contentId string) interface {
			Files() interface {
				Post(ctx context.Context, body graphmodels.MobileAppContentFileable, config *ApplicationsResourceModel) (string, error)
			}
		}
	}).ByMobileAppContentId(*contentVersion.GetId()).
		Files().
		Post(ctx, graphmodels.NewMobileAppContentFile(), nil)
	if err != nil {
		return "", err
	}

	contentFile := graphmodels.NewMobileAppContentFile()
	fileBytes, err := getFileBytes(object)
	if err != nil {
		return "", err
	}

	fileSize := int64(len(fileBytes))
	contentFile.SetSizeInBytes(&fileSize)

	return uploadUrl, nil
}

func getFileBytes(object *ApplicationsResourceModel) ([]byte, error) {
	var filePath string

	switch object.ApplicationType.ValueString() {
	case "MacOSPkgApp":
		if object.MacOSPkgApp == nil || object.MacOSPkgApp.PackageInstallerFileSource.IsNull() {
			return nil, fmt.Errorf("package_installer_file_source is required for MacOSPkgApp")
		}
		filePath = object.MacOSPkgApp.PackageInstallerFileSource.ValueString()
	// Add other app types here
	default:
		return nil, fmt.Errorf("unsupported application type for file upload: %s", object.ApplicationType.ValueString())
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", filePath, err.Error())
	}
	return fileBytes, nil
}
