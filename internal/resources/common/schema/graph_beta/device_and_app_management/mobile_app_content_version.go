package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func MobileAppContentVersionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The current committed content version of the app, including its files.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the content version.",
			},
			"files": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The files within this content version.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":                         schema.StringAttribute{Computed: true},
						"size":                         schema.Int64Attribute{Computed: true},
						"size_encrypted":               schema.Int64Attribute{Computed: true},
						"upload_state":                 schema.StringAttribute{Computed: true},
						"is_committed":                 schema.BoolAttribute{Computed: true},
						"is_dependency":                schema.BoolAttribute{Computed: true},
						"azure_storage_uri":            schema.StringAttribute{Computed: true},
						"azure_storage_uri_expiration": schema.StringAttribute{Computed: true},
						"created_date_time":            schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}
