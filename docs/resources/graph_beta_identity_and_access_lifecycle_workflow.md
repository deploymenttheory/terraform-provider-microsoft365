---
page_title: "microsoft365_graph_beta_identity_and_access_lifecycle_workflow Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft Entra ID Governance lifecycle workflows using the /identityGovernance/lifecycleWorkflows/workflows endpoint. Lifecycle workflows automate joiner, mover and leaver tasks for users. Changes to display_name, description, is_enabled and is_scheduling_enabled are applied with PATCH; changes to category, execution_conditions or tasks create a new workflow version in place (createNewVersion). On destroy the workflow is deleted and then permanently removed from the lifecycle workflows deleted items container.
---

# microsoft365_graph_beta_identity_and_access_lifecycle_workflow (Resource)

Manages Microsoft Entra ID Governance lifecycle workflows using the `/identityGovernance/lifecycleWorkflows/workflows` endpoint. Lifecycle workflows automate joiner, mover and leaver tasks for users. Changes to `display_name`, `description`, `is_enabled` and `is_scheduling_enabled` are applied with PATCH; changes to `category`, `execution_conditions` or `tasks` create a new workflow version in place (`createNewVersion`). On destroy the workflow is deleted and then permanently removed from the lifecycle workflows deleted items container.

## Microsoft Documentation

