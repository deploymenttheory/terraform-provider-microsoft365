package graphBetaServicePrincipalSynchronizationSchemaJson_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// customApplicationTemplateId is the "Custom" (non-gallery) application template, which supports SCIM provisioning.
	customApplicationTemplateId = "8adf8e6e-67b2-4cf2-a259-e3dc5476c621"

	// scimSynchronizationTemplateId is the synchronization template the acceptance job is created from. Creating the job
	// needs no SCIM endpoint or credentials; the job stays paused and is never started.
	scimSynchronizationTemplateId = "scim"
)

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

// synchronizationFixture is an acc-test- application, its service principal and a synchronization job created for one test.
type synchronizationFixture struct {
	ApplicationObjectId string
	ServicePrincipalId  string
	JobId               string
	DefaultRules        []any
	DefaultDirectories  []any
}

// synchronizationJobPreCheck instantiates an acc-test- custom application, creates a SCIM synchronization job on its
// service principal and records the job's default schema. A t.Cleanup deletes the application and its service
// principal by their exact IDs and purges both from deleted items.
func synchronizationJobPreCheck(t *testing.T) *synchronizationFixture {
	t.Helper()
	mocks.TestAccPreCheck(t)

	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	displayName := "acc-test-sync-schema-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	fixture, err := instantiateApplication(ctx, client, displayName)
	if fixture != nil {
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cleanupCancel()
			if err := deleteFixture(cleanupCtx, client, fixture); err != nil {
				t.Errorf("delete acceptance application %s (service principal %s): %v", fixture.ApplicationObjectId, fixture.ServicePrincipalId, err)
				return
			}
			t.Logf("Deleted and purged acceptance application %s and service principal %s", fixture.ApplicationObjectId, fixture.ServicePrincipalId)
		})
	}
	require.NoError(t, err)
	t.Logf("Created acceptance application %s (%s) with service principal %s", displayName, fixture.ApplicationObjectId, fixture.ServicePrincipalId)

	fixture.JobId, err = createSynchronizationJob(ctx, client, fixture.ServicePrincipalId)
	require.NoError(t, err)
	t.Logf("Created synchronization job %s", fixture.JobId)

	schema, err := waitForRawSchema(ctx, client, fixture.ServicePrincipalId, fixture.JobId)
	require.NoError(t, err)
	fixture.DefaultRules, _ = schema["synchronizationRules"].([]any)
	fixture.DefaultDirectories, _ = schema["directories"].([]any)
	require.NotEmpty(t, fixture.DefaultRules, "the default schema has no synchronization rules")
	require.NotEmpty(t, fixture.DefaultDirectories, "the default schema has no directories")

	return fixture
}

// sendJSON sends a Graph request through the test client's adapter and decodes the JSON response when one is wanted.
func sendJSON(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, method abstractions.HttpMethod, path string, body any, response any) error {
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, "{+baseurl}"+path, map[string]string{"baseurl": client.GetAdapter().GetBaseUrl()})
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if body != nil {
		content, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestInfo.Content = content
		requestInfo.Headers.TryAdd("Content-Type", "application/json")
	}
	if response == nil {
		return client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping)
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return err
	}
	content, ok := result.([]byte)
	if !ok {
		return fmt.Errorf("unexpected response type %T", result)
	}
	return json.Unmarshal(content, response)
}

// statusCode returns the HTTP status of a Graph error, or 0.
func statusCode(err error) int {
	var apiErr *abstractions.ApiError
	if errors.As(err, &apiErr) {
		return apiErr.ResponseStatusCode
	}
	var odataErr *odataerrors.ODataError
	if errors.As(err, &odataErr) {
		return odataErr.ResponseStatusCode
	}
	return 0
}

