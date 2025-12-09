# Import an existing agent identity into Terraform
# The import ID format is: {agent_identity_id}/{agent_identity_blueprint_id}
#
# Where:
# - {agent_identity_id} is the Object ID of the agent identity service principal
# - {agent_identity_blueprint_id} is the Application (client) ID of the blueprint

terraform import microsoft365_graph_beta_agents_agent_identity.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"

