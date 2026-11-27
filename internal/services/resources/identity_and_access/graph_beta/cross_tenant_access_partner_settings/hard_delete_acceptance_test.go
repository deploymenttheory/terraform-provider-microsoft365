package graphBetaCrossTenantAccessPartnerSettings_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func TestAccResourceCrossTenantAccessPartnerSettings_PermanentDeletion(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC must be set to 1")
	}
	id := os.Getenv("M365_CROSS_TENANT_PARTNER_TEST_TENANT_ID")
	if id == "" {
		t.Skip("Set M365_CROSS_TENANT_PARTNER_TEST_TENANT_ID to a lab partner tenant with no existing active or deleted configuration")
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			client, err := acceptance.TestGraphClient()
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err = client.Policies().CrossTenantAccessPolicy().Partners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(id).Get(ctx, nil)
			if err == nil || grapherrors.GraphError(ctx, err).StatusCode != 404 {
				t.Fatal("Lab partner must be confirmed absent from the active store before testing")
			}
			_, err = client.Policies().DeletedItems().CrossTenantPartners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(id).Get(ctx, nil)
			if err == nil || grapherrors.GraphError(ctx, err).StatusCode != 404 {
				t.Fatal("Lab partner must be confirmed absent from policy deleted items before testing")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCrossTenantAccessPartnerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf("provider \"microsoft365\" {}\nresource %q \"hard_delete\" {\n tenant_id = %q\n hard_delete = true\n}", resourceType, id),
				Check:  resource.TestCheckResourceAttr(resourceType+".hard_delete", "tenant_id", id),
			},
			{
				ResourceName:            resourceType + ".hard_delete",
				ImportState:             true,
				ImportStateId:           id + ":hard_delete=true",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
