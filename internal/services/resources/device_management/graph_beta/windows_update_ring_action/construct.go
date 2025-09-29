package graphBetaWindowsUpdateRingAction

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ActionType represents the type of action to perform
type ActionType string

const (
	ActionPauseFeatureUpdates    ActionType = "pause_feature_updates"
	ActionResumeFeatureUpdates   ActionType = "resume_feature_updates"
	ActionExtendFeatureUpdates   ActionType = "extend_feature_updates_pause"
	ActionRollbackFeatureUpdates ActionType = "rollback_feature_updates"
	ActionPauseQualityUpdates    ActionType = "pause_quality_updates"
	ActionResumeQualityUpdates   ActionType = "resume_quality_updates"
	ActionRollbackQualityUpdates ActionType = "rollback_quality_updates"
)

// ActionRequest represents an action to be performed
type ActionRequest struct {
	ActionType ActionType
	Value      bool
}

// determineActionsToPerform analyzes the model to determine what actions need to be performed
func determineActionsToPerform(ctx context.Context, data *WindowsUpdateRingActionResourceModel) []ActionRequest {
	tflog.Debug(ctx, fmt.Sprintf("Determining actions to perform for %s", ResourceName))

	var actions []ActionRequest

	// Check feature update actions
	if !data.PauseFeatureUpdates.IsNull() && data.PauseFeatureUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionPauseFeatureUpdates,
			Value:      true,
		})
	}

	if !data.ResumeFeatureUpdates.IsNull() && data.ResumeFeatureUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionResumeFeatureUpdates,
			Value:      true,
		})
	}

	if !data.ExtendFeatureUpdatesPause.IsNull() && data.ExtendFeatureUpdatesPause.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionExtendFeatureUpdates,
			Value:      true,
		})
	}

	if !data.RollbackFeatureUpdates.IsNull() && data.RollbackFeatureUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionRollbackFeatureUpdates,
			Value:      true,
		})
	}

	// Check quality update actions
	if !data.PauseQualityUpdates.IsNull() && data.PauseQualityUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionPauseQualityUpdates,
			Value:      true,
		})
	}

	if !data.ResumeQualityUpdates.IsNull() && data.ResumeQualityUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionResumeQualityUpdates,
			Value:      true,
		})
	}

	if !data.RollbackQualityUpdates.IsNull() && data.RollbackQualityUpdates.ValueBool() {
		actions = append(actions, ActionRequest{
			ActionType: ActionRollbackQualityUpdates,
			Value:      true,
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d actions to perform", len(actions)))
	return actions
}

// constructPatchRequest creates a PATCH request body for pause/resume/rollback actions
func constructPatchRequest(ctx context.Context, actionType ActionType) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing PATCH request for action: %s", actionType))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()

	switch actionType {
	case ActionPauseFeatureUpdates:
		pauseValue := true
		requestBody.SetFeatureUpdatesPaused(&pauseValue)
	case ActionResumeFeatureUpdates:
		pauseValue := false
		requestBody.SetFeatureUpdatesPaused(&pauseValue)
	case ActionRollbackFeatureUpdates:
		rollbackValue := true
		requestBody.SetFeatureUpdatesWillBeRolledBack(&rollbackValue)
	case ActionPauseQualityUpdates:
		pauseValue := true
		requestBody.SetQualityUpdatesPaused(&pauseValue)
	case ActionResumeQualityUpdates:
		pauseValue := false
		requestBody.SetQualityUpdatesPaused(&pauseValue)
	case ActionRollbackQualityUpdates:
		rollbackValue := true
		requestBody.SetQualityUpdatesWillBeRolledBack(&rollbackValue)
	default:
		return nil, fmt.Errorf("unsupported PATCH action type: %s", actionType)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("PATCH JSON for action %s", actionType), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log PATCH object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// updateActionMetadata updates the metadata fields in the model after performing an action
func updateActionMetadata(ctx context.Context, data *WindowsUpdateRingActionResourceModel, actionType ActionType) {
	tflog.Debug(ctx, fmt.Sprintf("Updating action metadata for action: %s", actionType))

	now := time.Now().UTC().Format(time.RFC3339)
	data.LastActionPerformed = types.StringValue(fmt.Sprintf("%s", actionType))
	data.LastActionTimestamp = types.StringValue(now)

	tflog.Debug(ctx, fmt.Sprintf("Updated metadata: action=%s, timestamp=%s", actionType, now))
}
