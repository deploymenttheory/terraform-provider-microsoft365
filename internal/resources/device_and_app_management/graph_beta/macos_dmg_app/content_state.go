package graphBetaMacOSDmgApp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// waitForCommitCompletion polls the commit state until completion or timeout for macOS DMG apps
func waitForCommitCompletion(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appId, contentVersionId string, retryTimeout time.Duration) error {
	tflog.Debug(ctx, fmt.Sprintf("Waiting for commit completion for DMG app %s, content version %s", appId, contentVersionId))

	return retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		// Get the mobile app to check its committed content version
		mobileApp, err := client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appId).
			GraphMacOSDmgApp().
			Get(ctx, nil)

		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to get mobile app: %v", err))
			return retry.RetryableError(fmt.Errorf("failed to get mobile app: %v", err))
		}

		if mobileApp.GetCommittedContentVersion() == nil {
			tflog.Debug(ctx, "Mobile app does not have a committed content version yet, retrying...")
			return retry.RetryableError(fmt.Errorf("mobile app does not have a committed content version yet"))
		}

		committedVersion := *mobileApp.GetCommittedContentVersion()
		tflog.Debug(ctx, fmt.Sprintf("Expected version: %s, Committed version: %s", contentVersionId, committedVersion))

		if committedVersion == contentVersionId {
			tflog.Debug(ctx, "File commit completed successfully")
			return nil
		}

		tflog.Debug(ctx, "File commit still in progress, retrying...")
		return retry.RetryableError(fmt.Errorf("file commit still in progress"))
	})
}