// instantiateApplication creates an application and service principal from the custom application template. The
// instantiate response uses objectId for the object IDs.
func instantiateApplication(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string) (*synchronizationFixture, error) {
	var response struct {
		Application struct {
			ID       string `json:"id"`
			ObjectID string `json:"objectId"`
		} `json:"application"`
		ServicePrincipal struct {
			ID       string `json:"id"`
			ObjectID string `json:"objectId"`
		} `json:"servicePrincipal"`
	}
	path := "/applicationTemplates/" + customApplicationTemplateId + "/instantiate"
	if err := sendJSON(ctx, client, abstractions.POST, path, map[string]any{"displayName": displayName}, &response); err != nil {
		return nil, fmt.Errorf("instantiate application template: %w", err)
	}

	fixture := &synchronizationFixture{ApplicationObjectId: response.Application.ID, ServicePrincipalId: response.ServicePrincipal.ID}
	if fixture.ApplicationObjectId == "" {
		fixture.ApplicationObjectId = response.Application.ObjectID
	}
	if fixture.ServicePrincipalId == "" {
		fixture.ServicePrincipalId = response.ServicePrincipal.ObjectID
	}
	if fixture.ApplicationObjectId == "" || fixture.ServicePrincipalId == "" {
		return nil, fmt.Errorf("instantiate response did not include the application and service principal IDs")
	}
	return fixture, nil
}

// createSynchronizationJob creates a SCIM synchronization job, retrying while the new service principal replicates.
func createSynchronizationJob(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, servicePrincipalId string) (string, error) {
	var lastErr error
	for {
		var response struct {
			ID string `json:"id"`
		}
		path := "/servicePrincipals/" + servicePrincipalId + "/synchronization/jobs"
		lastErr = sendJSON(ctx, client, abstractions.POST, path, map[string]any{"templateId": scimSynchronizationTemplateId}, &response)
		if lastErr == nil && response.ID != "" {
			return response.ID, nil
		}
		if lastErr != nil {
			status := statusCode(lastErr)
			if status != 401 && status != 403 && status != 404 {
				return "", fmt.Errorf("create synchronization job: %w", lastErr)
			}
		}
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("create synchronization job: %w (last error: %v)", ctx.Err(), lastErr)
		case <-time.After(5 * time.Second):
		}
	}
}

// getRawSchema reads the synchronization job schema as JSON.
func getRawSchema(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, servicePrincipalId, jobId string) (map[string]any, error) {
	requestInfo, err := client.ServicePrincipals().ByServicePrincipalId(servicePrincipalId).Synchronization().Jobs().BySynchronizationJobId(jobId).Schema().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err
	}
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected synchronization schema response type %T", result)
	}
	var schema map[string]any
	if err := json.Unmarshal(content, &schema); err != nil {
		return nil, err
	}
	return schema, nil
}

// waitForRawSchema reads the schema of a new synchronization job, retrying while the job replicates.
func waitForRawSchema(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, servicePrincipalId, jobId string) (map[string]any, error) {
	for {
		schema, err := getRawSchema(ctx, client, servicePrincipalId, jobId)
		if err == nil {
			return schema, nil
		}
		status := statusCode(err)
		if status != 401 && status != 403 && status != 404 {
			return nil, fmt.Errorf("read synchronization schema (HTTP %d): %w", status, err)
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("read synchronization schema (HTTP %d, %v): %w", status, err, ctx.Err())
		case <-time.After(5 * time.Second):
		}
	}
}

// deleteFixture deletes the application by its exact ID, which also deletes its service principal and synchronization
// job, then purges the application and service principal from deleted items. A 404 right after creation can come from
// a replica that has not seen the object yet, so every step repeats until three consecutive reads agree.
func deleteFixture(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, fixture *synchronizationFixture) error {
	if err := deleteApplication(ctx, client, fixture.ApplicationObjectId); err != nil {
		return err
	}
	for _, id := range []string{fixture.ApplicationObjectId, fixture.ServicePrincipalId} {
		if err := purgeDeletedItem(ctx, client, id); err != nil {
			return err
		}
	}
	return nil
}

// deleteApplication deletes the application until it is no longer returned as an active directory object.
func deleteApplication(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, id string) error {
	consecutive := 0
	for {
		err := client.Applications().ByApplicationId(id).Delete(ctx, nil)
		if err != nil && statusCode(err) != 404 {
			return fmt.Errorf("delete application %s: %w", id, err)
		}
		_, getErr := client.DirectoryObjects().ByDirectoryObjectId(id).Get(ctx, nil)
		if statusCode(getErr) == 404 {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("application %s was not deleted: %w", id, ctx.Err())
		case <-time.After(5 * time.Second):
		}
	}
}

