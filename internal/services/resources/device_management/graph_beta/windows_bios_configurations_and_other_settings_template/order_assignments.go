package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var errMissingFilterCreatedDateTime = errors.New("missing assignment filter createdDateTime")

// orderAssignmentsByFilterCreationTime orders filtered hardware assignments by filter creation time.
// The hardware endpoint rejects otherwise valid scope tags when distinct filters are not oldest first.
func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) orderAssignmentsByFilterCreationTime(
	ctx context.Context,
	body devicemanagement.HardwareConfigurationsItemAssignPostRequestBodyable,
) error {
	assignments := body.GetHardwareConfigurationAssignments()
	createdTimes := make(map[string]time.Time)
	positions := make([]int, 0)
	filtered := make([]graphmodels.HardwareConfigurationAssignmentable, 0)
	for index, assignment := range assignments {
		filterID := assignment.GetTarget().GetDeviceAndAppManagementAssignmentFilterId()
		if filterID == nil || *filterID == "" ||
			*filterID == "00000000-0000-0000-0000-000000000000" {
			continue
		}
		createdTimes[*filterID] = time.Time{}
		positions = append(positions, index)
		filtered = append(filtered, assignment)
	}
	if len(createdTimes) < 2 {
		return nil
	}

	for filterID := range createdTimes {
		filter, err := r.client.
			DeviceManagement().
			AssignmentFilters().
			ByDeviceAndAppManagementAssignmentFilterId(filterID).
			Get(ctx, nil)
		if err != nil {
			return fmt.Errorf("could not read assignment filter %s: %w", filterID, err)
		}
		if filter == nil || filter.GetCreatedDateTime() == nil ||
			filter.GetCreatedDateTime().IsZero() {
			return fmt.Errorf("%w: %s", errMissingFilterCreatedDateTime, filterID)
		}
		createdTimes[filterID] = *filter.GetCreatedDateTime()
	}

	sort.SliceStable(filtered, func(first, second int) bool {
		firstID := *filtered[first].GetTarget().GetDeviceAndAppManagementAssignmentFilterId()
		secondID := *filtered[second].GetTarget().GetDeviceAndAppManagementAssignmentFilterId()
		return createdTimes[firstID].Before(createdTimes[secondID])
	})
	for index, position := range positions {
		assignments[position] = filtered[index]
	}
	body.SetHardwareConfigurationAssignments(assignments)
	return nil
}
