#!/bin/bash

# Import an existing Agent Identity Blueprint using the Object ID (id)
# The ID can be found in the Microsoft Entra admin center under:
# Applications > App registrations > [Your Blueprint] > Overview > Object ID

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint.example 00000000-0000-0000-0000-000000000000

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
