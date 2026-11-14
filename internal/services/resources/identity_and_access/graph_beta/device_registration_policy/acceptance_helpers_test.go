package graphBetaIdentityAndAccessDeviceRegistrationPolicy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
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

// originalPolicy is the tenant policy recorded before the first acceptance test changes it.
// It is captured once per test binary so a failed restore cannot become the next test's baseline.
var originalPolicy struct {
	once sync.Once
	body map[string]any
	err  error
}

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// deviceRegistrationPolicyPreCheck records the tenant device registration policy and registers a
// cleanup that PUTs the recorded policy back and verifies the restore with GET (rule R9).
func deviceRegistrationPolicyPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalPolicy.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalPolicy.body, originalPolicy.err = getRawDeviceRegistrationPolicy(ctx, client)
	})
	require.NoError(t, originalPolicy.err)
	original := originalPolicy.body

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := restoreDeviceRegistrationPolicy(ctx, client, original); err != nil {
			t.Errorf("restore original device registration policy: %v", err)
		}
	})
}

// getRawDeviceRegistrationPolicy reads the policy as JSON so it can be sent back unchanged.
func getRawDeviceRegistrationPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]any, error) {
	requestInfo, err := client.Policies().DeviceRegistrationPolicy().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected device registration policy response type %T", result)
	}
	var policy map[string]any
	if err := json.Unmarshal(content, &policy); err != nil {
		return nil, err
	}
	delete(policy, "@odata.context")
	return policy, nil
}

// restoreDeviceRegistrationPolicy PUTs the recorded policy and waits until three consecutive reads match it.
func restoreDeviceRegistrationPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, original map[string]any) error {
	content, err := json.Marshal(original)
	if err != nil {
		return err
	}
	requestInfo, err := client.Policies().DeviceRegistrationPolicy().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return err
	}
	requestInfo.Method = abstractions.PUT
	requestInfo.Content = content
	requestInfo.Headers.TryAdd("Content-Type", "application/json")
	if err := client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping); err != nil {
		return fmt.Errorf("PUT original policy: %w", err)
	}

	expected := normalisePolicy(original)
	consecutive := 0
	for {
		current, err := getRawDeviceRegistrationPolicy(ctx, client)
		if err == nil && reflect.DeepEqual(expected, normalisePolicy(current)) {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("device registration policy did not return to its original settings: last read %v, error %v", current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// normalisePolicy sorts string lists so group and user order does not affect the comparison.
func normalisePolicy(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalised := make(map[string]any, len(typed))
		for key, nested := range typed {
			normalised[key] = normalisePolicy(nested)
		}
		return normalised
	case []any:
		strings := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return typed
			}
			strings = append(strings, text)
		}
		slices.Sort(strings)
		return strings
	default:
		return value
	}
}

// checkPolicyStillExists asserts that destroy left the tenant device registration policy in place.
func checkPolicyStillExists(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	policy, err := client.Policies().DeviceRegistrationPolicy().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("expected the device registration policy to remain after destroy: %w", err)
	}
	if policy == nil || policy.GetId() == nil || *policy.GetId() != "deviceRegistrationPolicy" {
		return fmt.Errorf("expected the device registration policy to remain after destroy, got %v", policy)
	}
	return nil
}
