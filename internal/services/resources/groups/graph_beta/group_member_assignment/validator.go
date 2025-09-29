package graphBetaGroupMemberAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ValidateGroupMemberAssignment is the main validation function that orchestrates all validation checks
func ValidateGroupMemberAssignment(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	object GroupMemberAssignmentResourceModel,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Starting validation for Group Member Assignment")

	memberObjectType := object.MemberObjectType.ValueString()
	memberId := object.MemberID.ValueString()
	groupId := object.GroupID.ValueString()

	tflog.Debug(ctx, "Validating member assignment", map[string]any{
		"group_id":           groupId,
		"member_id":          memberId,
		"member_object_type": memberObjectType,
	})

	// Get the target group to determine its type
	targetGroup, err := getTargetGroup(ctx, client, groupId)
	if err != nil {
		return err
	}

	// Determine the group type
	groupType := determineGroupType(targetGroup)
	tflog.Debug(ctx, fmt.Sprintf("Target group %s is of type: %s", groupId, groupType))

	// Check if member already exists in the group
	if err := validateMemberUniqueness(ctx, client, groupId, memberId, isUpdate); err != nil {
		return err
	}

	// Validate member compatibility with group type
	if err := validateMemberCompatibility(ctx, client, memberObjectType, groupType, memberId); err != nil {
		return err
	}

	tflog.Debug(ctx, "All validation checks passed for Group Member Assignment")
	return nil
}

// getTargetGroup retrieves the target group for validation
func getTargetGroup(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	groupId string,
) (graphmodels.Groupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving target group: %s", groupId))

	targetGroup, err := client.Groups().ByGroupId(groupId).Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve target group %s: %w", groupId, err)
	}

	return targetGroup, nil
}

// validateMemberUniqueness ensures that the member doesn't already exist in the group
func validateMemberUniqueness(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	groupId string,
	memberId string,
	isUpdate bool,
) error {
	tflog.Debug(ctx, "Validating member uniqueness constraint")

	// For updates, we don't need to check uniqueness as the member should already exist
	if isUpdate {
		tflog.Debug(ctx, "Update operation detected, skipping member uniqueness validation")
		return nil
	}

	memberExists, err := checkMemberExists(ctx, client, groupId, memberId)
	if err != nil {
		return fmt.Errorf("failed to check member existence: %w", err)
	}

	if memberExists {
		return fmt.Errorf("member %s already exists in group %s", memberId, groupId)
	}

	tflog.Debug(ctx, "Member uniqueness validation passed")
	return nil
}

// validateMemberCompatibility validates if a member type can be added to a specific group type
func validateMemberCompatibility(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	memberObjectType string,
	groupType string,
	memberId string,
) error {
	tflog.Debug(ctx, "Validating member compatibility constraint")

	switch groupType {
	case "Microsoft365":
		// Microsoft 365 groups only support Users
		if memberObjectType != "User" {
			return fmt.Errorf("microsoft 365 groups only support User members. Cannot add %s to Microsoft 365 group", memberObjectType)
		}

	case "Security":
		// Security groups support: User, Security groups, Device, Service principal, Organizational contact
		switch memberObjectType {
		case "User", "Device", "ServicePrincipal", "OrganizationalContact":
			// These are always allowed
			return nil
		case "Group":
			// For groups, we need to verify the member group is a Security group
			if err := validateSecurityGroupMember(ctx, client, memberId); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported member object type %s for Security groups", memberObjectType)
		}

	case "MailEnabledSecurity":
		return fmt.Errorf("mail-enabled security groups are read-only and do not support adding members")

	case "Distribution":
		return fmt.Errorf("distribution groups are read-only and do not support adding members")

	default:
		return fmt.Errorf("unknown or unsupported group type: %s", groupType)
	}

	tflog.Debug(ctx, "Member compatibility validation passed")
	return nil
}

// validateSecurityGroupMember validates that a group member is a Security group
func validateSecurityGroupMember(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	memberId string,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating that member group %s is a Security group", memberId))

	memberGroup, err := client.Groups().ByGroupId(memberId).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve member group %s: %w", memberId, err)
	}

	memberGroupType := determineGroupType(memberGroup)
	if memberGroupType != "Security" {
		return fmt.Errorf("only Security groups can be added as members to Security groups. Member group %s is of type %s", memberId, memberGroupType)
	}

	tflog.Debug(ctx, "Security group member validation passed")
	return nil
}

// checkMemberExists checks if a member already exists in a group
func checkMemberExists(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	groupId string,
	memberId string,
) (bool, error) {
	tflog.Debug(ctx, "Checking if member exists in group", map[string]any{
		"group_id":  groupId,
		"member_id": memberId,
	})

	// Get all members of the group and search for the specific member
	members, err := client.
		Groups().
		ByGroupId(groupId).
		Members().
		Get(ctx, nil)

	if err != nil {
		return false, fmt.Errorf("failed to retrieve group members: %w", err)
	}

	if members != nil && members.GetValue() != nil {
		for _, member := range members.GetValue() {
			if member.GetId() != nil && *member.GetId() == memberId {
				tflog.Debug(ctx, "Member found in group", map[string]any{
					"group_id":  groupId,
					"member_id": memberId,
				})
				return true, nil
			}
		}
	}

	tflog.Debug(ctx, "Member not found in group", map[string]any{
		"group_id":  groupId,
		"member_id": memberId,
	})
	return false, nil
}

// determineGroupType determines the type of group based on its properties
func determineGroupType(group graphmodels.Groupable) string {
	groupTypes := group.GetGroupTypes()
	mailEnabled := group.GetMailEnabled()
	securityEnabled := group.GetSecurityEnabled()

	// Check for Microsoft 365 group (Unified)
	for _, groupType := range groupTypes {
		if groupType == "Unified" {
			return "Microsoft365"
		}
	}

	// Check for other group types
	if mailEnabled != nil && securityEnabled != nil {
		if *mailEnabled && *securityEnabled {
			return "MailEnabledSecurity" // Read-only
		}
		if *mailEnabled && !*securityEnabled {
			return "Distribution" // Read-only
		}
		if !*mailEnabled && *securityEnabled {
			return "Security"
		}
	}

	return "Unknown"
}
