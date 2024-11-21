FROM mcr.microsoft.com/microsoftgraph/powershell:latest
WORKDIR /app
VOLUME ["./scripts:/app/scripts"]
ENTRYPOINT ["pwsh", "-Command", "$tenantId = $env:TENANT_ID; $clientId = $env:CLIENT_ID; $clientSecret = $env:CLIENT_SECRET; . /app/scripts/ExportGraphPermissions.ps1 -TenantId $tenantId -ClientId $clientId -ClientSecret $clientSecret"]
