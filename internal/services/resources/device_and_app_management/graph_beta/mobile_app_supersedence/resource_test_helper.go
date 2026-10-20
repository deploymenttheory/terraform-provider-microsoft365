package graphBetaMobileAppSupersedence

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// MobileAppSupersedenceTestResource implements the shared acceptance test interface.
type MobileAppSupersedenceTestResource struct{}

// Exists verifies the relationship through the same collection the service supports.
func (r MobileAppSupersedenceTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	client, err := acceptance.TestGraphClient() //nolint:contextcheck // The shared acceptance client factory does not accept a context.
	if err != nil {
		return nil, fmt.Errorf("create acceptance client: %w", err)
	}
	sourceID := state.Attributes["source_id"]
	if sourceID == "" {
		sourceID = apprelationships.SourceID(state.ID)
	}
	relationships, err := apprelationships.List(ctx, client, sourceID)
	if err == nil {
		_, err = apprelationships.Find(relationships, state.ID, "")
	}
	found := err == nil
	if err != nil && !errors.Is(err, apprelationships.ErrNotFound) &&
		grapherrors.GraphError(ctx, err).StatusCode != 404 {
		return nil, fmt.Errorf("check app supersedence: %w", err)
	}
	return &found, nil
}
