<!-- BEGIN_TF_DOCS -->
### Requirements

No requirements.

### Providers

| Name | Version |
|------|---------|
| <a name="provider_azuread"></a> [azuread](#provider_azuread) | n/a |
| <a name="provider_microsoft365"></a> [microsoft365](#provider_microsoft365) | n/a |

### Modules

No modules.

### Resources

| Name | Type |
|------|------|
| [microsoft365_graph_beta_device_and_app_management_device_management_script.example](https://registry.terraform.io/providers/hashicorp/microsoft365/latest/docs/resources/graph_beta_device_and_app_management_device_management_script) | resource |
| [azuread_group.example_group](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/data-sources/group) | data source |

### Inputs

No inputs.

### Outputs

| Name | Description |
|------|-------------|
| <a name="output_existing_script_assignments"></a> [existing_script_assignments](#output_existing_script_assignments) | Assignments of the existing Device Management Script |
| <a name="output_existing_script_display_name"></a> [existing_script_display_name](#output_existing_script_display_name) | Display name of the existing Device Management Script |
| <a name="output_existing_script_group_assignments"></a> [existing_script_group_assignments](#output_existing_script_group_assignments) | Group assignments of the existing Device Management Script |
| <a name="output_existing_script_last_modified"></a> [existing_script_last_modified](#output_existing_script_last_modified) | Last modified date of the existing Device Management Script |
| <a name="output_new_script_id"></a> [new_script_id](#output_new_script_id) | ID of the newly created Device Management Script |
<!-- END_TF_DOCS -->