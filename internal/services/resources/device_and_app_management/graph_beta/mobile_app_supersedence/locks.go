package graphBetaMobileAppSupersedence

import (
	"context"
	"fmt"

	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
)

// lockAppRelationships also coordinates with supersedence writers for the same app.
func lockAppRelationships(ctx context.Context, sourceID string) (func(), error) {
	unlock, err := apprelationships.Lock(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("lock app supersedence relationships: %w", err)
	}
	return unlock, nil
}
