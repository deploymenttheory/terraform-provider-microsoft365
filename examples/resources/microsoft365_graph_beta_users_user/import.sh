#!/bin/bash

# Import an existing user into Terraform
# The import ID format is: {user_id}[:hard_delete=true|false]
#
# Where:
# - {user_id} is the unique identifier for the user
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_users_user.example 00000000-0000-0000-0000-000000000000

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_users_user.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
