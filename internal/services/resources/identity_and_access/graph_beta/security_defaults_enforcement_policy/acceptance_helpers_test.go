package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy_test

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

// originalPolicy is the tenant policy recorded before the first acceptance test runs.
// It is captured once per test binary so a failed restore cannot become the next test's baseline.
var originalPolicy struct {
	once sync.Once
	body map[string]any
	err  error
}

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// securityDefaultsPreCheck records the tenant security defaults policy and registers a cleanup that PATCHes the
// recorded isEnabled back when it changed and verifies with GET that the policy matches the snapshot (R9).
//
// Enabling security defaults enforces MFA registration for every user in the tenant and conflicts with the
// Conditional Access acceptance suites, so the acceptance tests only ever apply is_enabled = false. They are
// skipped when security defaults are already enabled, because applying false would then disable them.
func securityDefaultsPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalPolicy.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalPolicy.body, originalPolicy.err = getRawSecurityDefaultsPolicy(ctx, client)
	})
	require.NoError(t, originalPolicy.err)
	original := originalPolicy.body

	if original["isEnabled"] != false {
		t.Skipf("security defaults are not disabled in this tenant (isEnabled=%v); the acceptance tests only apply is_enabled = false and never change an enabled tenant", original["isEnabled"])
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := restoreSecurityDefaultsPolicy(ctx, client, original); err != nil {
			t.Errorf("restore original security defaults policy: %v", err)
		}
	})
}

// getRawSecurityDefaultsPolicy reads the policy as JSON so it can be compared and sent back unchanged.
func getRawSecurityDefaultsPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]any, error) {
	requestInfo, err := client.Policies().IdentitySecurityDefaultsEnforcementPolicy().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected security defaults policy response type %T", result)
	}
	var policy map[string]any
	if err := json.Unmarshal(content, &policy); err != nil {
		return nil, err
	}
	delete(policy, "@odata.context")
	return policy, nil
}

// comparableProperties are the policy properties a restore is verified against. securityDefaultsUpsell is left out:
// it is a read-only nudge that Microsoft Entra ID maintains on its own (its dueDateTime moves with service time),
// so it is neither written by the resource nor restorable.
var comparableProperties = []string{"id", "displayName", "description", "isEnabled", "controlTypes"}

// policyMatches reports whether the comparable properties of current equal those of original.
func policyMatches(original, current map[string]any) bool {
	for _, property := range comparableProperties {
		if !reflect.DeepEqual(original[property], current[property]) {
			return false
		}
	}
	return true
}

// restoreSecurityDefaultsPolicy PATCHes the recorded isEnabled (the only writable property) back when it differs
// from the current value and waits until three consecutive reads match the snapshot.
func restoreSecurityDefaultsPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, original map[string]any) error {
	isEnabled, ok := original["isEnabled"].(bool)
	if !ok {
		return fmt.Errorf("snapshot has no Boolean isEnabled: %v", original)
	}
	current, err := getRawSecurityDefaultsPolicy(ctx, client)
	if err != nil {
		return fmt.Errorf("read current policy: %w", err)
	}
	if current["isEnabled"] != isEnabled {
		content, err := json.Marshal(map[string]any{"isEnabled": isEnabled})
		if err != nil {
			return err
		}
		requestInfo, err := client.Policies().IdentitySecurityDefaultsEnforcementPolicy().ToGetRequestInformation(ctx, nil)
		if err != nil {
			return err
		}
		requestInfo.Method = abstractions.PATCH
		requestInfo.Content = content
		requestInfo.Headers.TryAdd("Content-Type", "application/json")
		if err := client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping); err != nil {
			return fmt.Errorf("PATCH original isEnabled: %w", err)
		}
	}

	consecutive := 0
	for {
		current, err = getRawSecurityDefaultsPolicy(ctx, client)
		if err == nil && policyMatches(original, current) {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("security defaults policy did not return to its original values: last read %v, error %v", current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// checkTenantSecurityDefaultsDisabled asserts directly against Graph that security defaults are still disabled.
func checkTenantSecurityDefaultsDisabled(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	policy, err := client.Policies().IdentitySecurityDefaultsEnforcementPolicy().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("expected the security defaults policy to remain readable: %w", err)
	}
	if policy == nil || policy.GetIsEnabled() == nil || *policy.GetIsEnabled() {
		return fmt.Errorf("expected security defaults to remain disabled in the tenant, got %v", policy)
	}
	return nil
}