// purgeDeletedItem permanently deletes one object from deleted items and waits until neither deleted items nor the
// active directory objects return it.
func purgeDeletedItem(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, id string) error {
	consecutive := 0
	for {
		err := client.Directory().DeletedItems().ByDirectoryObjectId(id).Delete(ctx, nil)
		if err != nil && statusCode(err) != 404 {
			return fmt.Errorf("purge deleted item %s: %w", id, err)
		}
		_, deletedErr := client.Directory().DeletedItems().ByDirectoryObjectId(id).Get(ctx, nil)
		_, activeErr := client.DirectoryObjects().ByDirectoryObjectId(id).Get(ctx, nil)
		if statusCode(deletedErr) == 404 && statusCode(activeErr) == 404 {
			consecutive++
			if consecutive >= 3 {
				return nil
			}
		} else {
			consecutive = 0
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("deleted item %s was not purged: %w", id, ctx.Err())
		case <-time.After(5 * time.Second):
		}
	}
}

// copyJSON returns an independent copy of a decoded JSON value.
func copyJSON(value []any) []any {
	content, _ := json.Marshal(value)
	var copied []any
	_ = json.Unmarshal(content, &copied)
	return copied
}

// encodeJSON returns the JSON string of a decoded value.
func encodeJSON(t *testing.T, value any) string {
	t.Helper()
	content, err := json.Marshal(value)
	require.NoError(t, err)
	return string(content)
}

// findObjectMapping returns the object mapping of the first rule whose source object is sourceObjectName.
func findObjectMapping(rules []any, sourceObjectName string) map[string]any {
	for _, ruleItem := range rules {
		rule, _ := ruleItem.(map[string]any)
		objectMappings, _ := rule["objectMappings"].([]any)
		for _, item := range objectMappings {
			mapping, _ := item.(map[string]any)
			if mapping["sourceObjectName"] == sourceObjectName {
				return mapping
			}
		}
	}
	return nil
}

// findDirectoryObject returns the first object of the directory named directoryName whose name ends with suffix.
func findDirectoryObject(directories []any, directoryName, suffix string) map[string]any {
	for _, directoryItem := range directories {
		directory, _ := directoryItem.(map[string]any)
		if directory["name"] != directoryName {
			continue
		}
		objects, _ := directory["objects"].([]any)
		for _, item := range objects {
			object, _ := item.(map[string]any)
			if name, _ := object["name"].(string); strings.HasSuffix(name, suffix) {
				return object
			}
		}
	}
	return nil
}

// customAttributeName is the SCIM extension attribute the maximal acceptance schema adds to the target directory.
const customAttributeName = "urn:ietf:params:scim:schemas:extension:acctest:2.0:User:costCenter"

// minimalSchema returns the default schema with the group object mapping disabled.
func minimalSchema(t *testing.T, fixture *synchronizationFixture) (string, string) {
	t.Helper()
	rules := copyJSON(fixture.DefaultRules)
	groups := findObjectMapping(rules, "Group")
	require.NotNil(t, groups, "the default schema has no group object mapping")
	groups["enabled"] = false
	return encodeJSON(t, rules), encodeJSON(t, fixture.DefaultDirectories)
}

