package graphBetaApplications

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// initializeAndUploadContent handles content initialization and file upload
func (r *ApplicationsResource) initializeAndUploadContent(ctx context.Context, object *ApplicationsResourceModel, resp *resource.CreateResponse) error {
	appType := object.ApplicationType.ValueString()

	var contentVersionsBuilder interface{}
	var filesBuilder interface{}

	switch appType {
	case "AndroidLobApp":
		androidLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphAndroidLobApp()
		contentVersionsBuilder = androidLobAppBuilder.ContentVersions()
		filesBuilder = androidLobAppBuilder.ContentVersions()
	case "IosLobApp":
		iosLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphIosLobApp()
		contentVersionsBuilder = iosLobAppBuilder.ContentVersions()
		filesBuilder = iosLobAppBuilder.ContentVersions()
	case "MacOSDmgApp":
		macOSDmgAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSDmgApp()
		contentVersionsBuilder = macOSDmgAppBuilder.ContentVersions()
		filesBuilder = macOSDmgAppBuilder.ContentVersions()
	case "MacOSLobApp":
		macOSLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSLobApp()
		contentVersionsBuilder = macOSLobAppBuilder.ContentVersions()
		filesBuilder = macOSLobAppBuilder.ContentVersions()
	case "MacOSPkgApp":
		macOSPkgAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp()
		contentVersionsBuilder = macOSPkgAppBuilder.ContentVersions()
		filesBuilder = macOSPkgAppBuilder.ContentVersions()
	case "ManagedAndroidLobApp":
		managedAndroidLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedAndroidLobApp()
		contentVersionsBuilder = managedAndroidLobAppBuilder.ContentVersions()
		filesBuilder = managedAndroidLobAppBuilder.ContentVersions()
	case "ManagedIOSLobApp":
		managedIosLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedIOSLobApp()
		contentVersionsBuilder = managedIosLobAppBuilder.ContentVersions()
		filesBuilder = managedIosLobAppBuilder.ContentVersions()
	case "ManagedMobileLobApp":
		managedMobileLobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphManagedMobileLobApp()
		contentVersionsBuilder = managedMobileLobAppBuilder.ContentVersions()
		filesBuilder = managedMobileLobAppBuilder.ContentVersions()
	case "Win32LobApp":
		win32LobAppBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWin32LobApp()
		contentVersionsBuilder = win32LobAppBuilder.ContentVersions()
		filesBuilder = win32LobAppBuilder.ContentVersions()
	case "WindowsAppX":
		windowsAppXBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWindowsAppX()
		contentVersionsBuilder = windowsAppXBuilder.ContentVersions()
		filesBuilder = windowsAppXBuilder.ContentVersions()
	default:
		return fmt.Errorf("unsupported application type for content upload: %s", appType)
	}

	if contentVersionsBuilder == nil || filesBuilder == nil {
		return fmt.Errorf("failed to initialize request builders for app type %s", appType)
	}

	// Step 1: Initialize content version for upload
	contentVersion, err := contentVersionsBuilder.(interface {
		Post(ctx context.Context, body graphmodels.MobileAppContentable, config *RequestConfiguration) (graphmodels.MobileAppContentable, error)
	}).Post(ctx, graphmodels.NewMobileAppContent(), nil)
	if err != nil {
		return fmt.Errorf("failed to initialize content version for resource %s: %s", object.ID.ValueString(), err.Error())
	}

	// Step 2: Initialize file upload
	uploadUrl, err := filesBuilder.(interface {
		ByMobileAppContentId(contentId string) *MobileAppContentVersionsFilesRequestBuilder
	}).ByMobileAppContentId(*contentVersion.GetId()).
		Files().
		Post(ctx, graphmodels.NewMobileAppContentFile(), nil)
	if err != nil {
		return fmt.Errorf("failed to initialize file upload for content version %s: %s", *contentVersion.GetId(), err.Error())
	}

	// Step 3: Perform the actual file upload
	fileBytes, err := getFileBytes(object.FileName.ValueString())
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", object.FileName.ValueString(), err.Error())
	}

	err = uploadFile(uploadUrl, fileBytes)
	if err != nil {
		return fmt.Errorf("failed to upload file to %s: %s", uploadUrl, err.Error())
	}

	return nil
}

// Utility functions for file handling
func getFileBytes(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func uploadFile(uploadUrl string, fileBytes []byte) error {
	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	return nil
}
