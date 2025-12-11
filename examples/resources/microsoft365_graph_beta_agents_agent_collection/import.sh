#!/bin/bash

# Import an existing Agent Collection using the Object ID (id)
# The ID can be found via the Graph API:
# GET https://graph.microsoft.com/beta/agentRegistry/agentCollections

# Note: Reserved collections (Global and Quarantined) cannot be managed via Terraform
# - Global: 00000000-0000-0000-0000-000000000001
# - Quarantined: 00000000-0000-0000-0000-000000000002

terraform import microsoft365_graph_beta_agents_agent_collection.example 00000000-0000-0000-0000-000000000000
