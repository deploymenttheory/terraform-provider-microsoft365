package graphBetaServicePrincipalSynchronizationSchemaJson_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Every test creates its own acc-test- custom application, service principal and SCIM synchronization job in its
// PreCheck, and a t.Cleanup deletes them by their exact IDs and purges them from deleted items. The job is never
// started, so no SCIM endpoint or credentials are needed. CheckDestroy asserts that destroy reset the job schema to
// the template default recorded before the test.

// checkImportedSchemaJSON asserts that an import reads non-empty JSON arrays for both schema collections.
func checkImportedSchemaJSON(states []*terraform.InstanceState) error {
	if len(states) != 1 {
		return fmt.Errorf("expected one imported instance, got %d", len(states))
	}
	for _, key := range []string{"synchronization_rules_json", "directories_json"} {
		var collection []any
		if err := json.Unmarshal([]byte(states[0].Attributes[key]), &collection); err != nil || len(collection) == 0 {
			return fmt.Errorf("expected imported %s to be a non-empty JSON array", key)
		}
	}
	return nil
}

// Lifecycle: replaces the template schema, updates it with a custom attribute and mapping, imports it and resets it to
// the template default on destroy.
func TestAccResourceServicePrincipalSynchronizationSchemaJson_01_Lifecycle(t *testing.T) {
	fixture := synchronizationJobPreCheck(t)
	resourceId := fixture.ServicePrincipalId + "/" + fixture.JobId
	minimalRules, minimalDirectories := minimalSchema(t, fixture)
	maximalRules, maximalDirectories := maximalSchema(t, fixture)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDefaultSchemaRestoredInGraph(fixture),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Replacing the job schema with the group object mapping disabled")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				ConfigVariables: config.Variables{
					"service_principal_id":       config.StringVariable(fixture.ServicePrincipalId),
					"job_id":                     config.StringVariable(fixture.JobId),
					"synchronization_rules_json": config.StringVariable(minimalRules),
					"directories_json":           config.StringVariable(minimalDirectories),
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue(resourceId),
					check.That(resourceType+".test").Key("service_principal_id").HasValue(fixture.ServicePrincipalId),
					check.That(resourceType+".test").Key("job_id").HasValue(fixture.JobId),
					check.That(resourceType+".test").Key("version").MatchesRegex(regexp.MustCompile(`.+`)),
					check.That(resourceType+".test").Key("synchronization_rules_json").HasValue(minimalRules),
					check.That(resourceType+".test").Key("directories_json").HasValue(minimalDirectories),
					checkGroupMappingEnabled(fixture, false),
					checkCustomAttribute(fixture, false),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adding a custom SCIM attribute and its attribute mapping")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				ConfigVariables: config.Variables{
					"service_principal_id":       config.StringVariable(fixture.ServicePrincipalId),
					"job_id":                     config.StringVariable(fixture.JobId),
					"synchronization_rules_json": config.StringVariable(maximalRules),
					"directories_json":           config.StringVariable(maximalDirectories),
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("synchronization_rules_json").HasValue(maximalRules),
					check.That(resourceType+".test").Key("directories_json").HasValue(maximalDirectories),
					checkGroupMappingEnabled(fixture, false),
					checkCustomAttribute(fixture, true),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the synchronization schema")
				},
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: resourceId,
				ConfigVariables: config.Variables{
					"service_principal_id":       config.StringVariable(fixture.ServicePrincipalId),
					"job_id":                     config.StringVariable(fixture.JobId),
					"synchronization_rules_json": config.StringVariable(maximalRules),
					"directories_json":           config.StringVariable(maximalDirectories),
				},
				ImportStateVerify: true,
				// Graph returns the collections with its own directory versions and Unsynchronized metadata. They are
				// semantically equal to the configuration but not the same string, so they are checked separately.
				ImportStateVerifyIgnore: []string{"timeouts", "synchronization_rules_json", "directories_json"},
				ImportStateCheck:        checkImportedSchemaJSON,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Returning to the minimal schema")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				ConfigVariables: config.Variables{
					"service_principal_id":       config.StringVariable(fixture.ServicePrincipalId),
					"job_id":                     config.StringVariable(fixture.JobId),
					"synchronization_rules_json": config.StringVariable(minimalRules),
					"directories_json":           config.StringVariable(minimalDirectories),
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("synchronization_rules_json").HasValue(minimalRules),
					checkCustomAttribute(fixture, false),
				),
			},
		},
	})
}
