package graphBetaCrossTenantAccessPolicy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestUnitResourceCrossTenantAccessPolicy_RefreshDestroyFlag(t *testing.T) {
	for _, flag := range []types.Bool{types.BoolNull(), types.BoolUnknown(), types.BoolValue(false), types.BoolValue(true)} {
		data := CrossTenantAccessPolicyResourceModel{RestoreDefaultsOnDestroy: flag}
		MapRemoteResourceStateToTerraform(
			t.Context(),
			&data,
			graphmodels.NewCrossTenantAccessPolicy(),
		)
		if data.RestoreDefaultsOnDestroy.IsNull() || data.RestoreDefaultsOnDestroy.IsUnknown() {
			t.Fatal("refresh must initialise a missing Terraform-only destroy flag")
		}
		if data.RestoreDefaultsOnDestroy.ValueBool() != flag.ValueBool() {
			t.Fatal("refresh must retain an explicitly configured destroy flag")
		}
	}
}
