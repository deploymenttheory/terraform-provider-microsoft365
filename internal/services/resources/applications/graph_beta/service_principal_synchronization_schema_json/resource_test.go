package graphBetaServicePrincipalSynchronizationSchemaJson_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalSynchronizationSchemaJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_synchronization_schema_json"
	schemaMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_synchronization_schema_json/mocks"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalSynchronizationSchemaJson.ResourceName

	// testResource is the test resource implementation for synchronization schemas
	testResource = graphBetaServicePrincipalSynchronizationSchemaJson.ServicePrincipalSynchronizationSchemaJsonTestResource{}

	mockResourceId = schemaMocks.ServicePrincipalID + "/" + schemaMocks.JobID
)

func setupMockEnvironment() (*mocks.Mocks, *schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	schemaMock := &schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock{}
	schemaMock.RegisterMocks()
	return mockClient, schemaMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	schemaMock := &schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock{}
	schemaMock.RegisterErrorMocks()
	return mockClient, schemaMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// checkLastPut asserts the body sent in the most recent PUT request.
func checkLastPut(mock *schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock, check func(body map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		puts := mock.PutRequests()
		if len(puts) == 0 {
			return fmt.Errorf("expected at least one PUT request")
		}
		return check(puts[len(puts)-1])
	}
}

// checkPutCount asserts the number of PUT requests sent so far.
func checkPutCount(mock *schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock, want int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := len(mock.PutRequests()); got != want {
			return fmt.Errorf("expected %d PUT requests, got %d", want, got)
		}
		return nil
	}
}

// objectMapping returns an object mapping of the first synchronization rule in a PUT body by name.
func objectMapping(body map[string]any, name string) (map[string]any, error) {
	rules, _ := body["synchronizationRules"].([]any)
	if len(rules) != 1 {
		return nil, fmt.Errorf("expected one synchronization rule in the PUT body, got %d", len(rules))
	}
	rule, _ := rules[0].(map[string]any)
	objectMappings, _ := rule["objectMappings"].([]any)
	for _, item := range objectMappings {
		mapping, _ := item.(map[string]any)
		if mapping["name"] == name {
			return mapping, nil
		}
	}
	return nil, fmt.Errorf("object mapping %q not found in the PUT body", name)
}

// checkDefaultSchemaRestored asserts that destroy reset the schema to the template default with DELETE.
func checkDefaultSchemaRestored(mock *schemaMocks.ServicePrincipalSynchronizationSchemaJsonMock) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if mock.DeleteCount() != 1 {
			return fmt.Errorf("expected one DELETE request on destroy, got %d", mock.DeleteCount())
		}
		if !mock.IsDefaultSchema() {
			return fmt.Errorf("expected the synchronization schema to be reset to the template default")
		}
		return nil
	}
}

// checkImportedSchema asserts that an import reads the schema and the collections Graph returns.
func checkImportedSchema(states []*terraform.InstanceState) error {
	if len(states) != 1 {
		return fmt.Errorf("expected one imported instance, got %d", len(states))
	}
	attributes := states[0].Attributes
	expected := map[string]string{
		"id":                   mockResourceId,
		"service_principal_id": schemaMocks.ServicePrincipalID,
		"job_id":               schemaMocks.JobID,
	}
	for key, want := range expected {
		if got := attributes[key]; got != want {
			return fmt.Errorf("expected imported %s to be %q, got %q", key, want, got)
		}
	}
	for _, key := range []string{"synchronization_rules_json", "directories_json"} {
		var collection []any
		if err := json.Unmarshal([]byte(attributes[key]), &collection); err != nil || len(collection) == 0 {
			return fmt.Errorf("expected imported %s to be a non-empty JSON array, got %q", key, attributes[key])
		}
	}
	if !regexp.MustCompile(`^Date:`).MatchString(attributes["version"]) {
		return fmt.Errorf("expected imported version to be set by the last PUT, got %q", attributes["version"])
	}
	return nil
}

// Minimal: replaces the template schema with one that disables the group object mapping, then imports it. Graph
// changes the directory versions and the Unsynchronized metadata, which must not show as drift.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkDefaultSchemaRestored(schemaMock),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue(mockResourceId),
					check.That(resourceType+".test").Key("service_principal_id").HasValue(schemaMocks.ServicePrincipalID),
					check.That(resourceType+".test").Key("job_id").HasValue(schemaMocks.JobID),
					check.That(resourceType+".test").Key("version").MatchesRegex(regexp.MustCompile(`^Date:`)),
					check.That(resourceType+".test").Key("synchronization_rules_json").MatchesRegex(regexp.MustCompile(`"enabled":false`)),
					check.That(resourceType+".test").Key("directories_json").MatchesRegex(regexp.MustCompile(`"version":"69F1E34D-A39F-4A48-B627-3AE0AD8E0DF6"`)),
					checkPutCount(schemaMock, 1),
					checkLastPut(schemaMock, func(body map[string]any) error {
						if len(body) != 2 {
							return fmt.Errorf("expected only synchronizationRules and directories in the PUT body, got %v", body)
						}
						groups, err := objectMapping(body, "Provision Microsoft Entra ID Groups")
						if err != nil {
							return err
						}
						if groups["enabled"] != false {
							return fmt.Errorf("expected the group object mapping to be disabled, got %v", groups["enabled"])
						}
						directories, _ := body["directories"].([]any)
						if len(directories) != 2 {
							return fmt.Errorf("expected two directories in the PUT body, got %d", len(directories))
						}
						return nil
					}),
				),
			},
			{
				ResourceName:     resourceType + ".test",
				ImportState:      true,
				ImportStateId:    mockResourceId,
				ImportStateCheck: checkImportedSchema,
			},
		},
	})
}

