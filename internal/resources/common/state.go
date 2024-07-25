package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ResourceWithID is an interface that represents a resource with an ID.
type ResourceWithID interface {
	GetTypeName() string
}

// StateWithID is an interface that represents a state model with an ID.
type StateWithID interface {
	GetID() string
}

// HandleReadStateError handles errors during the read operation.
func HandleReadStateError(ctx context.Context, resp *resource.ReadResponse, resource ResourceWithID, state StateWithID, err error) {
	if IsNotFoundError(err) || strings.Contains(err.Error(), "An error has occurred") {
		tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found, removing from state", resource.GetTypeName(), state.GetID()))
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("Could not update %s with ID %s: %s", resource.GetTypeName(), state.GetID(), err.Error()),
		)
	}
}

// HandleUpdateStateError handles errors during the update operation.
func HandleUpdateStateError(ctx context.Context, resp *resource.UpdateResponse, resource ResourceWithID, state StateWithID, err error) {
	if IsNotFoundError(err) || strings.Contains(err.Error(), "An error has occurred") {
		tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found, removing from state", resource.GetTypeName(), state.GetID()))
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.AddError(
			"Error updating resource",
			fmt.Sprintf("Could not update %s with ID %s: %s", resource.GetTypeName(), state.GetID(), err.Error()),
		)
	}
}
