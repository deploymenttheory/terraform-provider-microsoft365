# Import an existing agent identity blueprint service principal
# The import ID format is: {service_principal_id}[:hard_delete=true|false]
#
# Where:
# - {service_principal_id} is the service principal's object ID
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example "00000000-0000-0000-0000-000000000000"

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
