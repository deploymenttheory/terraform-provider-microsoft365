package graphBetaIdentityAndAccessLifecycleWorkflowSettings_test

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

// restorableProperties are the settings the resource can change; only these are PATCHed back on restore.
var restorableProperties = []string{"workflowScheduleIntervalInHours", "emailSettings"}

// lifecycleWorkflowSettingsPreCheck records the tenant lifecycle workflow settings and registers a cleanup that
// PATCHes the recorded values back and verifies with GET that the whole settings object matches the snapshot (R9).
func lifecycleWorkflowSettingsPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalSettings.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalSettings.body, originalSettings.err = getRawLifecycleWorkflowSettings(ctx, client)
	})
	require.NoError(t, originalSettings.err)
	original := originalSettings.body

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := restoreLifecycleWorkflowSettings(ctx, client, original); err != nil {
			t.Errorf("restore original lifecycle workflow settings: %v", err)
		}
	})
}

// getRawLifecycleWorkflowSettings reads the settings as JSON so they can be sent back unchanged.
func getRawLifecycleWorkflowSettings(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]any, error) {
	requestInfo, err := client.IdentityGovernance().LifecycleWorkflows().Settings().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected lifecycle workflow settings response type %T", result)
	}
	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		return nil, err
	}
	delete(settings, "@odata.context")
	return settings, nil
}

// restoreLifecycleWorkflowSettings PATCHes the recorded values and waits until three consecutive reads match the snapshot.
func restoreLifecycleWorkflowSettings(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, original map[string]any) error {
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
	requestInfo, err := client.IdentityGovernance().LifecycleWorkflows().Settings().ToGetRequestInformation(ctx, nil)
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
		current, err := getRawLifecycleWorkflowSettings(ctx, client)
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
			return fmt.Errorf("lifecycle workflow settings did not return to their original values: last read %v, error %v", current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// checkSettingsStillExist asserts that destroy left the tenant lifecycle workflow settings in place.
func checkSettingsStillExist(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	settings, err := client.IdentityGovernance().LifecycleWorkflows().Settings().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("expected the lifecycle workflow settings to remain after destroy: %w", err)
	}
	if settings == nil || settings.GetWorkflowScheduleIntervalInHours() == nil {
		return fmt.Errorf("expected the lifecycle workflow settings to remain after destroy, got %v", settings)
	}
	return nil
}
