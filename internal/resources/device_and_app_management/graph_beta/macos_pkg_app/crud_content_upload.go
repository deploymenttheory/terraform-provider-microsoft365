package graphBetaMacOSPKGApp

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
// 	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
// )

// func initializeContentIfNeeded(ctx context.Context, r *ApplicationsResource, object *ApplicationsResourceModel, resp *resource.CreateResponse) error {
// 	var fileSourcePath *types.String

// 	switch object.ApplicationType.ValueString() {
// 	case "MacOSPkgApp":
// 		if object.MacOSPkgApp != nil {
// 			fileSourcePath = &object.MacOSPkgApp.PackageInstallerFileSource
// 		}
// 	}

// 	if fileSourcePath != nil && !fileSourcePath.IsNull() {
// 		contentVersion, err := r.initializeContentUpload(ctx, object)
// 		if err != nil {
// 			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
// 			return err
// 		}

// 		_, err = r.contentFileUpload(ctx, object, contentVersion)
// 		if err != nil {
// 			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (r *ApplicationsResource) initializeContentUpload(ctx context.Context, object *ApplicationsResourceModel) (graphmodels.MobileAppContentable, error) {
// 	builder, err := r.NewContentVersionsBuilder(object.ApplicationType.ValueString(), object.ID.ValueString())
// 	if err != nil {
// 		return nil, err
// 	}

// 	content := graphmodels.NewMobileAppContent()

// 	switch object.ApplicationType.ValueString() {
// 	case "AndroidLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphAndroidLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "IosLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphIosLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "MacOSDmgApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSDmgAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "MacOSLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "MacOSPkgApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSPkgAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "ManagedAndroidLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedAndroidLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "ManagedIOSLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedIOSLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "ManagedMobileLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedMobileLobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "Win32LobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphWin32LobAppContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	case "WindowsAppX":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphWindowsAppXContentVersionsRequestBuilder).Post(ctx, content, nil)
// 	default:
// 		return nil, fmt.Errorf("unsupported application type: %s", object.ApplicationType.ValueString())
// 	}
// }

// // contentFileUpload handles the file upload for the application content
// func (r *ApplicationsResource) contentFileUpload(ctx context.Context, object *ApplicationsResourceModel, contentVersion graphmodels.MobileAppContentable) (graphmodels.MobileAppContentFileable, error) {
// 	builder, err := r.NewContentVersionsBuilder(object.ApplicationType.ValueString(), object.ID.ValueString())
// 	if err != nil {
// 		return nil, err
// 	}

// 	contentFile := graphmodels.NewMobileAppContentFile()
// 	fileName := filepath.Base(object.MacOSPkgApp.PackageInstallerFileSource.ValueString())
// 	contentFile.SetName(&fileName)

// 	fileInfo, err := os.Stat(object.MacOSPkgApp.PackageInstallerFileSource.ValueString())
// 	if err != nil {
// 		return nil, err
// 	}
// 	size := fileInfo.Size()
// 	contentFile.SetSize(&size)

// 	switch object.ApplicationType.ValueString() {
// 	case "AndroidLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphAndroidLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "IosLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphIosLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "MacOSDmgApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSDmgAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "MacOSLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "MacOSPkgApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphMacOSPkgAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "ManagedAndroidLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedAndroidLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "ManagedIOSLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedIOSLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "ManagedMobileLobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphManagedMobileLobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "Win32LobApp":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphWin32LobAppContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	case "WindowsAppX":
// 		return builder.(*deviceappmanagement.MobileAppsItemGraphWindowsAppXContentVersionsRequestBuilder).
// 			ByMobileAppContentId(*contentVersion.GetId()).
// 			Files().
// 			Post(ctx, contentFile, nil)
// 	default:
// 		return nil, fmt.Errorf("unsupported application type: %s", object.ApplicationType.ValueString())
// 	}
// }

// // NewContentVersionsBuilder returns a new content versions builder based on the application type
// func (r *ApplicationsResource) NewContentVersionsBuilder(appType, appID string) (interface{}, error) {
// 	switch appType {
// 	case "AndroidLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphAndroidLobApp().ContentVersions(), nil
// 	case "IosLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphIosLobApp().ContentVersions(), nil
// 	case "MacOSDmgApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSDmgApp().ContentVersions(), nil
// 	case "MacOSLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSLobApp().ContentVersions(), nil
// 	case "MacOSPkgApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphMacOSPkgApp().ContentVersions(), nil
// 	case "ManagedAndroidLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedAndroidLobApp().ContentVersions(), nil
// 	case "ManagedIOSLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedIOSLobApp().ContentVersions(), nil
// 	case "ManagedMobileLobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphManagedMobileLobApp().ContentVersions(), nil
// 	case "Win32LobApp":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphWin32LobApp().ContentVersions(), nil
// 	case "WindowsAppX":
// 		return r.client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).GraphWindowsAppX().ContentVersions(), nil
// 	default:
// 		return nil, fmt.Errorf("unsupported application type: %s", appType)
// 	}
// }
