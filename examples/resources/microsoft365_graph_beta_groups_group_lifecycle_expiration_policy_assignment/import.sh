#!/bin/bash

# Import a group lifecycle policy assignment using the group ID
# The ID for this resource is simply the group ID (UUID format)

terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment.marketing "12345678-1234-1234-1234-123456789012"

# Note: The group must already be assigned to the tenant's lifecycle policy
# before it can be imported into Terraform state

