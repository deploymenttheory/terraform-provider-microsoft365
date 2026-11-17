package graphBetaM365AdminSharepointSettings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// originalSettings are the tenant settings recorded before the first acceptance test changes them.
// They are captured once per test binary so a failed restore cannot become the next test's baseline.
var originalSettings struct {
	once sync.Once
	body map[string]any
	err  error
}

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// restorableProperties are the only settings the acceptance tests change; only these are PATCHed back on restore.
// The tests never change sharing, Loop, legacy authentication or unmanaged device settings.
var restorableProperties = []string{"deletedUserPersonalSiteRetentionPeriodInDays", "excludedFileExtensionsForSyncApp", "isMacSyncAppEnabled", "isSyncButtonHiddenOnPersonalSite", "idleSessionSignOut"}

// sharepointSettingsPreCheck records the tenant SharePoint settings and registers a cleanup that PATCHes the
// recorded values back and verifies with GET that the whole settings object matches the snapshot (R9).
func sharepointSettingsPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalSettings.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalSettings.body, originalSettings.err = getRawSharepointSettings(ctx, client)
	})
	require.NoError(t, originalSettings.err)
	original := originalSettings.body

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := restoreSharepointSettings(ctx, client, original); err != nil {
			t.Errorf("restore original SharePoint settings: %v", err)
		}
	})
}

// getRawSharepointSettings reads the settings as JSON so they can be sent back unchanged.
func getRawSharepointSettings(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]any, error) {
	requestInfo, err := client.Admin().Sharepoint().Settings().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected SharePoint settings response type %T", result)
	}
	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		return nil, err
	}
	delete(settings, "@odata.context")
	return settings, nil
}

// restoreSharepointSettings PATCHes the recorded values and waits until three consecutive reads match the snapshot.
func restoreSharepointSettings(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, original map[string]any) error {
	patch := map[string]any{}
	for _, property := range restorableProperties {
		if value, ok := original[property]; ok && value != nil {
			patch[property] = value
		}
	}
	content, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	requestInfo, err := client.Admin().Sharepoint().Settings().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return err
	}
	requestInfo.Method = abstractions.PATCH
	requestInfo.Content = content
	requestInfo.Headers.TryAdd("Content-Type", "application/json")
	if err := client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping); err != nil {
		return fmt.Errorf("PATCH original settings: %w", err)
	}

	consecutive := 0
	for {
		current, err := getRawSharepointSettings(ctx, client)
		if err == nil && reflect.DeepEqual(original, current) {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("SharePoint settings did not return to their original values: last read %v, error %v", current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// checkSettingsStillExist asserts that destroy left the tenant SharePoint settings in place.
func checkSettingsStillExist(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	settings, err := client.Admin().Sharepoint().Settings().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("expected the SharePoint settings to remain after destroy: %w", err)
	}
	if settings == nil || settings.GetDeletedUserPersonalSiteRetentionPeriodInDays() == nil {
		return fmt.Errorf("expected the SharePoint settings to remain after destroy, got %v", settings)
	}
	return nil
}
