#!/bin/bash

# Import an existing Agent Instance using the Object ID (id)
# The ID can be found in the Microsoft Entra admin center or via the Graph API:
# GET https://graph.microsoft.com/beta/agentRegistry/agentInstances

terraform import microsoft365_graph_beta_agents_agent_instance.example 00000000-0000-0000-0000-000000000000
