# Import an existing agent user into Terraform
# The import ID format is: {agent_user_id}[:hard_delete=true|false]
#
# Where:
# - {agent_user_id} is the unique identifier for the agent user
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_agents_agent_user.example "00000000-0000-0000-0000-000000000000"

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_agents_agent_user.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
