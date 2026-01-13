package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
)

type GroupLicenseAssignmentTestResource struct{}

// Exists checks whether the specific group license assignment exists in Microsoft Graph.
//
// This implementation directly accesses Microsoft Graph SDK model fields using getter methods
// instead of using the shared CheckResourceExistsByArrayMembership helper. This is necessary
// because msgraph SDK models (Userable, Groupable, etc.) use private fields with getter methods
// rather than exported struct fields. When json.Marshal is called on these SDK models, it returns
// empty objects {} because the JSON encoder only serializes exported fields. The shared helper
// relies on JSON marshaling to navigate nested structures, which fails for SDK models, causing
// the assignedLicenses field to appear as "not found" even when it exists. By directly accessing
// group.GetAssignedLicenses() and license.GetSkuId(), we bypass JSON serialization entirely and
// work directly with the SDK's type-safe getter methods.
func (r GroupLicenseAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	groupId := state.Attributes["group_id"]
	skuId := state.Attributes["sku_id"]

	if groupId == "" {
		return nil, fmt.Errorf("group_id not found in state")
	}
	if skuId == "" {
		return nil, fmt.Errorf("sku_id not found in state")
	}

	group, err := graphClient.Groups().ByGroupId(groupId).Get(ctx, &groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "assignedLicenses"},
		},
	})
	if err != nil {
		errMsg := err.Error()
		if errMsg != "" && (strings.Contains(errMsg, "does not exist") ||
			strings.Contains(errMsg, "ResourceNotFound") ||
			strings.Contains(errMsg, "404")) {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	assignedLicenses := group.GetAssignedLicenses()
	for _, license := range assignedLicenses {
		if license == nil {
			continue
		}
		licenseSkuId := license.GetSkuId()
		if licenseSkuId != nil && licenseSkuId.String() == skuId {
			exists := true
			return &exists, nil
		}
	}

	exists := false
	return &exists, nil
}
