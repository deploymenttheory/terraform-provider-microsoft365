package graphBetaIdentityAndAccessAuthenticationMethodsPolicy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
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

// restoredMethodIDs are the method configurations the acceptance tests change. Cleanup PATCHes each one back
// to its recorded value; every other part of the policy is only verified.
var restoredMethodIDs = []string{"Email"}

// restoredMethodFields are the writable fields of a restored method configuration that cleanup sends back.
var restoredMethodFields = []string{"@odata.type", "state", "allowExternalIdToUseEmailOtp", "includeTargets", "excludeTargets"}

// originalPolicy is the tenant policy recorded before the first acceptance test changes it.
// It is captured once per test binary so a failed restore cannot become the next test's baseline.
var originalPolicy struct {
	once sync.Once
	body map[string]any
	err  error
}

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// authenticationMethodsPolicyPreCheck records the tenant authentication methods policy and registers a cleanup
// that PATCHes the changed method configurations back and verifies with GET that the whole policy matches (rule R9).
func authenticationMethodsPolicyPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalPolicy.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalPolicy.body, originalPolicy.err = getRawAuthenticationMethodsPolicy(ctx, client)
	})
	require.NoError(t, originalPolicy.err)
	original := originalPolicy.body

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := restoreAuthenticationMethodsPolicy(ctx, client, original); err != nil {
			t.Errorf("restore original authentication methods policy: %v", err)
		}
	})
}

// getRawAuthenticationMethodsPolicy reads the policy, including its method configurations, as JSON.
func getRawAuthenticationMethodsPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]any, error) {
	requestInfo, err := client.Policies().AuthenticationMethodsPolicy().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected authentication methods policy response type %T", result)
	}
	var policy map[string]any
	if err := json.Unmarshal(content, &policy); err != nil {
		return nil, err
	}
	return policy, nil
}

// methodConfiguration returns the method configuration with the given ID from a raw policy.
func methodConfiguration(policy map[string]any, id string) (map[string]any, bool) {
	configurations, _ := policy["authenticationMethodConfigurations"].([]any)
	for _, item := range configurations {
		configuration, _ := item.(map[string]any)
		if configuration["id"] == id {
			return configuration, true
		}
	}
	return nil, false
}

// restoreAuthenticationMethodsPolicy PATCHes each changed method configuration back to its recorded value and
// waits until three consecutive reads of the whole policy match the recorded policy.
func restoreAuthenticationMethodsPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, original map[string]any) error {
	for _, id := range restoredMethodIDs {
		configuration, ok := methodConfiguration(original, id)
		if !ok {
			return fmt.Errorf("recorded policy has no %s method configuration", id)
		}
		body := map[string]any{}
		for _, field := range restoredMethodFields {
			if value, present := configuration[field]; present {
				body[field] = value
			}
		}
		content, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestInfo, err := client.Policies().AuthenticationMethodsPolicy().AuthenticationMethodConfigurations().ByAuthenticationMethodConfigurationId(id).ToPatchRequestInformation(ctx, nil, nil)
		if err != nil {
			return err
		}
		requestInfo.Content = content
		requestInfo.Headers.TryAdd("Content-Type", "application/json")
		if err := client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping); err != nil {
			return fmt.Errorf("PATCH original %s method configuration: %w", id, err)
		}
	}

	expected := normalisePolicy(original)
	consecutive := 0
	var current map[string]any
	var err error
	for {
		current, err = getRawAuthenticationMethodsPolicy(ctx, client)
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
			return fmt.Errorf("authentication methods policy did not return to its original settings: last read %v, error %v", current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// normalisePolicy drops values that change on every write (lastModifiedDateTime and OData context links) and
// sorts collections so target order does not affect the comparison.
func normalisePolicy(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalised := make(map[string]any, len(typed))
		for key, nested := range typed {
			if key == "lastModifiedDateTime" || strings.HasSuffix(key, "@odata.context") {
				continue
			}
			normalised[key] = normalisePolicy(nested)
		}
		return normalised
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			encoded, _ := json.Marshal(normalisePolicy(item))
			items = append(items, string(encoded))
		}
		slices.Sort(items)
		return items
	default:
		return value
	}
}

// checkPolicyStillExists asserts that destroy left the tenant authentication methods policy in place.
func checkPolicyStillExists(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	policy, err := client.Policies().AuthenticationMethodsPolicy().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("expected the authentication methods policy to remain after destroy: %w", err)
	}
	if policy == nil || policy.GetId() == nil || *policy.GetId() != "authenticationMethodsPolicy" {
		return fmt.Errorf("expected the authentication methods policy to remain after destroy, got %v", policy)
	}
	return nil
}

// checkUnmanagedMethodsUnchanged asserts that every method configuration other than Email still matches the recorded policy.
func checkUnmanagedMethodsUnchanged(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	current, err := getRawAuthenticationMethodsPolicy(ctx, client)
	if err != nil {
		return err
	}
	configurations, _ := originalPolicy.body["authenticationMethodConfigurations"].([]any)
	for _, item := range configurations {
		original, _ := item.(map[string]any)
		id, _ := original["id"].(string)
		if slices.Contains(restoredMethodIDs, id) {
			continue
		}
		actual, ok := methodConfiguration(current, id)
		if !ok {
			return fmt.Errorf("method configuration %s is missing", id)
		}
		if !reflect.DeepEqual(normalisePolicy(original), normalisePolicy(actual)) {
			return fmt.Errorf("method configuration %s changed although it is not managed by the test", id)
		}
	}
	return nil
}
