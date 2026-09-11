#!/bin/sh
# Replace {policy_id}/{rule_id} with the resource UUIDs.
terraform import microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule.example "{policy_id}/{rule_id}"
