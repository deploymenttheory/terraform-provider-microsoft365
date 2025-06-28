package graphBetaCloudPcAuditEvent

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.CloudPcAuditEventable) CloudPcAuditEventItemModel {
	model := CloudPcAuditEventItemModel{
		ID:               convert.GraphToFrameworkString(data.GetId()),
		DisplayName:      convert.GraphToFrameworkString(data.GetDisplayName()),
		ComponentName:    convert.GraphToFrameworkString(data.GetComponentName()),
		Activity:         convert.GraphToFrameworkString(data.GetActivity()),
		ActivityDateTime: convert.GraphToFrameworkTime(data.GetActivityDateTime()),
		ActivityType:     convert.GraphToFrameworkString(data.GetActivityType()),
		CorrelationId:    convert.GraphToFrameworkString(data.GetCorrelationId()),
	}

	// Enum fields: ActivityOperationType, ActivityResult, Category
	model.ActivityOperationType = convert.GraphToFrameworkEnum(data.GetActivityOperationType())
	model.ActivityResult = convert.GraphToFrameworkEnum(data.GetActivityResult())
	model.Category = convert.GraphToFrameworkEnum(data.GetCategory())

	if data.GetActor() != nil {
		model.Actor = mapActorToModel(data.GetActor())
	}

	if data.GetResources() != nil {
		for _, r := range data.GetResources() {
			model.Resources = append(model.Resources, mapResourceToModel(r))
		}
	}

	return model
}

func mapActorToModel(actor graphmodels.CloudPcAuditActorable) *CloudPcAuditActorModel {
	if actor == nil {
		return nil
	}
	m := &CloudPcAuditActorModel{
		ApplicationDisplayName: convert.GraphToFrameworkString(actor.GetApplicationDisplayName()),
		ApplicationId:          convert.GraphToFrameworkString(actor.GetApplicationId()),
		IpAddress:              convert.GraphToFrameworkString(actor.GetIpAddress()),
		RemoteTenantId:         convert.GraphToFrameworkString(actor.GetRemoteTenantId()),
		RemoteUserId:           convert.GraphToFrameworkString(actor.GetRemoteUserId()),
		ServicePrincipalName:   convert.GraphToFrameworkString(actor.GetServicePrincipalName()),
		UserId:                 convert.GraphToFrameworkString(actor.GetUserId()),
		UserPrincipalName:      convert.GraphToFrameworkString(actor.GetUserPrincipalName()),
	}
	// Enum: Type
	m.Type = convert.GraphToFrameworkEnum(actor.GetTypeEscaped())
	if perms := actor.GetUserPermissions(); perms != nil {
		for _, p := range perms {
			m.UserPermissions = append(m.UserPermissions, convert.GraphToFrameworkString(&p))
		}
	}
	if tags := actor.GetUserRoleScopeTags(); tags != nil {
		for _, t := range tags {
			m.UserRoleScopeTags = append(m.UserRoleScopeTags, mapUserRoleScopeTagToModel(t))
		}
	}
	return m
}

func mapUserRoleScopeTagToModel(tag graphmodels.CloudPcUserRoleScopeTagInfoable) CloudPcUserRoleScopeTagInfoModel {
	return CloudPcUserRoleScopeTagInfoModel{
		DisplayName:    convert.GraphToFrameworkString(tag.GetDisplayName()),
		RoleScopeTagId: convert.GraphToFrameworkString(tag.GetRoleScopeTagId()),
	}
}

func mapResourceToModel(resource graphmodels.CloudPcAuditResourceable) CloudPcAuditResourceModel {
	m := CloudPcAuditResourceModel{
		DisplayName:  convert.GraphToFrameworkString(resource.GetDisplayName()),
		ResourceId:   convert.GraphToFrameworkString(resource.GetResourceId()),
		ResourceType: convert.GraphToFrameworkString(resource.GetResourceType()),
	}
	if props := resource.GetModifiedProperties(); props != nil {
		for _, p := range props {
			m.ModifiedProperties = append(m.ModifiedProperties, mapPropertyToModel(p))
		}
	}
	return m
}

func mapPropertyToModel(prop graphmodels.CloudPcAuditPropertyable) CloudPcAuditPropertyModel {
	return CloudPcAuditPropertyModel{
		DisplayName: convert.GraphToFrameworkString(prop.GetDisplayName()),
		NewValue:    convert.GraphToFrameworkString(prop.GetNewValue()),
		OldValue:    convert.GraphToFrameworkString(prop.GetOldValue()),
	}
}