// maximalSchema returns the minimal schema plus a custom SCIM attribute on the target user object and a mapping from
// the Entra ID department attribute to it.
func maximalSchema(t *testing.T, fixture *synchronizationFixture) (string, string) {
	t.Helper()
	rules := copyJSON(fixture.DefaultRules)
	directories := copyJSON(fixture.DefaultDirectories)

	groups := findObjectMapping(rules, "Group")
	require.NotNil(t, groups, "the default schema has no group object mapping")
	groups["enabled"] = false

	users := findObjectMapping(rules, "User")
	require.NotNil(t, users, "the default schema has no user object mapping")
	targetObjectName, _ := users["targetObjectName"].(string)
	target := findDirectoryObject(directories, "customappsso", targetObjectName)
	require.NotNil(t, target, "the default schema has no target user object %s", targetObjectName)

	attributes, _ := target["attributes"].([]any)
	target["attributes"] = append(attributes, map[string]any{
		"anchor": false, "caseExact": false, "defaultValue": nil, "flowNullValues": false, "multivalued": false, "mutability": "ReadWrite",
		"name": customAttributeName, "required": false, "type": "String", "apiExpressions": []any{}, "metadata": []any{}, "referencedObjects": []any{},
	})
	attributeMappings, _ := users["attributeMappings"].([]any)
	users["attributeMappings"] = append(attributeMappings, map[string]any{
		"defaultValue": nil, "exportMissingReferences": false, "flowBehavior": "FlowWhenChanged", "flowType": "Always", "matchingPriority": 0,
		"targetAttributeName": customAttributeName,
		"source":              map[string]any{"expression": "[department]", "name": "department", "type": "Attribute", "parameters": []any{}},
	})
	return encodeJSON(t, rules), encodeJSON(t, directories)
}

// checkGraphSchema reads the job schema from Graph and runs check on it.
func checkGraphSchema(fixture *synchronizationFixture, check func(schema map[string]any) error) func(*terraform.State) error {
	return func(_ *terraform.State) error {
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		schema, err := waitForRawSchema(ctx, client, fixture.ServicePrincipalId, fixture.JobId)
		if err != nil {
			return fmt.Errorf("read synchronization schema: %w", err)
		}
		return check(schema)
	}
}

// checkGroupMappingEnabled asserts the enabled flag of the group object mapping in Graph.
func checkGroupMappingEnabled(fixture *synchronizationFixture, want bool) func(*terraform.State) error {
	return checkGraphSchema(fixture, func(schema map[string]any) error {
		rules, _ := schema["synchronizationRules"].([]any)
		groups := findObjectMapping(rules, "Group")
		if groups == nil {
			return fmt.Errorf("the synchronization schema has no group object mapping")
		}
		if groups["enabled"] != want {
			return fmt.Errorf("expected the group object mapping enabled to be %t, got %v", want, groups["enabled"])
		}
		return nil
	})
}

// checkCustomAttribute asserts whether the custom attribute and its mapping are in the Graph schema.
func checkCustomAttribute(fixture *synchronizationFixture, want bool) func(*terraform.State) error {
	return checkGraphSchema(fixture, func(schema map[string]any) error {
		content, err := json.Marshal(schema)
		if err != nil {
			return err
		}
		if got := strings.Count(string(content), customAttributeName); (got == 2) != want {
			return fmt.Errorf("expected the custom attribute present to be %t, found it %d times", want, got)
		}
		return nil
	})
}

// normaliseSchema removes the version properties Graph changes on every write, for comparing schemas.
func normaliseSchema(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalised := make(map[string]any, len(typed))
		for key, nested := range typed {
			if key == "version" || strings.HasPrefix(key, "@odata.") || strings.HasSuffix(key, "@odata.context") {
				continue
			}
			normalised[key] = normaliseSchema(nested)
		}
		return normalised
	case []any:
		normalised := make([]any, 0, len(typed))
		for _, item := range typed {
			normalised = append(normalised, normaliseSchema(item))
		}
		return normalised
	default:
		return value
	}
}

// checkDefaultSchemaRestoredInGraph asserts that destroy reset the job schema to the default recorded in the PreCheck.
func checkDefaultSchemaRestoredInGraph(fixture *synchronizationFixture) func(*terraform.State) error {
	return checkGraphSchema(fixture, func(schema map[string]any) error {
		expected := normaliseSchema(map[string]any{"synchronizationRules": fixture.DefaultRules, "directories": fixture.DefaultDirectories})
		actual := normaliseSchema(map[string]any{"synchronizationRules": schema["synchronizationRules"], "directories": schema["directories"]})
		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("expected destroy to reset the synchronization schema of job %s to the template default", fixture.JobId)
		}
		return nil
	})
}
