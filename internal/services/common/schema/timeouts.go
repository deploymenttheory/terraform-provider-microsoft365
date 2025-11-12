package schema

import (
	"context"

	datasourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Timeouts returns a common schema attribute for resource timeouts.
func Timeouts(ctx context.Context) resourceschema.Attribute {
	return resourcetimeouts.Attributes(ctx, resourcetimeouts.Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

// DatasourceTimeouts returns a common schema attribute for datasource timeouts.
func DatasourceTimeouts(ctx context.Context) datasourceschema.Attribute {
	return datasourcetimeouts.Attributes(ctx)
}
