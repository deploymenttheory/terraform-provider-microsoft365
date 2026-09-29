package graphBetaIdentityAndAccessB2bManagementPolicyAssignment_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy"
	assignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy_assignment"
)

var testResource = assignment.B2bManagementPolicyAssignmentTestResource{}

func testAccConfig(name string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + name)
	if err != nil {
		panic(fmt.Sprintf("failed to load B2B assignment config: %s", err))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceB2bManagementPolicyAssignment_01_ServicePrincipal(t *testing.T) {
	var assignmentID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: policy.ResourceName,
				TestResource: policy.B2bManagementPolicyTestResource{},
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 01_service_principal") },
				Config:    testAccConfig("resource_01_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".service_principal").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_service_principal.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".service_principal"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".service_principal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: testAccConfig("resource_01_service_principal.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceB2bManagementPolicyAssignment_02_Application(t *testing.T) {
	var assignmentID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: policy.ResourceName,
				TestResource: policy.B2bManagementPolicyTestResource{},
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 02_application") },
				Config:    testAccConfig("resource_02_application.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".application").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_application.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".application"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".application",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: testAccConfig("resource_02_application.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceB2bManagementPolicyAssignment_03_MinimalToMaximal(t *testing.T) {
	var assignmentID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: policy.ResourceName,
				TestResource: policy.B2bManagementPolicyTestResource{},
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 01_service_principal") },
				Config:    testAccConfig("resource_01_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".service_principal").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_service_principal.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".service_principal"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 03_both_targets") },
				Config:    testAccConfig("resource_03_both_targets.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".service_principal").ExistsInGraph(testResource),
					check.That(resourceType+".application").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_service_principal.target",
						"microsoft365_graph_beta_applications_application.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".service_principal"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".service_principal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				ResourceName:            resourceType + ".application",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: testAccConfig("resource_03_both_targets.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceB2bManagementPolicyAssignment_04_MaximalToMinimal(t *testing.T) {
	var assignmentID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: policy.ResourceName,
				TestResource: policy.B2bManagementPolicyTestResource{},
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 03_both_targets") },
				Config:    testAccConfig("resource_03_both_targets.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".service_principal").ExistsInGraph(testResource),
					check.That(resourceType+".application").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_service_principal.target",
						"microsoft365_graph_beta_applications_application.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".service_principal"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply assignments: 01_service_principal") },
				Config:    testAccConfig("resource_01_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(
							"B2B management policy assignment",
							30*time.Second,
						)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".service_principal").ExistsInGraph(testResource),
					checkAssignmentTargets(
						"microsoft365_graph_beta_applications_service_principal.target",
					),
					func(state *terraform.State) error {
						id := state.RootModule().Resources[resourceType+".service_principal"].Primary.ID
						if assignmentID != "" && assignmentID != id {
							return fmt.Errorf(
								"unchanged assignment was recreated: %s to %s",
								assignmentID,
								id,
							)
						}
						assignmentID = id
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".service_principal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: testAccConfig("resource_01_service_principal.tf"), PlanOnly: true},
		},
	})
}

// checkAssignmentTargets checks exact Graph membership, including removal of targets no longer configured.
func checkAssignmentTargets(addresses ...string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return err
		}
		expected := make(map[string]bool, len(addresses))
		for _, address := range addresses {
			expected[state.RootModule().Resources[address].Primary.ID] = true
		}
		policyID := state.RootModule().Resources[policy.ResourceName+".test"].Primary.ID
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		return crud.PollUntil(ctx, 3*time.Second, func(ctx context.Context) (bool, error) {
			builder := client.Policies().
				B2bManagementPolicies().
				ByB2bManagementPolicyId(policyID).
				AppliesTo()
			page, err := builder.Get(ctx, nil)
			if err != nil {
				return false, err
			}
			actual := make(map[string]bool)
			for page != nil {
				for _, object := range page.GetValue() {
					if object.GetId() != nil {
						actual[*object.GetId()] = true
					}
				}
				next := page.GetOdataNextLink()
				if next == nil || *next == "" {
					break
				}
				page, err = builder.WithUrl(*next).Get(ctx, nil)
				if err != nil {
					return false, err
				}
			}
			if len(actual) != len(expected) {
				return false, fmt.Errorf(
					"expected %d assignment targets, got %d",
					len(expected),
					len(actual),
				)
			}
			for id := range expected {
				if !actual[id] {
					return false, fmt.Errorf("expected assignment target %s not found", id)
				}
			}
			return true, nil
		})
	}
}
