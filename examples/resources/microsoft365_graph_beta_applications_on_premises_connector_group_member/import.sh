#!/bin/bash

# Import by connector ID. A connector belongs to exactly one connector group,
# so the connector group is read from Microsoft Graph.
terraform import microsoft365_graph_beta_applications_on_premises_connector_group_member.example "00000000-0000-0000-0000-000000000000"