- [Lifecycle workflows documentation](https://learn.microsoft.com/en-us/entra/id-governance/what-are-lifecycle-workflows)
- [workflow resource type](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflow?view=graph-rest-beta)
- [Create workflow](https://learn.microsoft.com/en-us/graph/api/identitygovernance-lifecycleworkflowscontainer-post-workflows?view=graph-rest-beta)
- [Get workflow](https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-get?view=graph-rest-beta)
- [Update workflow](https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-update?view=graph-rest-beta)
- [workflow: createNewVersion](https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-createnewversion?view=graph-rest-beta)
- [Delete workflow](https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-delete?view=graph-rest-beta)
- [Permanently delete a deleted workflow](https://learn.microsoft.com/en-us/graph/api/identitygovernance-deleteditemcontainer-delete?view=graph-rest-beta)
- [Lifecycle workflow built-in tasks](https://learn.microsoft.com/en-us/entra/id-governance/lifecycle-workflow-tasks)
- [List taskDefinitions](https://learn.microsoft.com/en-us/graph/api/identitygovernance-lifecycleworkflowscontainer-list-taskdefinitions?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `LifecycleWorkflows-Workflow.Read.All`
- `LifecycleWorkflows-Workflow.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

`LifecycleWorkflows-Workflow.ReadBasic.All` is listed as the least privileged read permission in Graph metadata, but it only allows listing workflows; reading a workflow's tasks needs `LifecycleWorkflows-Workflow.Read.All`.

Lifecycle workflows require a Microsoft Entra ID Governance license in the tenant.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Preview | Initial release |

## Example Usage

### Joiner: pre-hire onboarding

```terraform
# Joiner workflow that runs 7 days before the employee hire date
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "joiner" {
  category              = "joiner"
  display_name          = "Pre-hire onboarding"
  description           = "Enable the account and send a welcome email before day one"
  is_enabled            = false
  is_scheduling_enabled = false

  execution_conditions = {
    scope = {
      rule = "(department eq 'Marketing')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeHireDate"
        offset_in_days = -7
      }
    }
  }

  tasks = [
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
      arguments = [
        {
          name  = "enableOnPremisesAccount"
          value = "false"
        }
      ]
    },
    {
      task_definition_id = "70b29d51-b59a-4773-9280-8841dfd3f2ea"
      display_name       = "Send welcome email"
      continue_on_error  = true
    }
  ]
}
```

### Mover: attribute change

```terraform
# Mover workflow triggered by a department change
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "mover" {
  category     = "mover"
  display_name = "Department change"
  description  = "Notify the manager when a user changes department"

  execution_conditions = {
    scope = {
      rule = "(department eq 'Marketing')"
    }
    trigger = {
      attribute_change = {
        trigger_attribute = "department"
      }
    }
  }

  tasks = [
    {
      task_definition_id = "aab41899-9972-422a-9d97-f626014578b7"
      display_name       = "Send email to notify manager of user move"
    }
  ]
}
```

### Mover: group membership change

```terraform
# Mover workflow scoped to members of a group, triggered when a user is added to it
resource "microsoft365_graph_beta_groups_group" "sales_team" {
  display_name     = "Sales Team"
  mail_nickname    = "sales-team"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "group_membership" {
  category              = "mover"
  display_name          = "Sales team onboarding"
  description           = "Runs when a user is added to the Sales Team group"
  is_enabled            = true
  is_scheduling_enabled = false

  execution_conditions = {
    scope = {
      group_ids = [microsoft365_graph_beta_groups_group.sales_team.id]
    }
    trigger = {
      membership_change = {
        change_type = "add"
      }
    }
  }

  tasks = [
    {
      task_definition_id = "22085229-5809-45e8-97fd-270d28d66910"
      display_name       = "Add user to groups"
      description        = "Add the user to the sales resources group"
      arguments = [
        {
          name  = "groupID"
          value = "00000000-0000-0000-0000-000000000000"
        }
      ]
    },
    {
      task_definition_id = "aab41899-9972-422a-9d97-f626014578b7"
      display_name       = "Send email to notify manager of user move"
      continue_on_error  = true
      arguments = [
        {
          name  = "locale"
          value = "en-us"
        }
      ]
    }
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
```

### Leaver: last day

```terraform
# Leaver workflow that runs on the user's last day
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "leaver" {
  category     = "leaver"
  display_name = "Offboard on last day"
  description  = "Disable the account and remove licenses on the last day"

  execution_conditions = {
    scope = {
      rule = "(department eq 'Marketing')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeLeaveDateTime"
        offset_in_days = 0
      }
    }
  }

  tasks = [
    {
      task_definition_id = "1dfdfcc7-52fa-4c2e-bf3a-e3919cc12950"
      display_name       = "Disable user account"
      arguments = [
        {
          name  = "disableOnPremisesAccount"
          value = "false"
        }
      ]
    },
    {
      task_definition_id = "8fa97d28-3e52-4985-b3a9-a1126f9b8b4e"
      display_name       = "Remove all licenses for user"
      continue_on_error  = true
    },
    {
      task_definition_id = "509589a4-0466-4471-829e-49c5e502bdee"
      display_name       = "Revoke all refresh tokens for user"
    }
  ]
}
```

### On-demand workflow

```terraform
# Minimal on-demand workflow with a single task and only required attributes
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "minimal" {
  category     = "leaver"
  display_name = "Revoke sessions on demand"

  execution_conditions = {
    on_demand_execution_only = true
  }

  tasks = [
    {
      task_definition_id = "509589a4-0466-4471-829e-49c5e502bdee"
      display_name       = "Revoke all refresh tokens for user"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `category` (String) The category of the workflow. Possible values are: `joiner`, `mover`, `leaver`. Changing the category creates a new workflow version.
- `display_name` (String) The display name of the workflow.
- `execution_conditions` (Attributes) The conditions that determine when and for whom the workflow runs. Set `on_demand_execution_only = true` for a workflow that only runs on demand, or set both `scope` and `trigger` for a trigger and scope based workflow. (see [below for nested schema](#nestedatt--execution_conditions))
- `tasks` (Attributes List) The built-in tasks the workflow runs, in execution order. Custom task extensions are not supported by this resource. Changing tasks creates a new workflow version. (see [below for nested schema](#nestedatt--tasks))

### Optional

- `description` (String) A description of the workflow. Defaults to an empty string, which is what Graph returns when no description is set.
- `is_enabled` (Boolean) Whether the workflow is enabled. Defaults to `false`.
- `is_scheduling_enabled` (Boolean) Whether the workflow runs on its schedule. When `true`, the workflow processes in-scope users every schedule interval. Defaults to `false`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) When the workflow was created.
- `id` (String) The unique identifier of the workflow.
- `last_modified_date_time` (String) When the workflow was last modified.
- `version` (Number) The current version number of the workflow. Increases each time `category`, `execution_conditions` or `tasks` change.

<a id="nestedatt--execution_conditions"></a>
### Nested Schema for `execution_conditions`

Optional:

- `on_demand_execution_only` (Boolean) When `true`, the workflow can only be run on demand (`onDemandExecutionOnly`) and `scope` and `trigger` must not be set. Defaults to `false`.
- `scope` (Attributes) The users the workflow applies to. Set exactly one of `rule` or `group_ids`. (see [below for nested schema](#nestedatt--execution_conditions--scope))
- `trigger` (Attributes) The event that starts the workflow. Set exactly one of `time_based_attribute`, `attribute_change` or `membership_change`. (see [below for nested schema](#nestedatt--execution_conditions--trigger))

<a id="nestedatt--execution_conditions--scope"></a>
### Nested Schema for `execution_conditions.scope`

Optional:

- `group_ids` (Set of String) The IDs of the groups whose members are in scope (`groupBasedSubjectSet`). Required by a `membership_change` trigger.
- `rule` (String) A filter rule that selects users (`ruleBasedSubjectSet`), for example `(department eq 'Marketing')`.


<a id="nestedatt--execution_conditions--trigger"></a>
### Nested Schema for `execution_conditions.trigger`

Optional:

- `attribute_change` (Attributes) Runs the workflow when a user attribute changes (`attributeChangeTrigger`). Typically used by `mover` workflows. (see [below for nested schema](#nestedatt--execution_conditions--trigger--attribute_change))
- `membership_change` (Attributes) Runs the workflow when a user is added to or removed from a group in `scope.group_ids` (`membershipChangeTrigger`). (see [below for nested schema](#nestedatt--execution_conditions--trigger--membership_change))
- `time_based_attribute` (Attributes) Runs the workflow relative to a date attribute of the user (`timeBasedAttributeTrigger`). (see [below for nested schema](#nestedatt--execution_conditions--trigger--time_based_attribute))

<a id="nestedatt--execution_conditions--trigger--attribute_change"></a>
### Nested Schema for `execution_conditions.trigger.attribute_change`

Required:

- `trigger_attribute` (String) The name of the user attribute whose change triggers the workflow, for example `department` or `jobTitle`. Graph accepts a single trigger attribute.


<a id="nestedatt--execution_conditions--trigger--membership_change"></a>
### Nested Schema for `execution_conditions.trigger.membership_change`

Required:

- `change_type` (String) The membership change that triggers the workflow. Possible values are: `add`, `remove`.


<a id="nestedatt--execution_conditions--trigger--time_based_attribute"></a>
### Nested Schema for `execution_conditions.trigger.time_based_attribute`

Required:

- `attribute` (String) The user date attribute. Possible values are: `employeeHireDate`, `employeeLeaveDateTime`, `createdDateTime`.
- `offset_in_days` (Number) The number of days before (negative) or after (positive) the attribute date to run the workflow. Must be between -180 and 180.




<a id="nestedatt--tasks"></a>
### Nested Schema for `tasks`

Required:

- `display_name` (String) The display name of the task.
- `task_definition_id` (String) The ID of the built-in task definition, from `GET /identityGovernance/lifecycleWorkflows/taskDefinitions`. For example `6fc52c9d-398b-4305-9763-15f42c1676fc` (Enable user account).

Optional:

- `arguments` (Attributes Set) The arguments of the task, as defined by the task definition's parameters (for example `enableOnPremisesAccount` or `groupID`). (see [below for nested schema](#nestedatt--tasks--arguments))
- `continue_on_error` (Boolean) Whether the workflow continues with the next task when this task fails. Defaults to `false`.
- `description` (String) A description of the task. Defaults to an empty string.
- `is_enabled` (Boolean) Whether the task runs when the workflow runs. Defaults to `true`.

Read-Only:

- `execution_sequence` (Number) The position of the task in the workflow (1-based). Set from the order of `tasks`.

<a id="nestedatt--tasks--arguments"></a>
### Nested Schema for `tasks.arguments`

Required:

- `name` (String) The argument name.
- `value` (String) The argument value.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Updates and versions
- `display_name`, `description`, `is_enabled` and `is_scheduling_enabled` are updated with `PATCH`. The workflow version does not change.
- Changing `category`, `execution_conditions` or `tasks` calls `createNewVersion`. The workflow keeps its ID, `version` increases by one and Graph assigns new task IDs. The resource is never replaced.

### Execution conditions
- Set `on_demand_execution_only = true` for a workflow that only runs when started on demand. `scope` and `trigger` must then be omitted.
- Otherwise set both `scope` and `trigger`. `scope` takes either a `rule` (for example `(department eq 'Marketing')`) or `group_ids`.
- `attribute_change` accepts a single attribute; Graph rejects more than one trigger attribute.
- `membership_change` requires `scope.group_ids`.

### Tasks
- Tasks run in list order; `execution_sequence` is set from the position in `tasks`.
- `task_definition_id` must be a built-in task definition. Custom task extensions are not supported by this resource.
- Arguments must match the task definition's parameters. Graph validates argument values, for example user or group IDs must exist.

### Deletion
- Destroy deletes the workflow and then permanently removes it from `/identityGovernance/lifecycleWorkflows/deletedItems/workflows`. Deleted workflows otherwise stay restorable for 30 days and count toward the tenant's workflow limit. A workflow destroyed by Terraform cannot be restored.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using the workflow ID
terraform import microsoft365_graph_beta_identity_and_access_lifecycle_workflow.example 00000000-0000-0000-0000-000000000000
```
