package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func IntuneApplicationAssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The Intune Application Assignment configuration for managing applications deployments in Microsoft 365.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Key of the entity. This is read-only and automatically generated.",
				Computed:            true,
			},
		},
	}
}
