package graphBetaRoleManagementPolicy_test

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/policies"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// accRoleDefinitionId is the low-privilege built-in Usage Summary Reports Reader role whose PIM policy the acceptance
// tests change and restore.
const accRoleDefinitionId = "75934031-6c7e-415a-99d7-48dbd49e875e"

// managedRuleTypes are the rule types this resource manages and the acceptance cleanup restores.
var managedRuleTypes = []string{
	"#microsoft.graph.unifiedRoleManagementPolicyExpirationRule",
	"#microsoft.graph.unifiedRoleManagementPolicyEnablementRule",
	"#microsoft.graph.unifiedRoleManagementPolicyNotificationRule",
	"#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule",
}

// originalPolicy is the policy recorded before the first acceptance test changes it. It is captured once per test
// binary so a failed restore cannot become the next test's baseline.
var originalPolicy struct {
	once     sync.Once
	policyId string
	rules    []map[string]any
	err      error
}

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// roleManagementPolicyPreCheck records the rules of the acceptance role's policy and registers a cleanup that PATCHes
// the recorded rules back and verifies the restore with GET (rule R9). It returns the policy ID.
func roleManagementPolicyPreCheck(t *testing.T) string {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	originalPolicy.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		originalPolicy.policyId, originalPolicy.err = findAccPolicyId(ctx, client)
		if originalPolicy.err == nil {
			originalPolicy.rules, originalPolicy.err = getRawRules(ctx, client, originalPolicy.policyId)
		}
	})
	require.NoError(t, originalPolicy.err)
	policyId := originalPolicy.policyId
	original := originalPolicy.rules

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if err := restoreRules(ctx, client, policyId, original); err != nil {
			t.Errorf("restore original role management policy rules: %v", err)
		}
	})

	return policyId
}

// findAccPolicyId looks up the policy assigned to the acceptance role at the directory scope.
func findAccPolicyId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (string, error) {
	filter := fmt.Sprintf("scopeId eq '/' and scopeType eq 'DirectoryRole' and roleDefinitionId eq '%s'", accRoleDefinitionId)
	config := &policies.RoleManagementPolicyAssignmentsRequestBuilderGetRequestConfiguration{
		QueryParameters: &policies.RoleManagementPolicyAssignmentsRequestBuilderGetQueryParameters{Filter: &filter},
	}
	response, err := client.Policies().RoleManagementPolicyAssignments().Get(ctx, config)
	if err != nil {
		return "", err
	}
	if response == nil || len(response.GetValue()) != 1 || response.GetValue()[0].GetPolicyId() == nil {
		return "", fmt.Errorf("expected one policy assignment for role definition %s", accRoleDefinitionId)
	}
	return *response.GetValue()[0].GetPolicyId(), nil
}

// getRawRules reads the policy rules as JSON so they can be sent back unchanged.
func getRawRules(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId string) ([]map[string]any, error) {
	requestInfo, err := client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(policyId).Rules().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected role management policy rules response type %T", result)
	}
	var response struct {
		Value []map[string]any `json:"value"`
	}
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}
	if len(response.Value) == 0 {
		return nil, fmt.Errorf("role management policy %s returned no rules", policyId)
	}
	return response.Value, nil
}

// restoreRules PATCHes the recorded values of every managed rule type and waits until three consecutive reads of all
// rules match the recording.
func restoreRules(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId string, original []map[string]any) error {
	restore := make([]map[string]any, 0, len(original))
	for _, rule := range original {
		if !slices.Contains(managedRuleTypes, fmt.Sprint(rule["@odata.type"])) {
			continue
		}
		body := map[string]any{}
		for key, value := range rule {
			if key != "target" {
				body[key] = value
			}
		}
		restore = append(restore, body)
	}

	content, err := json.Marshal(map[string]any{"rules": restore})
	if err != nil {
		return err
	}
	requestInfo, err := client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(policyId).ToGetRequestInformation(ctx, nil)
	if err != nil {
		return err
	}
	requestInfo.Method = abstractions.PATCH
	requestInfo.Content = content
	requestInfo.Headers.TryAdd("Content-Type", "application/json")
	if err := client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping); err != nil {
		return fmt.Errorf("PATCH original rules: %w", err)
	}

	expected := normaliseRules(original)
	consecutive := 0
	for {
		current, err := getRawRules(ctx, client, policyId)
		if err == nil && reflect.DeepEqual(expected, normaliseRules(current)) {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("role management policy %s did not return to its original rules: last read %v, error %v", policyId, current, err)
		case <-time.After(2 * time.Second):
		}
	}
}

// normaliseRules indexes rules by ID and sorts string lists so ordering does not affect the comparison.
func normaliseRules(rules []map[string]any) map[string]any {
	normalised := make(map[string]any, len(rules))
	for _, rule := range rules {
		normalised[fmt.Sprint(rule["id"])] = normaliseValue(rule)
	}
	return normalised
}

func normaliseValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalised := make(map[string]any, len(typed))
		for key, nested := range typed {
			normalised[key] = normaliseValue(nested)
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

// checkPolicyStillExists asserts that destroy left the role management policy in place.
func checkPolicyStillExists(_ *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	policyId, err := findAccPolicyId(ctx, client)
	if err != nil {
		return fmt.Errorf("expected the role management policy to remain after destroy: %w", err)
	}
	if _, err := client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(policyId).Get(ctx, nil); err != nil {
		return fmt.Errorf("expected the role management policy %s to remain after destroy: %w", policyId, err)
	}
	return nil
}

// checkAccImportedAllRules asserts that an import of the acceptance policy tracks every supported rule.
func checkAccImportedAllRules(policyId string) func(states []*terraform.InstanceState) error {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one imported instance, got %d", len(states))
		}
		attributes := states[0].Attributes
		expected := map[string]string{
			"id":                             policyId,
			"policy_id":                      policyId,
			"role_definition_id":             accRoleDefinitionId,
			"scope_id":                       "/",
			"scope_type":                     "DirectoryRole",
			"expiration_rules.%":             "3",
			"enablement_rules.%":             "3",
			"notification_rules.%":           "9",
			"authentication_context_rules.%": "1",
		}
		for key, want := range expected {
			if got := attributes[key]; got != want {
				return fmt.Errorf("expected imported %s to be %q, got %q", key, want, got)
			}
		}
		return nil
	}
}