// Update: adds a custom target attribute and its mapping, then returns to the minimal schema.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_02_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkDefaultSchemaRestored(schemaMock),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 1),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 2),
					check.That(resourceType+".test").Key("synchronization_rules_json").MatchesRegex(regexp.MustCompile(`urn:ietf:params:scim:schemas:extension:acctest:2.0:User:costCenter`)),
					check.That(resourceType+".test").Key("directories_json").MatchesRegex(regexp.MustCompile(`urn:ietf:params:scim:schemas:extension:acctest:2.0:User:costCenter`)),
					checkLastPut(schemaMock, func(body map[string]any) error {
						users, err := objectMapping(body, "Provision Microsoft Entra ID Users")
						if err != nil {
							return err
						}
						attributeMappings, _ := users["attributeMappings"].([]any)
						if len(attributeMappings) != 3 {
							return fmt.Errorf("expected three user attribute mappings in the PUT body, got %d", len(attributeMappings))
						}
						return nil
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 3),
					checkLastPut(schemaMock, func(body map[string]any) error {
						users, err := objectMapping(body, "Provision Microsoft Entra ID Users")
						if err != nil {
							return err
						}
						attributeMappings, _ := users["attributeMappings"].([]any)
						if len(attributeMappings) != 2 {
							return fmt.Errorf("expected two user attribute mappings in the PUT body, got %d", len(attributeMappings))
						}
						return nil
					}),
				),
			},
		},
	})
}

// Reapply: re-applying the same JSON plans no change and sends no PUT. Rewriting the same schema as indented JSON
// is applied once and stored as configured, without an inconsistent result.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_03_Reapply(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 1),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 1),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal_heredoc.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 2),
					check.That(resourceType+".test").Key("synchronization_rules_json").MatchesRegex(regexp.MustCompile(`^\[\n  \{\n`)),
				),
			},
		},
	})
}

// Drift: a change made outside Terraform is detected on refresh and reverted by the next apply.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_04_Drift(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					schemaMock.SetObjectMappingEnabled("Provision Microsoft Entra ID Groups", true)
				},
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkPutCount(schemaMock, 2),
					checkLastPut(schemaMock, func(body map[string]any) error {
						groups, err := objectMapping(body, "Provision Microsoft Entra ID Groups")
						if err != nil {
							return err
						}
						if groups["enabled"] != false {
							return fmt.Errorf("expected the group object mapping to be disabled again, got %v", groups["enabled"])
						}
						return nil
					}),
				),
			},
		},
	})
}

// RemovedJob: a synchronization job removed outside Terraform drops the resource from state.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_05_RemovedJob(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					schemaMock.RemoveJob()
				},
				Config:             loadUnitTestTerraform("resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// MissingJob: a schema PUT to a job that does not exist surfaces the Graph error.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_06_MissingJob(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_missing_job.tf"),
				ExpectError: regexp.MustCompile(`does not exist under service principal`),
			},
		},
	})
}

// ImportInvalidID: an import ID that is not service_principal_id/job_id is rejected.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_07_ImportInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_minimal.tf"),
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: schemaMocks.ServicePrincipalID,
				ExpectError:   regexp.MustCompile(`Import ID must be in the format service_principal_id/job_id`),
			},
		},
	})
}

// Validation: invalid JSON, a JSON value that is not an array and an invalid service principal ID are rejected at plan.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_08_Validation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	cases := map[string]*regexp.Regexp{
		"resource_invalid_json.tf":                 regexp.MustCompile(`The value must be a JSON array`),
		"resource_invalid_not_array.tf":            regexp.MustCompile(`The value must be a JSON array, got a JSON object`),
		"resource_invalid_service_principal_id.tf": regexp.MustCompile(`must be a valid GUID`),
	}
	for filename, expectError := range cases {
		t.Run(filename, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      loadUnitTestTerraform(filename),
						ExpectError: expectError,
					},
				},
			})
		})
	}
	if len(schemaMock.PutRequests()) != 0 {
		t.Fatalf("expected no PUT requests for invalid configurations, got %d", len(schemaMock.PutRequests()))
	}
}

// ErrorHandling: Graph errors on PUT surface as diagnostics.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Bad Request - 400`),
			},
		},
	})
}

// ImportForbidden: a 403 on read during import reports the read permissions.
func TestUnitResourceServicePrincipalSynchronizationSchemaJson_10_ImportForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, schemaMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer schemaMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_minimal.tf"),
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: mockResourceId,
				ExpectError:   regexp.MustCompile(`Synchronization\.Read\.All`),
			},
		},
	})
}
