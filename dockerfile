FROM mcr.microsoft.com/microsoftgraph/powershell:latest

WORKDIR /app
COPY scripts/ExportGraphPermissions.ps1 /app/

ENTRYPOINT ["pwsh", "-Command", "$tenantId = $env:TENANT_ID; $clientId = $env:CLIENT_ID; $clientSecret = $env:CLIENT_SECRET; . /app/ExportGraphPermissions.ps1 -TenantId $tenantId -ClientId $clientId -ClientSecret $clientSecret"]