#!/bin/bash

# Import an existing group into Terraform
# The import ID format is: {group_id}[:hard_delete=true|false]
#
# Where:
# - {group_id} is the unique identifier for the group
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_groups_group.example 00000000-0000-0000-0000-000000000000

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_groups_group.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
