FROM mcr.microsoft.com/microsoftgraph/powershell:latest

WORKDIR /app
COPY scripts/ExportGraphPermissions.ps1 /app/

ENTRYPOINT ["pwsh", "-File", "/app/ExportGraphPermissions.ps1", "-TenantId", "$env:TENANT_ID", "-ClientId", "$env:CLIENT_ID", "-ClientSecret", "$env:CLIENT_SECRET"]