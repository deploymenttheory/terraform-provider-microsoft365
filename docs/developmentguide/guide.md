# Development Guide

## Resource Folder Hierarchy

The hierarchy for this project follows the graph / graph beta api documentation as of 05/09/2024. The design decision was taken to group the folders according to the documentation rather than the api endpoint path as it reflects current usage of the services.

The next decision was to then group resources by either graph beta (beta) or by graph (v1.0)

The next decision, beyond this top level grouping, is to ignore any sub categories as they are subject to change and will introduce a significant administrative overhead to maintain.

files and folder names are seperated with _

resources use the 

The resource folder hierarchy for this project is organized as follows:

```
internal/
└── resources/
    ├── users/
    │   ├── beta/
    │   │   └── user/
    │   │       ├── construct.go
    │   │       ├── crud.go
    │   │       ├── model.go
    │   │       ├── resource_model.json
    │   │       ├── resource.go
    │   │       ├── state.go
    │   │       └── validators.go
    │   └── v1.0/
    │       └── user/
    │           ├── construct.go
    │           ├── crud.go
    │           ├── model.go
    │           ├── resource_model.json
    │           ├── resource.go
    │           ├── state.go
    │           └── validators.go
    ├── groups/
    ├── applications/
    ├── backup_storage/
    ├── calendars/
    ├── change_notifications/
    ├── compliance/
    ├── cross-device_experiences/
    ├── customer_booking/
    ├── device_and_app_management/
    ├── education/
    ├── employee_experience/
    ├── extensions/
    ├── external_data_connections/
    ├── files/
    ├── financials/
    ├── identity_and_access/
    ├── industry_data_etl/
    ├── mail/
    ├── notes/
    ├── people_and_workplace_intelligence/
    ├── personal_contacts/
    ├── reports/
    ├── search/
    ├── security/
    ├── sites_and_lists/
    ├── tasks_and_lists/
    ├── tasks_and_plans/
    ├── teamwork_and_communications/
    ├── to-do_lists/
    └── workbooks_and_charts/
```

This structure helps to organize different types of resources within the project. Here's a brief explanation of each folder:

- `internal/`: This is the root folder for internal packages that are not meant to be imported by other projects.
  - `resources/`: This folder contains all the resource-related code.
    - `users/`: Contains code related to user resources.
    - `groups/`: Contains code related to group resources.
    - `applications/`: Contains code related to application resources.

Each resource type (users, groups, applications) can have its own set of files and subfolders as needed for implementation.