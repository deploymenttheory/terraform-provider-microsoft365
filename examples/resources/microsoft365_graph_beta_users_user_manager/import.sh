#!/bin/bash

# Import an existing user manager relationship into Terraform
# The import ID is the user_id (the user whose manager is being managed)

terraform import microsoft365_graph_beta_users_user_manager.example 00000000-0000-0000-0000-000000000000

