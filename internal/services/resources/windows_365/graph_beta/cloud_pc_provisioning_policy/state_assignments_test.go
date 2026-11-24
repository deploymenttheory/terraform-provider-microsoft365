package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

const (
	testAssignmentID     = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	testGroupID          = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	testServicePlanID    = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	testAllotmentName    = "example-allotment"
	testAllotmentLicence = int32(5)
)

// mappedAssignments returns the assignments mapped to Terraform as model structs.
func mappedAssignments(t *testing.T, assignments []models.CloudPcProvisioningPolicyAssignmentable) []CloudPcProvisioningPolicyAssignmentModel {
	t.Helper()
	set := MapAssignmentsToTerraformSet(context.Background(), assignments)
	require.False(t, set.IsNull(), "assignments set should not be null")
	var result []CloudPcProvisioningPolicyAssignmentModel
	require.False(t, set.ElementsAs(context.Background(), &result, false).HasError())
	return result
}

// Regression test for #4928: a typed group target on an assignment whose ID differs from the group ID.
func TestUnitMapAssignmentsToTerraformSet_TypedFrontlineGroupTarget(t *testing.T) {
	odataType := "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
	target := models.NewCloudPcManagementGroupAssignmentTarget()
	target.SetOdataType(&odataType)
	target.SetGroupId(new(testGroupID))
	target.SetServicePlanId(new(testServicePlanID))
	target.SetAllotmentDisplayName(new(testAllotmentName))
	target.SetAllotmentLicensesCount(new(testAllotmentLicence))

	assignment := models.NewCloudPcProvisioningPolicyAssignment()
	assignment.SetId(new(testAssignmentID))
	assignment.SetTarget(target)

	result := mappedAssignments(t, []models.CloudPcProvisioningPolicyAssignmentable{assignment})
	require.Len(t, result, 1)
	require.Equal(t, CloudPcProvisioningPolicyAssignmentModel{
		Type:                  types.StringValue("groupAssignmentTarget"),
		GroupId:               types.StringValue(testGroupID),
		ServicePlanId:         types.StringValue(testServicePlanID),
		AllotmentLicenseCount: types.Int32Value(testAllotmentLicence),
		AllotmentDisplayName:  types.StringValue(testAllotmentName),
	}, result[0])
}

// Regression test for #4928: the read/import path deserializes the Graph response with the SDK before mapping.
func TestUnitMapAssignmentsToTerraformSet_ReadFromGraphResponse(t *testing.T) {
	raw := `{
  "id": "dddddddd-dddd-4ddd-8ddd-dddddddddddd",
  "provisioningType": "sharedByEntraGroup",
  "assignments": [
    {
      "id": "` + testAssignmentID + `",
      "target": {
        "@odata.type": "#microsoft.graph.cloudPcManagementGroupAssignmentTarget",
        "groupId": "` + testGroupID + `",
        "servicePlanId": "` + testServicePlanID + `",
        "allotmentDisplayName": "` + testAllotmentName + `",
        "allotmentLicensesCount": 5
      }
    }
  ]
}`
	node, err := jsonserialization.NewJsonParseNode([]byte(raw))
	require.NoError(t, err)
	parsed, err := node.GetObjectValue(models.CreateCloudPcProvisioningPolicyFromDiscriminatorValue)
	require.NoError(t, err)
	policy, ok := parsed.(models.CloudPcProvisioningPolicyable)
	require.True(t, ok)

	result := mappedAssignments(t, policy.GetAssignments())
	require.Len(t, result, 1)
	require.Equal(t, testGroupID, result[0].GroupId.ValueString())
	require.Equal(t, testServicePlanID, result[0].ServicePlanId.ValueString())
	require.Equal(t, testAllotmentName, result[0].AllotmentDisplayName.ValueString())
	require.Equal(t, testAllotmentLicence, result[0].AllotmentLicenseCount.ValueInt32())
	require.Equal(t, "groupAssignmentTarget", result[0].Type.ValueString())
}

// A dedicated policy's group target carries no Frontline fields, which stay null.
func TestUnitMapAssignmentsToTerraformSet_TypedDedicatedGroupTarget(t *testing.T) {
	odataType := "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
	target := models.NewCloudPcManagementGroupAssignmentTarget()
	target.SetOdataType(&odataType)
	target.SetGroupId(new(testGroupID))

	assignment := models.NewCloudPcProvisioningPolicyAssignment()
	assignment.SetId(new(testAssignmentID))
	assignment.SetTarget(target)

	result := mappedAssignments(t, []models.CloudPcProvisioningPolicyAssignmentable{assignment})
	require.Len(t, result, 1)
	require.Equal(t, testGroupID, result[0].GroupId.ValueString())
	require.True(t, result[0].ServicePlanId.IsNull())
	require.True(t, result[0].AllotmentDisplayName.IsNull())
	require.True(t, result[0].AllotmentLicenseCount.IsNull())
}

// The SDK has no exclusion target model, so that target keeps its groupId in additionalData.
func TestUnitMapAssignmentsToTerraformSet_ExclusionGroupTarget(t *testing.T) {
	odataType := "#microsoft.graph.cloudPcManagementExclusionGroupAssignmentTarget"
	target := models.NewCloudPcManagementAssignmentTarget()
	target.SetOdataType(&odataType)
	target.GetAdditionalData()["groupId"] = testGroupID

	assignment := models.NewCloudPcProvisioningPolicyAssignment()
	assignment.SetId(new(testAssignmentID))
	assignment.SetTarget(target)

	result := mappedAssignments(t, []models.CloudPcProvisioningPolicyAssignmentable{assignment})
	require.Len(t, result, 1)
	require.Equal(t, "exclusionGroupAssignmentTarget", result[0].Type.ValueString())
	require.Equal(t, testGroupID, result[0].GroupId.ValueString())
}
