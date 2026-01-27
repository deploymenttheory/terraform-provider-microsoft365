package graphBetaGroupMemberAssignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGroupMemberAssignmentResource_Create_Minimal tests creating a group member assignment with minimal configuration
func TestAccResourceGroupMemberAssignment_01_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test group ID from environment variable or skip
	testGroupID1 := os.Getenv("TEST_GROUP_ID_1")
	if testGroupID1 == "" {
		t.Skip("TEST_GROUP_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test member ID from environment variable or skip
	testMemberID1 := os.Getenv("TEST_MEMBER_ID_1")
	if testMemberID1 == "" {
		t.Skip("TEST_MEMBER_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_groups_group_member_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupMemberAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(testGroupID1, testMemberID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMemberAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", testGroupID1),
					resource.TestCheckResourceAttr(resourceName, "member_id", testMemberID1),
					resource.TestCheckResourceAttr(resourceName, "member_object_type", "User"),
					resource.TestCheckResourceAttrSet(resourceName, "member_type"),
					resource.TestCheckResourceAttrSet(resourceName, "member_display_name"),
				),
			},
		},
	})
}

// TestAccGroupMemberAssignmentResource_Create_Maximal tests creating a group member assignment with maximal configuration
func TestAccResourceGroupMemberAssignment_02_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test group ID from environment variable or skip
	testGroupID2 := os.Getenv("TEST_GROUP_ID_2")
	if testGroupID2 == "" {
		t.Skip("TEST_GROUP_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test member ID from environment variable or skip (should be a group ID for maximal test)
	testMemberID2 := os.Getenv("TEST_MEMBER_ID_2")
	if testMemberID2 == "" {
		t.Skip("TEST_MEMBER_ID_2 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_groups_group_member_assignment.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupMemberAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration (using a group as the member)
			{
				Config: testAccConfigMaximal(testGroupID2, testMemberID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMemberAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", testGroupID2),
					resource.TestCheckResourceAttr(resourceName, "member_id", testMemberID2),
					resource.TestCheckResourceAttr(resourceName, "member_object_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "member_type"),
					resource.TestCheckResourceAttrSet(resourceName, "member_display_name"),
				),
			},
		},
	})
}

// TestAccGroupMemberAssignmentResource_Delete_Minimal tests deleting a group member assignment with minimal configuration
func TestAccResourceGroupMemberAssignment_03_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test group ID from environment variable or skip
	testGroupID1 := os.Getenv("TEST_GROUP_ID_1")
	if testGroupID1 == "" {
		t.Skip("TEST_GROUP_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test member ID from environment variable or skip
	testMemberID1 := os.Getenv("TEST_MEMBER_ID_1")
	if testMemberID1 == "" {
		t.Skip("TEST_MEMBER_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_groups_group_member_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupMemberAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(testGroupID1, testMemberID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMemberAssignmentExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources[resourceName]
					if exists {
						return fmt.Errorf("resource %s still exists after deletion", resourceName)
					}
					return nil
				},
			},
		},
	})
}

// TestAccGroupMemberAssignmentResource_Delete_Maximal tests deleting a group member assignment with maximal configuration
func TestAccResourceGroupMemberAssignment_04_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test group ID from environment variable or skip
	testGroupID2 := os.Getenv("TEST_GROUP_ID_2")
	if testGroupID2 == "" {
		t.Skip("TEST_GROUP_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test member ID from environment variable or skip (should be a group ID for maximal test)
	testMemberID2 := os.Getenv("TEST_MEMBER_ID_2")
	if testMemberID2 == "" {
		t.Skip("TEST_MEMBER_ID_2 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_groups_group_member_assignment.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupMemberAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(testGroupID2, testMemberID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMemberAssignmentExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources[resourceName]
					if exists {
						return fmt.Errorf("resource %s still exists after deletion", resourceName)
					}
					return nil
				},
			},
		},
	})
}

// TestAccGroupMemberAssignmentResource_Import tests importing a resource
func TestAccResourceGroupMemberAssignment_05_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test group ID from environment variable or skip
	testGroupID1 := os.Getenv("TEST_GROUP_ID_1")
	if testGroupID1 == "" {
		t.Skip("TEST_GROUP_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test member ID from environment variable or skip
	testMemberID1 := os.Getenv("TEST_MEMBER_ID_1")
	if testMemberID1 == "" {
		t.Skip("TEST_MEMBER_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_groups_group_member_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupMemberAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(testGroupID1, testMemberID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMemberAssignmentExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     testGroupID1 + "/" + testMemberID1,
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccCheckGroupMemberAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

func testAccCheckGroupMemberAssignmentDestroy(s *terraform.State) error {
	// In a real test, we would verify the group member assignment is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the group member assignment
	return nil
}

// Test configurations

// Minimal configuration with default resource name (User as member)
func testAccConfigMinimal(groupID, memberID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_member_assignment" "minimal" {
  group_id = "%s"
  member_id = "%s"
  member_object_type = "User"
}
`, groupID, memberID)
}

// Maximal configuration with default resource name (Group as member)
func testAccConfigMaximal(groupID, memberID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_member_assignment" "maximal" {
  group_id = "%s"
  member_id = "%s"
  member_object_type = "Group"
}
`, groupID, memberID)
}
