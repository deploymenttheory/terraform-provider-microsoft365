package schema

import (
	"context"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	datasourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ResourceTimeouts returns a common schema attribute for resource timeouts.
// Supports create, read, update, and delete operations with configurable timeout durations.
func ResourceTimeouts(ctx context.Context) resourceschema.Attribute {
	return resourcetimeouts.Attributes(ctx, resourcetimeouts.Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

// DatasourceTimeouts returns a common schema attribute for datasource timeouts.
// Supports read operation with configurable timeout duration.
func DatasourceTimeouts(ctx context.Context) datasourceschema.Attribute {
	return datasourcetimeouts.Attributes(ctx)
}

// ActionTimeouts returns a common schema attribute for action timeouts.
// Supports invoke operation with configurable timeout duration.
// Actions have a single timeout for the invoke operation, unlike resources which have separate timeouts for CRUD operations.
// Default timeout values should be set when retrieving the timeout in the Invoke function using data.Timeouts.Invoke(ctx, defaultDuration).
func ActionTimeouts(ctx context.Context) actionschema.Attribute {
	return actiontimeouts.Attributes(ctx)
}
