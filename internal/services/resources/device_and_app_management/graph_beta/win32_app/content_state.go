package graphBetaWin32App

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// File commit completion constants
const (
	MaxCommitRetries  = 20
	InitialBackoff    = 1 * time.Second
	MaxBackoff        = 30 * time.Second
	BackoffMultiplier = 1.5
)

// WaitForFileCommitCompletion waits for a mobile app content file commit to complete with exponential backoff
// It intelligently handles different upload states and provides detailed logging
func WaitForFileCommitCompletion(
	ctx context.Context,
	contentBuilder *deviceappmanagement.MobileAppsItemGraphWin32LobAppContentVersionsRequestBuilder,
	contentVersionID string,
	fileID string,
	encryptionInfo *construct.EncryptionInfo,
	resp *resource.CreateResponse,
	permissions []string,
) error {
	tflog.Debug(ctx, "Starting file commit completion check with exponential backoff")

	backoff := InitialBackoff

	for i := 0; i < MaxCommitRetries; i++ {
		tflog.Debug(ctx, fmt.Sprintf("Commit status check attempt %d/%d (backoff: %v)", i+1, MaxCommitRetries, backoff))

		file, err := contentBuilder.
			ByMobileAppContentId(contentVersionID).
			Files().
			ByMobileAppContentFileId(fileID).
			Get(ctx, nil)

		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Error retrieving file status: %v", err))
			errors.HandleGraphError(ctx, err, resp, "WaitForFileCommitCompletion", permissions)
			return err
		}

		if file.GetUploadState() == nil {
			tflog.Debug(ctx, "File upload state is nil, retrying...")
			time.Sleep(calculateBackoffWithJitter(backoff))
			backoff = incrementBackoff(backoff)
			continue
		}

		state := *file.GetUploadState()
		tflog.Debug(ctx, fmt.Sprintf("Current file state: %s", state.String()))

		// Handle different states
		switch state {
		case graphmodels.COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "✅ File commit completed successfully")
			return nil

		case graphmodels.COMMITFILEPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "⏳ File commit is pending, waiting...")

		case graphmodels.COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "❌ File commit failed; attempting to retry commit request")
			err = retryCommitRequest(ctx, contentBuilder, contentVersionID, fileID, encryptionInfo)
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Error during commit retry: %v", err))
				// Continue the loop rather than returning to give it more chances
			}

		case graphmodels.COMMITFILETIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "⏰ File commit timed out; attempting to retry commit request")
			err = retryCommitRequest(ctx, contentBuilder, contentVersionID, fileID, encryptionInfo)
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Error during commit retry: %v", err))
			}

		case graphmodels.SUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "✅ Upload succeeded, but commit status not yet updated, continuing to wait...")

		case graphmodels.TRANSIENTERROR_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "⚠️ Transient error detected, retrying...")

		case graphmodels.ERROR_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "❌ Permanent error detected in upload state")
			return fmt.Errorf("permanent error detected in upload state")

		case graphmodels.UNKNOWN_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "❓ Unknown upload state detected, continuing to retry...")

		// Azure Storage URI states
		case graphmodels.AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE,
			graphmodels.AZURESTORAGEURIRENEWALSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "Azure storage URI is valid, but commit not yet started, attempting commit...")
			err = retryCommitRequest(ctx, contentBuilder, contentVersionID, fileID, encryptionInfo)
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Error during commit initiation: %v", err))
			}

		case graphmodels.AZURESTORAGEURIREQUESTPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE,
			graphmodels.AZURESTORAGEURIRENEWALPENDING_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "⏳ Azure storage URI request/renewal pending, waiting...")

		case graphmodels.AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE,
			graphmodels.AZURESTORAGEURIRENEWALFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "❌ Azure storage URI request/renewal failed")
			return fmt.Errorf("azure storage URI request/renewal failed: %s", state.String())

		case graphmodels.AZURESTORAGEURIREQUESTTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE,
			graphmodels.AZURESTORAGEURIRENEWALTIMEDOUT_MOBILEAPPCONTENTFILEUPLOADSTATE:
			tflog.Debug(ctx, "⏰ Azure storage URI request/renewal timed out")
			return fmt.Errorf("azure storage URI request/renewal timed out: %s", state.String())

		default:
			tflog.Debug(ctx, fmt.Sprintf("Unhandled state: %s, continuing to retry...", state.String()))
		}

		// Check if we've reached max retries
		if i == MaxCommitRetries-1 {
			return fmt.Errorf("file commit did not complete after %d attempts. Last state: %s",
				MaxCommitRetries, state.String())
		}

		// Apply backoff with jitter
		sleepTime := calculateBackoffWithJitter(backoff)
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v before next check...", sleepTime))
		time.Sleep(sleepTime)

		// Increase backoff for next iteration
		backoff = incrementBackoff(backoff)
	}

	return fmt.Errorf("reached end of retry loop without success or explicit failure")
}

// Helper function to retry commit request
func retryCommitRequest(
	ctx context.Context,
	contentBuilder *deviceappmanagement.MobileAppsItemGraphWin32LobAppContentVersionsRequestBuilder,
	contentVersionID string,
	fileID string,
	encryptionInfo *construct.EncryptionInfo,
) error {
	commitBody, err := construct.CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
	if err != nil {
		return fmt.Errorf("error constructing commit request: %v", err)
	}

	err = contentBuilder.
		ByMobileAppContentId(contentVersionID).
		Files().
		ByMobileAppContentFileId(fileID).
		Commit().
		Post(ctx, commitBody, nil)

	if err != nil {
		return fmt.Errorf("error posting commit request: %v", err)
	}

	tflog.Debug(ctx, "Successfully initiated commit request")
	return nil
}

// Calculate backoff with jitter
func calculateBackoffWithJitter(backoff time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(backoff) / 2))
	return backoff + jitter
}

// Increment backoff with capping
func incrementBackoff(backoff time.Duration) time.Duration {
	newBackoff := time.Duration(float64(backoff) * BackoffMultiplier)
	if newBackoff > MaxBackoff {
		return MaxBackoff
	}
	return newBackoff
}
