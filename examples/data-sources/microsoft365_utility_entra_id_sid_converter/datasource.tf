# Example usage of the Entra ID SID Converter data source
# This data source allows bidirectional conversion between SIDs and Object IDs

# Example 1: Convert a SID to Object ID
data "microsoft365_utility_entra_id_sid_converter" "sid_to_objectid" {
  sid = "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
}

# Output the converted Object ID
output "object_id_from_sid" {
  description = "The Object ID converted from the SID"
  value       = data.microsoft365_utility_entra_id_sid_converter.sid_to_objectid.object_id
  # Expected output: 73d664e4-0886-4a73-b745-c694da45ddb4
}

# Example 2: Convert an Object ID to SID
data "microsoft365_utility_entra_id_sid_converter" "objectid_to_sid" {
  object_id = "73d664e4-0886-4a73-b745-c694da45ddb4"
}

# Output the converted SID
output "sid_from_object_id" {
  description = "The SID converted from the Object ID"
  value       = data.microsoft365_utility_entra_id_sid_converter.objectid_to_sid.sid
  # Expected output: S-1-12-1-1943430372-1249052806-2496021943-3034400218
}

# Example 3: Use in a hybrid identity scenario
# Convert a synced user's SID to get their Entra ID Object ID

variable "on_prem_user_sid" {
  description = "The SID of the on-premises user"
  type        = string
  default     = "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
}

data "microsoft365_utility_entra_id_sid_converter" "synced_user" {
  sid = var.on_prem_user_sid
}

# Use the Object ID in a conditional access policy or other Entra ID resource
output "entra_object_id" {
  description = "The Entra ID Object ID for the synced user"
  value       = data.microsoft365_utility_entra_id_sid_converter.synced_user.object_id
}

# Example 4: Batch conversion using for_each
variable "user_sids" {
  description = "Map of user names to their on-premises SIDs"
  type        = map(string)
  default = {
    "john.doe"   = "S-1-12-1-1000000000-2000000000-3000000000-4000000000"
    "jane.smith" = "S-1-12-1-1234567890-987654321-1111111111-2222222222"
  }
}

data "microsoft365_utility_entra_id_sid_converter" "batch_conversion" {
  for_each = var.user_sids
  sid      = each.value
}

# Output all converted Object IDs
output "batch_object_ids" {
  description = "Map of usernames to their Entra ID Object IDs"
  value = {
    for name, converter in data.microsoft365_utility_entra_id_sid_converter.batch_conversion :
    name => converter.object_id
  }
}

