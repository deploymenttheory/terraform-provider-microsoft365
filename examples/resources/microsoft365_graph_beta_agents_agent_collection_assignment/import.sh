#!/bin/bash

# Import an existing Agent Collection Assignment using the composite ID
# Format: {agent_instance_id}/{agent_collection_id}

# The IDs can be found via the Graph API:
# GET https://graph.microsoft.com/beta/agentRegistry/agentInstances
# GET https://graph.microsoft.com/beta/agentRegistry/agentCollections

terraform import microsoft365_graph_beta_agents_agent_collection_assignment.example 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
