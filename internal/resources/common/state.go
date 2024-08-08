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

// HandleReadErrorIfNotFound handles errors during the read operation.
func HandleReadErrorIfNotFound(ctx context.Context, resp *resource.ReadResponse, resource ResourceWithID, state StateWithID, err error) {
	handleIsNotFoundStateError(ctx, "read", resource, state, err, func() {
		resp.State.RemoveResource(ctx)
	})
	resp.Diagnostics.AddError("Error reading resource", fmt.Sprintf("Could not read %s with ID %s: %s", resource.GetTypeName(), state.GetID(), err.Error()))
}

// HandleUpdateErrorIfNotFound handles errors during the update operation.
func HandleUpdateErrorIfNotFound(ctx context.Context, resp *resource.UpdateResponse, resource ResourceWithID, state StateWithID, err error) {
	handleIsNotFoundStateError(ctx, "update", resource, state, err, func() {
		resp.State.RemoveResource(ctx)
	})
	resp.Diagnostics.AddError("Error updating resource", fmt.Sprintf("Could not update %s with ID %s: %s", resource.GetTypeName(), state.GetID(), err.Error()))
}

// handleIsNotFoundStateError handles errors during CRUD operations internally.
func handleIsNotFoundStateError(ctx context.Context, operation string, resource ResourceWithID, state StateWithID, err error, removeResource func()) {
	if IsNotFoundError(err) || strings.Contains(err.Error(), "An error has occurred") {
		tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found, removing from state", resource.GetTypeName(), state.GetID()))
		removeResource()
	} else {
		fmt.Printf("Error %sing resource: %s with ID %s: %s", operation, resource.GetTypeName(), state.GetID(), err.Error())
	}
}
