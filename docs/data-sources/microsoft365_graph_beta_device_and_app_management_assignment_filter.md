---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: "Intune"
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

The Microsoft 365 Intune assignment filter data source provides information about a specific assignment filter.

## Example Usage

{{ tffile "examples/microsoft365_graph_beta/microsoft365_graph_beta_device_and_app_management_assignment_filter/datasource.tf" }}

{{ .SchemaMarkdown | trimspace }}