<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_microsoft365"></a> [microsoft365](#requirement\_microsoft365) | ~> 1.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_microsoft365"></a> [microsoft365](#provider\_microsoft365) | ~> 1.0.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [microsoft365_graph_beta_device_and_app_management_assignment_filter.example](https://registry.terraform.io/providers/deploymenttheory/terraform-provider-microsoft365/latest/docs/resources/graph_beta_device_and_app_management_assignment_filter) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_client_id"></a> [client\_id](#input\_client\_id) | The client ID for the Entra ID application | `string` | n/a | yes |
| <a name="input_client_secret"></a> [client\_secret](#input\_client\_secret) | The client secret for the Entra ID application | `string` | n/a | yes |
| <a name="input_cloud"></a> [cloud](#input\_cloud) | The Microsoft cloud environment to use | `string` | `"public"` | no |
| <a name="input_tenant_id"></a> [tenant\_id](#input\_tenant\_id) | The Microsoft 365 tenant ID | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END_TF_DOCS -->