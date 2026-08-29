package graphBetaWin32App

import (
	"context"
	"fmt"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sentinels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (r *Win32LobAppResource) publishContentVersion(ctx context.Context, appID string, installer *preparedWin32Content, response any, operation string) bool {
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		addContentDiagnostic(response, "Error publishing application content", fmt.Errorf("the content upload context does not have a deadline"))
		return false
	}

	contentBuilder := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appID).
		GraphWin32LobApp().
		ContentVersions()

	contentVersion, err := contentBuilder.Post(ctx, graphmodels.NewMobileAppContent(), nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}
	tflog.Debug(ctx, fmt.Sprintf("Content version created with ID: %s", *contentVersion.GetId()))

	createdFile, err := contentBuilder.
		ByMobileAppContentId(*contentVersion.GetId()).
		Files().
		Post(ctx, installer.contentFile(), nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}

	err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
		file, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			ByMobileAppContentFileId(*createdFile.GetId()).
			Get(ctx, nil)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("%w: %v", sentinels.ErrFileStatusFailed, err))
		}
		if file.GetUploadState() == nil {
			return retry.RetryableError(sentinels.ErrUploadStateNil)
		}

		state := *file.GetUploadState()
		if state == graphmodels.AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
			return nil
		}
		if state == graphmodels.AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
			return retry.NonRetryableError(sentinels.ErrAzureStorageURIRequestFailed)
		}
		return retry.RetryableError(fmt.Errorf("%w, current state: %s", sentinels.ErrWaitingForAzureStorageURI, state.String()))
	})
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}

	fileStatus, err := contentBuilder.
		ByMobileAppContentId(*contentVersion.GetId()).
		Files().
		ByMobileAppContentFileId(*createdFile.GetId()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}
	if fileStatus.GetAzureStorageUri() == nil {
		errors.HandleKiotaGraphError(ctx, sentinels.ErrAzureStorageURINil, response, operation, r.WritePermissions)
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("Uploading prepared encrypted package content: %s", installer.fileName))
	if err := construct.UploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), installer.encryptedFilePath); err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}

	err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
		commitBody, err := construct.CommitUploadedMobileAppWithEncryptionMetadata(installer.encryptionInfo)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("%w: %v", sentinels.ErrCommitRequestConstruction, err))
		}
		if err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			ByMobileAppContentFileId(*createdFile.GetId()).
			Commit().
			Post(ctx, commitBody, nil); err != nil {
			return retry.RetryableError(fmt.Errorf("%w: %v", sentinels.ErrFileCommitFailed, err))
		}
		return nil
	})
	if err != nil {
		addContentDiagnostic(response, "Error committing application content file", err)
		return false
	}

	if err := WaitForFileCommitCompletion(ctx, contentBuilder, *contentVersion.GetId(), *createdFile.GetId(), installer.encryptionInfo, response, r.WritePermissions); err != nil {
		addContentDiagnostic(response, "Error waiting for application content file commit", err)
		return false
	}

	updatePayload := graphmodels.NewWin32LobApp()
	updatePayload.SetCommittedContentVersion(contentVersion.GetId())
	if _, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appID).
		Patch(ctx, updatePayload, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, response, operation, r.WritePermissions)
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("Published content version %s for existing application %s", *contentVersion.GetId(), appID))
	return true
}

func addContentDiagnostic(response any, summary string, err error) {
	switch response := response.(type) {
	case *resource.CreateResponse:
		response.Diagnostics.AddError(summary, err.Error())
	case *resource.UpdateResponse:
		response.Diagnostics.AddError(summary, err.Error())
	}
}
