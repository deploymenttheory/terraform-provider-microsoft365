package graphBetaMacOSDepEnrollmentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type MacOSDepEnrollmentProfileTestResource struct{}

func (r MacOSDepEnrollmentProfileTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		profileResource := &MacOSDepEnrollmentProfileResource{client: client}
		depID, err := profileResource.resolveDepOnboardingSettingsId(ctx, types.StringValue(state.Attributes["dep_onboarding_settings_id"]))
		if err != nil {
			return err
		}
		_, err = client.
			DeviceManagement().
			DepOnboardingSettings().
			ByDepOnboardingSettingId(depID).
			EnrollmentProfiles().
			ByEnrollmentProfileId(state.ID).
			Get(ctx, nil)
		return err
	})
}
